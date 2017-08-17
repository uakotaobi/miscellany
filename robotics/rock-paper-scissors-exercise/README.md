This is a simple program that plays a game of rock-paper-scissors between two
players -- at least, that's what it's supposed to do.  However, it has several
problems that it is up to you to fix.

The purpose of this exercise is to learn how to read and modify existing,
buggy code to order to develop a new feature.

## Phase 1: Get the code to build

Your first task is to **fix the build error** we're having.  It has something to
do with CMake, but the previous developer wasn't sure now to fix the problem.
But you've done this sort of thing before; we're confident you'll figure it
out.

Once you have completed this task, please use Git to tag your code with the
tag `phase1`.

## Phase 2: Implement a human player

Now that you've run the program, you'll notice that it's _boring_; it only has
two computer players and the human user never gets a turn.

Please fix this by **implementing a human player.** Use the header in
[PlayerBase.h](/robotics/rock-paper-scissors-exercise/include/PlayerBase.h) as a guide.  You should write a class that inheritcs from `PlayerBase`, then create both a header file ("human.h" or
the like) and an implementation file ("human.cpp") for it.  Update the CMake build
accordingly.  At the very minimum, your human player should have a destructor
function (even if it does nothing) and a `play()` function that asks the human
player what he or she wants to choose and returns that as a std::string.

`PlayerBase.h` supplies utility functions that can help
you, but you don't have to call them if you don't want to:

```C++
bool valid(const std::string& play)
```

Returns true if the input, converted to lowercase, matches "rock", "paper",
or "scissors", and false otherwise.

```C++
bool defeats(const string& first, const string& second, bool print)
```

Returns true if "first beats second".  So, for instance, calling
`defeats("rock", "paper")` will return false because rock doesn't beat
paper.

The `print` argument is optional, and defaults to false.  All it does is
print a message like "Paper covers rock."

Once you have completed this task, please use Git to tag your code with the
tag `phase2`.


## Phase 3: Implement a more interesting AI player

Have you noticed a pattern in the way the computer players make their moves?
It's very repetitive.  See what you can do to make the computer player harder to anticipate.  This doesn't have to be anything sophisticated (that will come later)&mdash;for now, it just needs to be unpredictable.

Once you have completed this task, please use Git to tag your code with the
tag "phase3".

## Phase 4: Rock-Paper-Scissors-Lizard-Spock

Change the rules of the game so that:

* Rock smashes scissors and crushes **lizard**;
* Paper covers rock and disproves **Spock**;
* Scissors cuts paper and maims **lizard**;
* **Lizard** poisons **Spock** and eats paper;
* **Spock** smashes scissors and vaporizes rock.

Be sure to update the `ComputerPlayer` and `HumanPlayer` classes.

Once you have completed this task, please use Git to tag your code with the
tag `phase4`.

## Phase 5 (optional)

It turns out that humans are terrible at Rock-Paper-Scissors precisely because
they do *not* play at random.

To exploit this, your ComputerPlayer object has a `remember()` function that doesn't do anything
right now.  It is called at the end of every round, revealing both what the
opponent and the computer player itself played.

Can you find a way to use `remember()` to make the computer player stronger?

If you complete this task, please use Git to tag your code with the tag
`phase5`.
