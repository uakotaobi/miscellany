#include <iostream>
#include <algorithm>  // transform()
#include <limits>     // numeric_limits<streamsize>
#include <locale>     // tolower()
#include "HumanPlayer.h"

using std::numeric_limits;
using std::streamsize;
using std::transform;
using std::locale;
using std::string;
using std::cout;
using std::cin;

HumanPlayer::HumanPlayer() : name_("Human") { }
HumanPlayer::~HumanPlayer() { }
string HumanPlayer::name() const { return name_; }

string HumanPlayer::play() {

    string input;
    bool done = false;

    while (!done) {
        cout << "(R)ock, (P)aper, or (S)cissors? ";
        cin >> input;
        if (!cin) {
            cout << "Invalid choice.\n";

            // Clear any error flags
            cin.clear();

            // Ignore bad input up to the next newline
            cin.ignore(numeric_limits<streamsize>::max());

        } else {

            // Convert to lowercase.
            locale currentLocale("");
            transform(input.begin(), input.end(), input.begin(), [&currentLocale](char c) { return tolower(c, currentLocale); });

            if (input.length() == 1) {

                // Convertible to a single letter?
                switch(input[0]) {
                    case 'r':
                        input = "rock";
                        done = true;
                        break;
                    case 'p':
                        input = "paper";
                        done = true;
                        break;
                    case 's':
                        input = "scissors";
                        done = true;
                        break;
                    default:
                        cout << "Invalid choice \'" + input + "\' (Please enter 'R', 'P', or 'S'.)\n";
                }

            } else if (valid(input)) {

                // The user entered one of the three valid choices (in full.)
                done = true;

            } else {

                cout << "Invalid choice \"" + input + "\".\n";

            }
        }

    } // end (while not done)

    return input;
}
