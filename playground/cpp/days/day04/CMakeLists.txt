# Day 4 tests start real threads. A deadlocked solution would otherwise
# hang forever, so every test gets a time limit.
find_package(Threads REQUIRED)

foreach(name bounded_blocking_queue thread_pool work_stealing_pool)
  add_executable(day04_${name}_test tests/${name}_test.cpp)
  target_include_directories(day04_${name}_test PRIVATE include)
  target_link_libraries(day04_${name}_test PRIVATE gtest_main Threads::Threads)
  gtest_discover_tests(day04_${name}_test TEST_PREFIX "day04_" PROPERTIES TIMEOUT 60)
endforeach()
