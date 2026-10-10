#include "day04/thread_pool.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <atomic>
#include <chrono>
#include <future>
#include <mutex>
#include <set>
#include <stdexcept>
#include <string>
#include <thread>
#include <vector>

using namespace std::chrono_literals;

namespace {
void sleepMs(int n) { std::this_thread::sleep_for(std::chrono::milliseconds(n)); }
}  // namespace

// ---- running tasks ----

TEST(ThreadPool, ZeroWorkersIsRejected) {
    EXPECT_THROW({ ThreadPool pool(0); }, std::invalid_argument);
}

TEST(ThreadPool, ReportsItsWorkerCount) {
    ThreadPool pool(3);
    EXPECT_EQ(pool.workerCount(), 3u);
}

TEST(ThreadPool, RunsEveryTask) {
    std::atomic<int> count{0};
    {
        ThreadPool pool(4);
        for (int i = 0; i < 10000; ++i) pool.execute([&] { count.fetch_add(1); });
        pool.shutdown();
    }
    EXPECT_EQ(count.load(), 10000);
}

TEST(ThreadPool, SubmitReturnsTheResultInAFuture) {
    ThreadPool pool(2);
    std::future<int> answer = pool.submit([] { return 6 * 7; });
    std::future<std::string> text = pool.submit([] { return std::string("done"); });
    ASSERT_EQ(answer.wait_for(3s), std::future_status::ready);
    EXPECT_EQ(answer.get(), 42);
    EXPECT_EQ(text.get(), "done");
}

TEST(ThreadPool, ATaskThatThrowsDeliversTheExceptionThroughItsFuture) {
    ThreadPool pool(2);
    std::future<int> bad = pool.submit([]() -> int { throw std::logic_error("boom"); });
    EXPECT_THROW(bad.get(), std::logic_error);
    // The pool is still alive.
    std::future<int> good = pool.submit([] { return 1; });
    ASSERT_EQ(good.wait_for(3s), std::future_status::ready);
    EXPECT_EQ(good.get(), 1);
}

TEST(ThreadPool, ARawTaskThatThrowsDoesNotKillItsWorker) {
    // With one worker, a dead worker means no task ever runs again. (An
    // exception that escapes a std::thread ends the whole program.)
    ThreadPool pool(1);
    pool.execute([] { throw std::runtime_error("boom"); });
    std::atomic<bool> ran{false};
    pool.execute([&] { ran = true; });
    EXPECT_TRUE(waitUntil([&] { return ran.load(); }, 3s)) << "the worker did not survive";
}

TEST(ThreadPool, TasksRunInParallel) {
    // 4 tasks that each wait until all 4 are running: they only finish if
    // 4 workers really run at the same time.
    constexpr int kWorkers = 4;
    ThreadPool pool(kWorkers);
    std::atomic<int> arrived{0};
    std::atomic<int> passed{0};
    for (int i = 0; i < kWorkers; ++i) {
        pool.execute([&] {
            arrived.fetch_add(1);
            if (waitUntil([&] { return arrived.load() == kWorkers; }, 3s)) passed.fetch_add(1);
        });
    }
    pool.shutdown();
    EXPECT_EQ(passed.load(), kWorkers) << "the tasks did not run at the same time";
}

TEST(ThreadPool, NeverUsesMoreThreadsThanWorkers) {
    ThreadPool pool(3);
    std::mutex m;
    std::set<std::thread::id> ids;
    for (int i = 0; i < 60; ++i) {
        pool.execute([&] {
            sleepMs(2);
            std::lock_guard<std::mutex> lk(m);
            ids.insert(std::this_thread::get_id());
        });
    }
    pool.shutdown();
    EXPECT_LE(ids.size(), 3u) << "more threads than workers ran tasks";
    EXPECT_GE(ids.size(), 2u) << "the tasks did not spread over the workers";
    EXPECT_EQ(ids.count(std::this_thread::get_id()), 0u) << "a task ran on the caller's thread";
}

TEST(ThreadPool, ATaskCanSubmitAnotherTask) {
    ThreadPool pool(2);
    std::promise<int> p;
    std::future<int> f = p.get_future();
    pool.execute([&] { pool.execute([&] { p.set_value(7); }); });
    ASSERT_EQ(f.wait_for(3s), std::future_status::ready);
    EXPECT_EQ(f.get(), 7);
}

TEST(ThreadPool, ManySubmittersAtOnce) {
    // 8 threads submit 2,000 tasks each into a small queue of a small pool.
    constexpr int kSubmitters = 8;
    constexpr int kPerSubmitter = 2000;
    std::atomic<int> count{0};
    ThreadPool pool(4, 16);
    runTogetherWithin(30, kSubmitters, [&](int) {
        for (int i = 0; i < kPerSubmitter; ++i) pool.execute([&] { count.fetch_add(1); });
    });
    pool.shutdown();
    EXPECT_EQ(count.load(), kSubmitters * kPerSubmitter);
}

// ---- the queue limit ----

TEST(ThreadPool, ExecuteWaitsWhileTheQueueIsFull) {
    Gate gate;
    std::atomic<int> count{0};
    ThreadPool pool(1, 2);
    std::atomic<bool> started{false};
    pool.execute([&] {
        started = true;
        gate.wait();
        count.fetch_add(1);
    });
    ASSERT_TRUE(waitUntil([&] { return started.load(); }, 3s));
    pool.execute([&] { count.fetch_add(1); });  // waits in the queue
    pool.execute([&] { count.fetch_add(1); });  // waits in the queue: now full

    Background<bool> third([&] {
        pool.execute([&] { count.fetch_add(1); });
        return true;
    });
    EXPECT_FALSE(third.waitFor(150ms)) << "execute() returned although the queue was full";

    gate.open();
    EXPECT_TRUE(third.waitFor(3s)) << "execute() was not woken up when room appeared";
    third.result();
    pool.shutdown();
    EXPECT_EQ(count.load(), 4);
}

TEST(ThreadPool, TrySubmitRejectsWhenTheQueueIsFull) {
    Gate gate;
    std::atomic<int> ranA{0}, ranB{0};
    ThreadPool pool(1, 1);
    std::atomic<bool> started{false};
    pool.execute([&] {
        started = true;
        gate.wait();
    });
    ASSERT_TRUE(waitUntil([&] { return started.load(); }, 3s));

    EXPECT_TRUE(pool.trySubmit([&] { ranA.fetch_add(1); })) << "there was room";
    EXPECT_FALSE(pool.trySubmit([&] { ranB.fetch_add(1); })) << "the queue was full";

    gate.open();
    pool.shutdown();
    EXPECT_EQ(ranA.load(), 1);
    EXPECT_EQ(ranB.load(), 0) << "a rejected task must not run";
}

// ---- shutdown ----

TEST(ThreadPool, ShutdownRunsEveryAcceptedTask) {
    std::atomic<int> count{0};
    ThreadPool pool(2);
    for (int i = 0; i < 100; ++i) {
        pool.execute([&] {
            sleepMs(1);
            count.fetch_add(1);
        });
    }
    pool.shutdown();
    EXPECT_EQ(count.load(), 100) << "shutdown() returned before the queue was empty";
}

TEST(ThreadPool, TheDestructorShutsDownToo) {
    std::atomic<int> count{0};
    {
        ThreadPool pool(2);
        for (int i = 0; i < 50; ++i) {
            pool.execute([&] {
                sleepMs(1);
                count.fetch_add(1);
            });
        }
    }
    EXPECT_EQ(count.load(), 50);
}

TEST(ThreadPool, ShutdownStopsIdleWorkers) {
    ThreadPool pool(8);
    std::atomic<int> count{0};
    for (int i = 0; i < 8; ++i) pool.execute([&] { count.fetch_add(1); });
    ASSERT_TRUE(waitUntil([&] { return count.load() == 8; }, 3s));
    sleepMs(50);  // let the workers fall asleep
    Background<bool> stopper([&] {
        pool.shutdown();
        return true;
    });
    EXPECT_TRUE(stopper.waitFor(3s)) << "shutdown() did not wake the sleeping workers";
    stopper.result();
}

TEST(ThreadPool, ShutdownIsIdempotent) {
    ThreadPool pool(2);
    std::atomic<int> count{0};
    pool.execute([&] { count.fetch_add(1); });
    pool.shutdown();
    pool.shutdown();
    EXPECT_EQ(count.load(), 1);
}

TEST(ThreadPool, ANewTaskAfterShutdownIsRejected) {
    ThreadPool pool(2);
    pool.shutdown();
    std::atomic<bool> ran{false};
    EXPECT_THROW(pool.execute([&] { ran = true; }), std::runtime_error);
    EXPECT_FALSE(pool.trySubmit([&] { ran = true; }));
    EXPECT_THROW(pool.submit([] { return 1; }), std::runtime_error);
    sleepMs(20);
    EXPECT_FALSE(ran.load());
}

TEST(ThreadPool, ShutdownWakesASubmitterThatWaitsForRoom) {
    Gate gate;
    ThreadPool pool(1, 1);
    std::atomic<bool> started{false};
    pool.execute([&] {
        started = true;
        gate.wait();
    });
    ASSERT_TRUE(waitUntil([&] { return started.load(); }, 3s));
    pool.execute([] {});  // the queue is full now

    Background<bool> blocked([&] {
        try {
            pool.execute([] {});
            return false;
        } catch (const std::runtime_error&) {
            return true;
        }
    });
    EXPECT_FALSE(blocked.waitFor(150ms)) << "execute() should wait while the queue is full";

    Background<bool> stopper([&] {
        pool.shutdown();
        return true;
    });
    EXPECT_TRUE(blocked.waitFor(3s)) << "shutdown() did not wake the waiting submitter";
    gate.open();
    EXPECT_TRUE(stopper.waitFor(3s));
    EXPECT_TRUE(stopper.result());
    EXPECT_TRUE(blocked.result()) << "the woken execute() must throw std::runtime_error";
}
