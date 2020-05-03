#ifndef COMPUTER_PLAYER_H__
#define COMPUTER_PLAYER_H__

#include <string>
#include "PlayerBase.h"

using std::string;

// A computer player that isn't very good.  Please make it better.
class ComputerPlayer : public PlayerBase {
    public:

        ComputerPlayer();
        ~ComputerPlayer() { }
        std::string name() const;

    public:

        void remember(const string& myLastPlay, const string& theirLastPlay);
        std::string play();

    private:

        static int counter;
        int id_;
        string name_;

};

#endif // #ifndef COMPUTER_PLAYER_H__
