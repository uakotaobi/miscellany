#include <vector>
#include <memory>
#include <string>
#include <iostream>
#include "ComputerPlayer.h"

using std::unique_ptr;
using std::vector;
using std::string;
using std::cout;

// A simple driver for a simple game.
int main() {

    const int winsNeeded = 4;
    const int maxRounds = 2 * winsNeeded - 1;

    // No winner yet.
    int winner = -1;
    cout << "Best " << winsNeeded << " out of " << maxRounds << ".  Go! " << string(40, '-') << "\n";

    // The game only supports two players.
    vector<unique_ptr<PlayerBase>> players;
    players.push_back(unique_ptr<PlayerBase> (new ComputerPlayer()));
    players.push_back(unique_ptr<PlayerBase> (new ComputerPlayer()));

    vector<int> victories(players.size());

    while (winner == -1) {
        cout << "\n";

        string a = players[0]->play();
        string b = players[1]->play();

        cout << ">> " << players[0]->name() << ": " << a << "\n";
        cout << ">> " << players[1]->name() << ": " << b << "\n";

        bool firstPlayerDefeatedSecondPlayer = defeats(a, b, true);
        cout << ".";
        if (a != b) {
            cout << "  "
                 << (firstPlayerDefeatedSecondPlayer ? players[0]->name() : players[1]->name())
                 << " wins.";
        }
        cout << "\n";

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
            cout << "The score is " << victories[0] << " ("
                 << players[0]->name() << ") to " << victories[1] << " ("
                 << players[1]->name() << ").\n";

            if (victories[roundVictorIndex] >= winsNeeded) {
                winner = roundVictorIndex;
            }
        }

    } // end (while no one has yet won the game)

    cout << "\n" << players[winner]->name() << " wins the game! " << string(40, '-') << "\n";
}
