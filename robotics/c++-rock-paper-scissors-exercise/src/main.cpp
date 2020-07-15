#include <vector>
#include <memory>   // unique_ptr
#include <string>
#include <sstream>  // istringstream, ostringstream
#include <iostream>
#include "HumanPlayer.h"
#include "ComputerPlayer.h"

using std::istringstream;
using std::ostringstream;
using std::unique_ptr;
using std::vector;
using std::string;
using std::cout;

// Source: https://www.asciiart.eu/people/body-parts/hand-gestures
// Artist: Veronica Karlsson (VK)
string graphics[] = {
        // 0: Paper, left
        "    _______\n"
        "---'   ____)____\n"
        "          ______)\n"
        "          _______)\n"
        "         _______)\n"
        "---.__________)",
        // 1: Paper, right
        "       _______    \n"
        "  ____(____   `---\n"
        " (______          \n"
        "(_______          \n"
        " (_______         \n"
        "   (__________,---",
        // 2: Scissors, left
        "    _______       \n"
        "---'   ____)____  \n"
        "          ______) \n"
        "       __________)\n"
        "      (____)      \n"
        "---.__(___)       ",
        // 3: Scissors, right
        "       _______    \n"
        "  ____(____   `---\n"
        " (______          \n"
        "(__________       \n"
        "      (____)      \n"
        "       (___)__,---",
        // 4: Rock, left
        "    _______  \n"
        "---'   ____) \n"
        "      (_____)\n"
        "      (_____)\n"
        "      (____) \n"
        "---.__(___)  ",
        // 5: Rock, right
        "  _______    \n"
        " (____   `---\n"
        "(_____)      \n"
        "(_____)      \n"
        " (____)      \n"
        "  (___)__,---"
};

// Prints a message stating what the given players played, and who won as a
// result.
//
// If the verbosity flag is 0, the printed message will be minimal.
//
// If the verbosity flag is 2 or higher, the printed message is more elaborate
// and includes ASCII art.
bool printOutcome(const string& leftPlayerName, const string& leftPlay,
                  const string& rightPlayerName, const string& rightPlay, int verbosity) {


    bool leftPlayerDefeatedRightPlayer = false;


    switch(verbosity) {
        case 0:

            cout << ">> " << leftPlayerName << ": " << leftPlay << ", " << rightPlayerName << ": " << rightPlay << " | ";
            leftPlayerDefeatedRightPlayer = defeats(leftPlay, rightPlay);
            if (leftPlay != rightPlay) {
                cout << ", "
                     << (leftPlayerDefeatedRightPlayer ? leftPlayerName : rightPlayerName)
                     << " wins";
            }
            cout << "\n";
            break;

        case 1:

            cout << ">> " << leftPlayerName << ": " << leftPlay << "\n";
            cout << ">> " << rightPlayerName << ": " << rightPlay << "\n";

            cout << "\n  ";
            leftPlayerDefeatedRightPlayer = defeats(leftPlay, rightPlay);
            cout << ".";

            if (leftPlay != rightPlay) {
                cout << "\n  "
                     << (leftPlayerDefeatedRightPlayer ? leftPlayerName : rightPlayerName)
                     << " wins.";
            }
            cout << "\n";
            break;

        default: {

            int leftIndex = (leftPlay == "rock" ? 4 : (leftPlay == "scissors" ? 2 : 0));
            int rightIndex = (rightPlay == "rock" ? 5 : (rightPlay == "scissors" ? 3 : 1));

            const int WIDTH = 70;

            // Print the name line.
            cout << "\n" << leftPlayerName
                 << string(WIDTH - leftPlayerName.length() - rightPlayerName.length(), ' ')
                 << rightPlayerName << "\n\n";

            // Print the massive hand graphics.
            istringstream leftStream(graphics[leftIndex]), rightStream(graphics[rightIndex]);
            string leftLine, rightLine;
            int lineNumber = 0;

            while (leftStream && rightStream) {
                string line = string(70, ' ');

                if (getline(leftStream, leftLine)) {
                    line = leftLine + line.substr(leftLine.size());
                }

                if (getline(rightStream, rightLine)) {
                    line = line.substr(0, WIDTH - rightLine.size()) + rightLine;
                }

                if (lineNumber == 3) {
                    // Line 4 is kind of in the middle.
                    ostringstream out;
                    leftPlayerDefeatedRightPlayer = defeats(leftPlay, rightPlay, out);
                    string message = out.str();
                    if (message.length() % 2 != 0) {
                        // Force message length to be even to make
                        // calculations easier.
                        message += " ";
                    }

                    // Print an arrow that points toward the defeated player.
                    if (leftPlayerDefeatedRightPlayer) {
                        message = "    " + message;
                        if (message.back() != ' ') {
                            message += " -> ";
                        } else {
                            message += "->";
                        }
                    } else if (leftPlay != rightPlay) {
                        message = " <- " + message;
                        if (message.back() != ' ') {
                            message += "    ";
                        } else {
                            message += "  ";
                        }
                    }

                    // Graft the message into line 4.
                    size_t lineMid = line.length() / 2;
                    size_t messageMid = message.length() / 2;
                    line = line.substr(0, lineMid - messageMid) + message + line.substr(lineMid + messageMid);

                    // Add the winner string to the appropriate "hand."
                    if (leftPlayerDefeatedRightPlayer) {
                        line = "WINS" + line.substr(4);
                    } else if (leftPlay != rightPlay) {
                        line = line.substr(0, WIDTH - 4) + "WINS";
                    }
                }

                cout << line << "\n";
                lineNumber += 1;
            }
            break;
        }
    }

    return leftPlayerDefeatedRightPlayer;
}


// A simple driver for a simple game.
int main() {

    const int winsNeeded = 5;
    const int maxRounds = 2 * winsNeeded - 1;

    // No winner yet.
    int winner = -1;
    cout << "Best " << winsNeeded << " out of " << maxRounds << ".  Go! " << string(40, '-') << "\n";

    // The game only supports two players.
    vector<unique_ptr<PlayerBase>> players;
    players.push_back(unique_ptr<PlayerBase> (new HumanPlayer()));
    players.push_back(unique_ptr<PlayerBase> (new ComputerPlayer()));

    vector<int> victories(players.size());

    while (winner == -1) {
        cout << "\n";

        string a = players[0]->play();
        string b = players[1]->play();

        bool firstPlayerDefeatedSecondPlayer = printOutcome(players[0]->name(),
                                                            a,
                                                            players[1]->name(),
                                                            b,
                                                            2);

        // Allow AI players to remember the history of the other players'
        // moves, but only after they have been made.  (The computer players
        // cannot cheat.)
        players[0]->remember(a, b);
        players[1]->remember(b, a);

        // If we didn't tie then either the first player beat the second
        // player or vice versa.
        if (a != b) {

            int roundVictorIndex = (firstPlayerDefeatedSecondPlayer ? 0 : 1);
            victories[roundVictorIndex]++;
            cout << "  The score is " << victories[0] << " ("
                 << players[0]->name() << ") to " << victories[1] << " ("
                 << players[1]->name() << ").\n";

            if (victories[roundVictorIndex] >= winsNeeded) {
                winner = roundVictorIndex;
            }
        }

    } // end (while no one has yet won the game)

    cout << "\n" << players[winner]->name() << " wins the game! " << string(40, '-') << "\n";
}
