PROJECT(roshambo CXX)
cmake_minimum_required(VERSION 2.8)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_FLAGS "-Wall -Wextra")

# Cygwin kludge.
set(CXX "g++")

# Force CMake to recognize important subdirectories.
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY "./bin")
include_directories("./include")

# Add build target for the main game.  After compiling, run "bin/roshambo" to
# play.
set(SOURCES src/PlayerBase.cpp src/main.cpp src/ComputerPlayer.cpp src/HumanPlayer.cpp)
add_executable(roshambo ${SOURCES})
