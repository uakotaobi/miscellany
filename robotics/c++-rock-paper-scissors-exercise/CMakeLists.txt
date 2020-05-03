PROJECT(roshambo CXX)
cmake_minimum_required(VERSION 2.8)

if(CMAKE_COMPILER_IS_GNUCXX)
    # Enable C++14, with (nearly) maximum warnings
    add_definitions(-Wall -Wextra --std=c++14)
endif(CMAKE_COMPILER_IS_GNUCXX)

# Cygwin kludge.
set(CXX "g++")

# Force CMake to recognize important subdirectories.
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY "./bin")
include_directories("./include")

# Add build target for the main game.  After compiling, run "bin/roshambo" to
# play.
add_executable(roshambo src/PlayerBase.cpp
                        src/main.cpp)
