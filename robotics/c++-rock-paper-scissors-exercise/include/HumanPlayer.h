#ifndef HUMAN_PLAYER_H__
#define HUMAN_PLAYER_H__

#include <string>
#include "PlayerBase.h"

using std::string;

// A human player.
class HumanPlayer : public PlayerBase {
    public:

        HumanPlayer();
        ~HumanPlayer();
        std::string name() const;

    public:

        std::string play();

        // Your brain will do the remembering.
        void remember(const std::string&, const std::string&) { }

    private:
        string name_;

};

#endif // #ifndef HUMAN_PLAYER_H__
