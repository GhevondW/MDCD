# Day 3 tests start real threads. A deadlocked solution would otherwise
# hang forever, so every test gets a time limit.
find_package(Threads REQUIRED)

foreach(name ticket_dispenser account bounded_stack once_cache event_bus concurrent_booker)
  add_executable(day03_${name}_test tests/${name}_test.cpp)
  target_include_directories(day03_${name}_test PRIVATE include)
  target_link_libraries(day03_${name}_test PRIVATE gtest_main Threads::Threads)
  gtest_discover_tests(day03_${name}_test TEST_PREFIX "day03_" PROPERTIES TIMEOUT 60)
endforeach()
