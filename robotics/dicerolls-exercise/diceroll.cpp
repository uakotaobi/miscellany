// -*- mode: C++; compile-command: "g++ -Wall -g --std=c++14 -o diceroll diceroll.cpp" -*-
//
// By Uche A., 2017-07-23.  Public domain.
//
// Now available as a gist: https://gist.github.com/uakotaobi/ff798959c162aeb672afb62abf5474f8

#include <random>     // default_random_engine, uniform_int_distribution<T>
#include <chrono>     // high_resolution_clock
#include <vector>     // vector<T>
#include <iterator>   // ostream_iterator<T>
#include <iostream>   // cin, cout
#include <algorithm>  // generate(), accumulate(), copy()
#include <functional> // bind()

using namespace std;

int main() {
    int dice = 0;
    do {
        cout << "Number of dice to roll? ";
        cin >> dice;
        // Reset std::cin in case of bad integer input.
        cin.clear();
        cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
    } while (dice <= 0);

    vector<int> rollsToAdd(dice);

    auto clockTicks = chrono::high_resolution_clock::now().time_since_epoch().count();
    default_random_engine generator(clockTicks);
    uniform_int_distribution<int> distribution(1, 6);

    // Calling rollOneDie() with no args is the same as calling
    // distribution(generator).
    auto rollOneDie = bind(distribution, generator);

    // Call rollOneDie() on every element.
    generate(rollsToAdd.begin(), rollsToAdd.end(), rollOneDie);

    int sum = accumulate(rollsToAdd.begin(), rollsToAdd.end(), 0);

    if (dice > 1) {
        cout << "Sum of " << rollsToAdd.size() << " dice rolls: ";
        copy(rollsToAdd.begin(), rollsToAdd.end() - 1, ostream_iterator<int>(cout, " + "));
        cout << rollsToAdd.back() << " = ";
    } else {
        cout << "You rolled a ";
    }
    cout << sum << "\n";
}
