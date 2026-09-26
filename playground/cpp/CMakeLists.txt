cmake_minimum_required(VERSION 3.16)
project(mdcd_playground CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

include(cmake/CPM.cmake)

# MSVC/GoogleTest CRT mismatch fix: must be set before GoogleTest is configured.
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)

CPMAddPackage(
  NAME googletest
  GITHUB_REPOSITORY google/googletest
  GIT_TAG v1.18.0
  OPTIONS "INSTALL_GTEST OFF"
)

enable_testing()
include(GoogleTest)

add_subdirectory(days/day02)
add_subdirectory(days/day03)
