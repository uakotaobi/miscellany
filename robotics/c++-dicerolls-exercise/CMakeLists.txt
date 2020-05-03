PROJECT(diceroll C CXX)

# Not using find_package() here, but if we were, we'd need CMake v2.8+.
cmake_minimum_required(VERSION 2.8)

if(UNIX)
  if(CMAKE_COMPILER_IS_GNUCXX)
    # Add compiler flags to turn on (most) warnings and to use whichever
    # aspects of the C++14 standard have so far been implemented.
    add_definitions(-Wall -Wextra --std=c++14)
  endif(CMAKE_COMPILER_IS_GNUCXX)
endif(UNIX)

set(CXX "g++") # Only Cygwin needs this.

add_executable(diceroll ./diceroll.cpp)
