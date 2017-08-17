#include <sstream>
#include "ComputerPlayer.h"

using std::ostringstream;
using std::string;

int ComputerPlayer::counter = 0;

ComputerPlayer::ComputerPlayer() : id_(++counter), name_() {
    ostringstream stream;
    stream << "CPU #" << id_;
    name_ = stream.str();
}

std::string ComputerPlayer::name() const { return name_; }

string ComputerPlayer::play() {

    // Make a very sophisticated play.
    switch (id_ % 3) {
        case 0:
            return "rock";
        case 1:
            return "paper";
        case 2:
        default:
            return "scissors";
    }
}

void ComputerPlayer::remember(const string& myLastPlay, const string& theirLastPlay) {
    // Doesn't do anything right now.
}
