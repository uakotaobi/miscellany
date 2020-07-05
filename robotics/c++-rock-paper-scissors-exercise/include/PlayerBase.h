#ifndef PLAYER_BASE_H__
#define PLAYER_BASE_H__

#include <string>

class PlayerBase {

    public:
        // Ensures that base-class objects are destroyed properly (if any.)
        // Be sure to implement this, even if the implementation is empty.
        virtual ~PlayerBase();

        // Make a move against the opponent.
        //
        // The function should return a string for which valid() is true.
        virtual std::string play() = 0;

        // OPTIONAL: Report to a player what moves both it and an opponent
        // made.  To prevent cheating, this is only called after a round of
        // the game is already finished.
        //
        // This can be used to compile historical data.
        virtual void remember(const std::string& myLastPlay, const std::string& theirLastPlay);

        // OPTIONAL: Returns the name of this player.  You may wish to
        // override this to make it more unique.
        virtual std::string name() const { return "NAME"; }
};

////////////////////////////////////
// Stand-alone utility functions. //
////////////////////////////////////

// Returns true if play is one of the valid moves for this game.  The function
// is case-insensitive.
bool valid(const std::string& play);

// Returns true if the first play beats the second.
//
// * defeats("paper", "rock") returns true.
// * defeats("rock", "scissors") returns true.
// * defeats("scissors", "paper") returns true.
//
// Any other permutation returns false.
//
// The print argument, if true, causes a sentence like "Paper covers
// rock" or "Tie" to be printed to standard output.  Note that the
// sentence is printed regardless of whether first beat second or
// whether second beat first.
//
// An invalid argument (like "nuke") throws an exception.
bool defeats(const std::string& first, const std::string& second, bool print=false);

#endif // (#ifndef PLAYER_BASE_H__)
