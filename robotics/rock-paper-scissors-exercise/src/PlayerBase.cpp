#include <map>
#include <locale>
#include <iostream>
#include <stdexcept>
#include "PlayerBase.h"

using std::map;
using std::cout;
using std::string;
using std::locale;
using std::tolower;
using std::invalid_argument;

PlayerBase::~PlayerBase() { }

void PlayerBase::remember(const string& myLastPlay, const string& theirLastPlay) { }


// Non-public utility function.  Converts the given string to lowercase using
// the current locale's casing rules.
string lowercase(const std::string& s) {

    locale currentLocale("");
    string result = s;

    for (auto iter = result.begin(); iter != result.end(); ++iter) {
        *iter = tolower(*iter, currentLocale);
    }
    return result;

}

// Is this string a valid play for the current rock-paper-scissors game?
bool valid(const string& play) {

    string s = lowercase(play);
    if (s == "rock" || s == "paper" || s == "scissors") {
        return true;
    }
    return false;
}


// Returns true if, and only if, the first play beats the second.
//
// An invalid argument (like "nuke") throws an exception.
//
// Returns false otherwise.
bool defeats(const string& first, const string& second, bool print) {

    if (!valid(first)) {
        throw invalid_argument("defeats(): first argument (\"" + first + "\") is not valid.");
    }

    if (!valid(second)) {
        throw invalid_argument("defeats(): second argument (\"" + second + "\") is not valid.");
    }

    map<string, int> indexTable;
    indexTable["rock"] = 0;
    indexTable["paper"] = 1;
    indexTable["scissors"] = 2;

    string messageTable[3][3] = {

        // Rock [second]          Paper [second]            Scissors [second]
        { "Tie",                  "",                       "Rock breaks scissors" }, // Rock [first]
        { "Paper covers rock",    "Tie",                    "" },                     // Paper [first]
        { "",                     "Scissors cuts paper",    "Tie" }                   // Scissors [first]
    };

    unsigned int firstIndex = indexTable[first];
    unsigned int secondIndex = indexTable[second];
    string message = messageTable[firstIndex][secondIndex];

    if (message == "") {

        // First didn't beat second, so print the message that would have been
        // printed if the tables were turned.
        if (print) {
            cout << messageTable[secondIndex][firstIndex];
        }
        return false;

    } else if (message == "Tie") {

        if (print) {
            cout << message;
        }

        // If first tied second, then first did not beat second.
        return false;

    }

    // First beat second.
    if (print) {
        cout << message;
    }
    return true;
}
