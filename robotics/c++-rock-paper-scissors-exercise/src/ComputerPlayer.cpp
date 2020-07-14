#include <sstream>
#include <algorithm> // copy, copy_n
#include <iostream>  // cout, fixed
#include <random>    // default_random_engine, uniform_int_distribution
#include <chrono>    // chrono::high_resolution_clock
#include <deque>
#include <fstream>   // ofstream, ifstream
#include <map>
#include <vector>
#include <iterator>  // back_inserter
#include "ComputerPlayer.h"

using std::chrono::high_resolution_clock;
using std::uniform_int_distribution;
using std::default_random_engine;
using std::istream_iterator;
using std::ostream_iterator;
using std::back_inserter;
using std::ostringstream;
using std::istringstream;
using std::ofstream;
using std::ifstream;
using std::getline;
using std::string;
using std::vector;
using std::deque;
using std::fixed;
using std::cout;
using std::map;
using std::min;

int ComputerPlayer::counter = 0;
const string MOVES_DB_FILE = "./.moves";
const unsigned int LOOKAHEAD_SIZE = 8;

class MoveDatabase {
    public:
        MoveDatabase(unsigned int lookaheadLimit);
        void remember(const string& humanPlay, const string& computerPlay);
        string predict() const;
        void save(const string& filename) const;
        void load(const string& filename);
    private:
        size_t limit;
        deque<string> mostRecentHumanMoves;

        // The keys for this map are vectors of strings representing the N
        // most human move, where N is the lookahead limit.
        //
        // The values for this map are simply a map counting the number of
        // times that the human, in the given situation, gave a particular
        // response.
        //
        // For example, if the lookahead were 3, the last human moves were
        // "rock, rock, scissors", and the current human move was "paper",
        // then frequencies[{ "rock", "rock", "scissors"}]["paper"] would be
        // incremented by 1.

        map<vector<string>, map<string, int>> frequencies;

};

ComputerPlayer::ComputerPlayer() : id_(++counter), name_() {
    ostringstream stream;
    stream << "CPU #" << id_;
    name_ = stream.str();
}

std::string ComputerPlayer::name() const { return name_; }

string ComputerPlayer::play() {

    MoveDatabase movedb(LOOKAHEAD_SIZE);
    string prediction = movedb.predict();
    if (prediction != "") {

        // If we have enough data to make a prediction, then do it.
        return prediction;

    } else {

        // Make a no-so-sophisticated, random play.
        default_random_engine generator(high_resolution_clock::now().time_since_epoch().count());
        uniform_int_distribution distribution(0, 2);
        const char* choices[] = { "rock", "paper", "scissors" };
        return choices[distribution(generator)];
    }
}

void ComputerPlayer::remember(const string& myLastPlay, const string& theirLastPlay) {


    // Doesn't do anything right now.
    MoveDatabase movedb(LOOKAHEAD_SIZE);
    movedb.remember(theirLastPlay, myLastPlay);
    movedb.save(MOVES_DB_FILE);
}


MoveDatabase::MoveDatabase(unsigned int lookaheadLimit) : limit(lookaheadLimit) {
    // Read in the moves database if we can.
    ifstream file(MOVES_DB_FILE);
    if (file) {
        load(MOVES_DB_FILE);
    }
}

void MoveDatabase::save(const string& filename) const {

    ofstream out(filename);

    // File header.
    out << "MVDB\n";

    // Lookahead size.
    out << limit << "\n";

    // Sequence of most recently-seen human moves.
    copy(mostRecentHumanMoves.begin(),
         mostRecentHumanMoves.end(),
         ostream_iterator<string>(out, " "));
    out << "\n";

    // The actual frequency database.  Each record is separated with a lone
    // asterisk.
    for (const auto& p : frequencies) {
        const vector<string>& sequence = p.first;
        const map<string, int>& table = p.second;
        copy(sequence.begin(),
             sequence.end(),
             ostream_iterator<string>(out, " "));
        out << "\n";
        for (const auto& p2 : table) {
            out << p2.first << " = " << p2.second << "\n";
        }
        out << "*\n";
    }
}

void MoveDatabase::load(const string& filename) {
    ifstream in(filename);

    // Read the file header.
    string header;
    in >> header;

    // Read in the lookahead size.
    // If ours is bigger, we ignore what we read.
    size_t limit_from_file;
    in >> limit_from_file;
    if (limit_from_file > limit) {
        limit = limit_from_file;
    }
    in.get();

    // Read in the sequence of most recently-seen human moves.
    string line;
    getline(in, line, '\n');
    istringstream stream(line);
    while (stream) {
        string move;
        if (stream >> move) {
            mostRecentHumanMoves.push_back(move);
        }
    }

    // Read in the frequency database.
    vector<string> key;
    while(in) {
        getline(in, line);
        if (line == "*") {
            key.clear();
        } else if (key.size() == 0) {
            // Read and parse the key itself.
            istringstream stream(line);
            istream_iterator<string> iter(stream), end;
            copy(iter, end, back_inserter(key));

            frequencies[key] = map<string, int>();
        } else {
            // Read the move and the count.
            istringstream stream(line);
            string humanMove;
            int count;
            char equals;
            stream >> humanMove >> equals >> count;

            frequencies[key][humanMove] = count;
        }
    }

    // Summarize:
    // cout << ">>> Read successfully.\n";
    // cout << ">>> Lookahead limit: " << limit << "\n";
    // cout << ">>> Most recent human moves: [";
    // copy(mostRecentHumanMoves.begin(),
    //      mostRecentHumanMoves.end(),
    //      ostream_iterator<string>(cout, " "));
    // cout << "]\n";
}

void MoveDatabase::remember(const string& humanPlay, const string&) {

    // Update the frequency database.  We want to know, for each sequence of
    // moves up to the lookahead limit, what the human entered under those
    // circumstances.
    for (unsigned int i = 1; i <= min(limit, mostRecentHumanMoves.size()); ++i) {
        vector<string> last_n_moves;
        copy_n(mostRecentHumanMoves.end() - i, i, back_inserter(last_n_moves));

        if (frequencies.find(last_n_moves) == frequencies.end()) {
            // This sequence has not been seen.
            frequencies[last_n_moves] = map<string, int>();
            // cout << "Allocating: frequencies[";
            // copy(last_n_moves.begin(),
            //      last_n_moves.end(),
            //      ostream_iterator<string>(cout, " "));
            // cout << "] = { }\n";
        }

        map<string, int>& table = frequencies[last_n_moves];
        if (table.find(humanPlay) == table.end()) {
            // This play has not been seen for this sequence.
            table[humanPlay] = 0;
            // cout << "Initializing: frequencies[";
            // copy(last_n_moves.begin(),
            //      last_n_moves.end(),
            //      ostream_iterator<string>(cout, " "));
            // cout << "][" << humanPlay << "] = 0\n";
        }

        // Increment existing record.
        table[humanPlay] += 1;
        // cout << "Incrementing: frequencies[";
        // copy(last_n_moves.begin(),
        //      last_n_moves.end(),
        //      ostream_iterator<string>(cout, " "));
        // cout << "][" << humanPlay << "] = " << table[humanPlay] << "\n";
    }

    // Throw the most recent human moves into a queue.
    mostRecentHumanMoves.push_back(humanPlay);

    while (mostRecentHumanMoves.size() > limit) {
        mostRecentHumanMoves.pop_front();
    }
}

string MoveDatabase::predict() const {

    if (mostRecentHumanMoves.empty()) {
        cout << ">> Frequency analysis: No moves recorded yet.\n";
        return "";
    }

    map<string, int> table;
    vector<string> last_n_moves;

    for (size_t i = min(mostRecentHumanMoves.size(), limit); i >= 1; --i) {

        // Find the largest sequence of the human's most recent moves that we
        // can in the moves database.
        last_n_moves.clear();
        copy_n(mostRecentHumanMoves.end() - i, i, back_inserter(last_n_moves));

        const auto iter = frequencies.find(last_n_moves);
        if (iter != frequencies.end()) {
            table = iter->second;

            if (table.empty()) {
                // cout << ">> Frequency analysis: not enough data recorded for sequence [";
                copy(last_n_moves.begin(),
                     last_n_moves.end(),
                     ostream_iterator<string>(cout, " "));
                cout << "].\n";
                return "";
            }
            break;
        }

        if (i == 1) {
            // The database is just too empty to use for predictions right
            // now.
            // cout << ">> Frequency analysis: Cannot find sequence in database (not populated enough.)\n";
            return "";
        }
    }

    // Frequency analysis: which play has the human historically done the most
    // in this situation?
    int maximum = 0, sum = 0;
    string mostFrequent = "";
    for (const auto& p : table) {
        string humanPlay = p.first;
        int count = p.second;

        sum += count;

        if (count > maximum) {
            maximum = count;
            mostFrequent = humanPlay;
        }
    }

    string response;
    if (mostFrequent == "rock") {
        response = "paper";
    } else if (mostFrequent == "scissors") {
        response = "rock";
    } else if (mostFrequent == "paper") {
        response = "scissors";
    }

    // Uncommenting this line should frustrate the human by going for ties
    // instead of beating them!
    //
    // response = mostFrequent;

    // cout << ">> Frequency analysis: humans have most commonly chosen " << mostFrequent << " after [";
    // copy(last_n_moves.begin(),
    //      last_n_moves.end(),
    //      ostream_iterator<string>(cout, " "));
    // cout.precision(2);
    // cout << "] (" << fixed << (static_cast<float>(100 * maximum) / sum) << "%).  Choosing ";
    // cout << response << ".\n";
    return response;
}
