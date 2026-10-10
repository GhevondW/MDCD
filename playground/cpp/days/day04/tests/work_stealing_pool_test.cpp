#include "day04/work_stealing_pool.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <algorithm>
#include <atomic>
#include <chrono>
#include <mutex>
#include <numeric>
#include <set>
#include <stdexcept>
#include <string>
#include <thread>
#include <vector>

using namespace std::chrono_literals;
using Clock = std::chrono::steady_clock;

namespace {

void sleepMs(int n) { std::this_thread::sleep_for(std::chrono::milliseconds(n)); }

// fib(n) with one task per call -- the recursive code a plain pool
// deadlocks on, because every task waits for its child.
long fib(WorkStealingPool& pool, int n) {
    if (n < 2) return n;
    auto child = pool.submit([&pool, n] { return fib(pool, n - 1); });
    long other = fib(pool, n - 2);
    return child.get() + other;
}

// Sum of v[lo, hi) by splitting the range in two until it is small.
long long sum(WorkStealingPool& pool, const std::vector<int>& v, std::size_t lo, std::size_t hi) {
    if (hi - lo <= 10000) return std::accumulate(v.begin() + lo, v.begin() + hi, 0LL);
    std::size_t mid = lo + (hi - lo) / 2;
    auto left = pool.submit([&pool, &v, lo, mid] { return sum(pool, v, lo, mid); });
    long long right = sum(pool, v, mid, hi);
    return left.get() + right;
}

// Every task counts itself and forks two more, down to depth 0, and
// never waits for them: 2^(depth+1) - 1 tasks in all.
void spawnTree(WorkStealingPool& pool, std::atomic<int>& count, int depth) {
    count.fetch_add(1);
    if (depth == 0) return;
    pool.execute([&pool, &count, depth] { spawnTree(pool, count, depth - 1); });
    pool.execute([&pool, &count, depth] { spawnTree(pool, count, depth - 1); });
}

void checkFib(std::size_t workers) {
    WorkStealingPool pool(workers);
    Background<long> run([&] { return pool.submit([&] { return fib(pool, 20); }).get(); });
    EXPECT_EQ(run.result(30), 6765) << "with " << workers << " worker(s)";
}

void checkSum(std::size_t workers) {
    std::vector<int> v(1000000);
    std::iota(v.begin(), v.end(), 1);
    WorkStealingPool pool(workers);
    Background<long long> run([&] {
        return pool.submit([&] { return sum(pool, v, 0, v.size()); }).get();
    });
    EXPECT_EQ(run.result(30), 500000500000LL) << "with " << workers << " worker(s)";
}

}  // namespace

// ---- fork/join work: the join must not block the worker ----

TEST(WorkStealingPool, ZeroWorkersIsRejected) {
    EXPECT_THROW({ WorkStealingPool pool(0); }, std::invalid_argument);
}

TEST(WorkStealingPool, ReportsItsWorkerCount) {
    WorkStealingPool pool(3);
    EXPECT_EQ(pool.workerCount(), 3u);
}

TEST(WorkStealingPool, RecursiveFibWorksOnOneWorker) {
    // Every task waits for its child. A worker that sleeps in get() has
    // nobody left to run the child: this is the deadlock from the slides.
    checkFib(1);
}

TEST(WorkStealingPool, RecursiveFibWorksOnTwoWorkers) { checkFib(2); }

TEST(WorkStealingPool, RecursiveFibWorksOnFourWorkers) { checkFib(4); }

TEST(WorkStealingPool, DivideAndConquerSumOnOneWorker) { checkSum(1); }

TEST(WorkStealingPool, DivideAndConquerSumOnFourWorkers) { checkSum(4); }

TEST(WorkStealingPool, ParentsWaitingForChildrenDoNotDeadlock) {
    // The slides' example: 2 workers, 2 parents that each wait for a child.
    WorkStealingPool pool(2);
    auto parent = [&pool] {
        auto child = pool.submit([] { return 1; });
        return child.get() + 1;
    };
    Background<int> run([&] {
        auto a = pool.submit(parent);
        auto b = pool.submit(parent);
        return a.get() + b.get();
    });
    EXPECT_EQ(run.result(30), 4);
}

TEST(WorkStealingPool, AnExceptionTravelsThroughTheHandle) {
    WorkStealingPool pool(2);
    Background<bool> run([&] {
        auto bad = pool.submit([]() -> int { throw std::logic_error("boom"); });
        try {
            bad.get();
            return false;
        } catch (const std::logic_error&) {
            return true;
        }
    });
    EXPECT_TRUE(run.result(30)) << "get() must rethrow the task's exception";
    // The pool is still alive.
    Background<int> after([&] { return pool.submit([] { return 5; }).get(); });
    EXPECT_EQ(after.result(30), 5);
}

TEST(WorkStealingPool, HelpUntilFromAnOutsideThreadWaitsForTheCondition) {
    WorkStealingPool pool(2);
    std::atomic<bool> flag{false};
    pool.execute([&] {
        sleepMs(50);
        flag = true;
    });
    Background<bool> run([&] {
        pool.helpUntil([&] { return flag.load(); });
        return flag.load();
    });
    EXPECT_TRUE(run.result(10)) << "helpUntil returned before the condition was true";
}

// ---- where tasks go and who runs them ----

TEST(WorkStealingPool, TasksFromOutsideRunOnWorkersNotOnTheCaller) {
    WorkStealingPool pool(3);
    std::mutex m;
    std::set<std::thread::id> ids;
    Background<int> run([&] {
        std::vector<WorkStealingPool::Handle<int>> handles;
        for (int i = 0; i < 200; ++i) {
            handles.push_back(pool.submit([&] {
                std::lock_guard<std::mutex> lk(m);
                ids.insert(std::this_thread::get_id());
                return 1;
            }));
        }
        int n = 0;
        for (auto& h : handles) n += h.get();
        return n;
    });
    EXPECT_EQ(run.result(30), 200);
    EXPECT_LE(ids.size(), 3u) << "more threads than workers ran tasks";
    EXPECT_EQ(ids.count(std::this_thread::get_id()), 0u);
}

TEST(WorkStealingPool, TasksFromOutsideAreNotSteals) {
    // They wait in the shared injection queue; taking from it is not
    // stealing from another worker.
    WorkStealingPool pool(4);
    Background<int> run([&] {
        std::vector<WorkStealingPool::Handle<int>> handles;
        for (int i = 0; i < 300; ++i) handles.push_back(pool.submit([] { return 1; }));
        int n = 0;
        for (auto& h : handles) n += h.get();
        return n;
    });
    EXPECT_EQ(run.result(30), 300);
    EXPECT_EQ(pool.steals(), 0u);
}

TEST(WorkStealingPool, AWorkerRunsItsOwnNewestTaskFirst) {
    // One worker, nobody to steal: the tasks a task forks wait in its
    // worker's deque and come out newest first.
    WorkStealingPool pool(1);
    std::mutex m;
    std::string order;
    pool.execute([&] {
        for (char c : std::string("ABC")) {
            pool.execute([&, c] {
                std::lock_guard<std::mutex> lk(m);
                order += c;
            });
        }
    });
    pool.shutdown();
    EXPECT_EQ(order, "CBA");
    EXPECT_EQ(pool.steals(), 0u);
}

TEST(WorkStealingPool, AThiefTakesTheOldestTaskFirst) {
    // Two workers. One runs a task that forks A, B, C and then stays busy
    // without helping (it just waits). The other worker has nothing to do,
    // so it steals -- from the front of the busy worker's deque.
    WorkStealingPool pool(2);
    std::mutex m;
    std::string order;
    std::atomic<int> finished{0};
    pool.execute([&] {
        for (char c : std::string("ABC")) {
            pool.execute([&, c] {
                {
                    std::lock_guard<std::mutex> lk(m);
                    order += c;
                }
                finished.fetch_add(1);
            });
        }
        waitUntil([&] { return finished.load() == 3; }, 3s);
    });
    pool.shutdown();
    EXPECT_EQ(order, "ABC") << "thieves must take the oldest task, not the newest";
    EXPECT_GE(pool.steals(), 1u);
}

namespace {
// The typical time of sleep_for(5ms) on this machine (it is much more than
// 5 ms on some systems).
double sleepUnitMs() {
    std::vector<double> samples;
    for (int i = 0; i < 5; ++i) {
        auto t0 = Clock::now();
        sleepMs(5);
        samples.push_back(std::chrono::duration<double, std::milli>(Clock::now() - t0).count());
    }
    std::sort(samples.begin(), samples.end());
    return samples[2];
}
}  // namespace

TEST(WorkStealingPool, ASkewedWorkloadSpreadsOverAllWorkers) {
    // All the work starts on ONE worker: a single task forks 120 children
    // that each sleep 5 ms, then joins them. The other 3 workers have to
    // steal to help. Alone it would take 120 sleeps; with 4 workers about
    // a quarter of that.
    constexpr int kChildren = 120;
    const double unit = sleepUnitMs();
    WorkStealingPool pool(4);
    std::mutex m;
    std::set<std::thread::id> ids;

    const auto start = Clock::now();
    Background<int> run([&] {
        return pool
            .submit([&] {
                std::vector<WorkStealingPool::Handle<int>> kids;
                for (int i = 0; i < kChildren; ++i) {
                    kids.push_back(pool.submit([&] {
                        {
                            std::lock_guard<std::mutex> lk(m);
                            ids.insert(std::this_thread::get_id());
                        }
                        sleepMs(5);
                        return 1;
                    }));
                }
                int n = 0;
                for (auto& k : kids) n += k.get();
                return n;
            })
            .get();
    });
    EXPECT_EQ(run.result(30), kChildren);
    const double elapsedMs = std::chrono::duration<double, std::milli>(Clock::now() - start).count();

    EXPECT_GT(pool.steals(), 0u) << "nobody stole anything";
    EXPECT_GE(ids.size(), 2u) << "all the children ran on one worker";
    EXPECT_LT(elapsedMs, 0.6 * kChildren * unit)
        << "took " << elapsedMs << " ms; one worker alone needs about " << kChildren * unit
        << " ms, so the other workers did not help";
}

// ---- shutdown ----

TEST(WorkStealingPool, ShutdownWaitsForWorkThatRunningTasksFork) {
    // The first task forks two, each of those forks two more, and so on,
    // and nobody waits for anybody. shutdown() is called right away: it
    // must not return while any of the 2,047 tasks is still to come.
    std::atomic<int> count{0};
    WorkStealingPool pool(4);
    pool.execute([&] { spawnTree(pool, count, 10); });
    pool.shutdown();
    EXPECT_EQ(count.load(), 2047);
}

TEST(WorkStealingPool, TheDestructorShutsDownToo) {
    std::atomic<int> count{0};
    {
        WorkStealingPool pool(3);
        pool.execute([&] { spawnTree(pool, count, 8); });
    }
    EXPECT_EQ(count.load(), 511);
}

TEST(WorkStealingPool, ShutdownIsIdempotentAndStopsIdleWorkers) {
    WorkStealingPool pool(8);
    std::atomic<int> count{0};
    for (int i = 0; i < 8; ++i) pool.execute([&] { count.fetch_add(1); });
    ASSERT_TRUE(waitUntil([&] { return count.load() == 8; }, 3s));
    sleepMs(50);  // let the workers fall asleep
    Background<bool> stopper([&] {
        pool.shutdown();
        pool.shutdown();
        return true;
    });
    EXPECT_TRUE(stopper.waitFor(3s)) << "shutdown() did not wake the idle workers";
    stopper.result();
}

TEST(WorkStealingPool, ANewTaskFromOutsideAfterShutdownIsRejected) {
    WorkStealingPool pool(2);
    pool.shutdown();
    std::atomic<bool> ran{false};
    EXPECT_THROW(pool.execute([&] { ran = true; }), std::runtime_error);
    EXPECT_THROW(pool.submit([] { return 1; }), std::runtime_error);
    sleepMs(20);
    EXPECT_FALSE(ran.load());
}

TEST(WorkStealingPool, ARawTaskThatThrowsDoesNotKillItsWorker) {
    WorkStealingPool pool(1);
    pool.execute([] { throw std::runtime_error("boom"); });
    std::atomic<bool> ran{false};
    pool.execute([&] { ran = true; });
    EXPECT_TRUE(waitUntil([&] { return ran.load(); }, 3s)) << "the worker did not survive";
}

TEST(WorkStealingPool, ManyThreadsSubmitAtOnce) {
    constexpr int kSubmitters = 8;
    constexpr int kPerSubmitter = 2000;
    std::atomic<int> count{0};
    WorkStealingPool pool(4);
    runTogetherWithin(30, kSubmitters, [&](int) {
        for (int i = 0; i < kPerSubmitter; ++i) pool.execute([&] { count.fetch_add(1); });
    });
    pool.shutdown();
    EXPECT_EQ(count.load(), kSubmitters * kPerSubmitter);
}
