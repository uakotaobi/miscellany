// #include <curses.h>
#include <SDL.h>
#include <thread>
#include <chrono>
#include <algorithm>
#include <array>
#include <vector>
#include <iostream>
#include <string>
#include <cmath>
#include <cstdint>  // uint8_t, uint16_t, uint32_t

using std::vector;
using std::max;
using namespace std;
using std::min;
using std::abs;
using std::sin;
using std::cos;

void pset(SDL_Surface* surface, int x, int y, char c, int color = 1) {
    int scrWidth = surface->w;
    int scrHeight = surface->h;

    if (x < 0 || x >= scrWidth || y < 0 || y >= scrHeight) {
        return;
    }
    int offset = scrWidth * y + x;

    uint8_t r, g, b;
    switch (color) {
        case 0: r = 0x00; g = 0x00; b = 0x00; break;
        case 1: r = 0x00; g = 0x00; b = 0xAA; break;
        case 2: r = 0x00; g = 0xAA; b = 0x00; break;
        case 3: r = 0x00; g = 0xAA; b = 0xAA; break;
        case 4: r = 0xAA; g = 0x00; b = 0x00; break;
        case 5: r = 0xAA; g = 0x00; b = 0xAA; break;
        case 6: r = 0xAA; g = 0xAA; b = 0x00; break;
        case 7: r = 0xAA; g = 0xAA; b = 0xAA; break;
    }
    uint32_t packed_pixel = SDL_MapRGB(surface-format, r, g, b);

    switch (surface->format->BytesPerPixel) {
        case 1: {
            uint8_t* target_pixel = reinterpret_cast<uint8_t*>(surface->pixels);
            target_pixel[offset] = static_cast<uint8_t>(packed_pixel);
            break;
        }
        case 2: {
            uint16_t* target_pixel = reinterpret_cast<uint16_t*>(surface->pixels);
            target_pixel[offset] = static_cast<uint16_t>(packed_pixel);
            break;
        }
        case 3: {
            uint32_t* target_pixel = reinterpret_cast<uint32_t*>(surface->pixels);
            target_pixel[offset] = target_pixel[offset] static_cast<uint16_t>(packed_pixel);

        }
    }
    //put a character at x y
    wmove(stdscr, y, x);
    waddch(stdscr, c | COLOR_PAIR(color));

    //need to do wrefresh() to show new info
    //wrefresh(stdscr);

}

void line(SDL_Surface* surface, int x1, int y1, int x2, int y2, char c, int color = 1) {
    double px = x1;
    double py = y1;
    double vx = x2 - x1;
    double vy = y2 - y1;
    double divisions = max(abs(vx), abs(vy));

    for (int i = 0; i < divisions; i++) {
        pset(surface, px, py, c, color);
        px += vx / divisions;
        py += vy / divisions;
    }

}

struct Point {
        double x;
        double y;
        double z;
};

const double length = 10;
vector<Point> vertices = {
    {-length,  length, length},
    { length,  length, length},
    { length, -length, length},
    {-length, -length, length},

    {-length,  length, -length},
    { length,  length, -length},
    { length, -length, -length},
    {-length, -length, -length},
};

vector<vector<int>> faces = {
    {0,1,2,3},
    {4,5,6,7},

    {0,1,5,4},
    {1,2,6,5},

    {2,3,7,6},
    {3,0,4,7}
};


// Multiplies a 4x4 affine matrix by a 3D vector, yielding another vector in
// return.
//
// Arguments:
// - matrix: A 4x4 transformation matrix with 16 elements in row-major order.
// - v:      A 1x4 vector with 4 elements in row-major order.
// Returns:
//   Returns a 1x4 transformed vector.
array<double, 4> mult(const array<double, 16>& matrix,
                      const array<double, 4>& v) {
    array<double, 4> product = {
        v[0] * matrix[0] + v[1] * matrix[1] + v[2] * matrix[2] + v[3] * matrix[3],
        v[0] * matrix[4] + v[1] * matrix[5] + v[2] * matrix[6] + v[3] * matrix[7],
        v[0] * matrix[8] + v[1] * matrix[9] + v[2] * matrix[10] + v[3] * matrix[11],
        v[0] * matrix[12] + v[1] * matrix[13] + v[2] * matrix[14] + v[3] * matrix[15],
    };
    return product;
}

array<double, 16> mult(const array<double, 16>& a,
                       const array<double, 16>& b) {
    return {
        // First row
        a[0] * b[0] + a[1] * b[4] + a[2] * b[8]  + a[3] * b[12],
        a[0] * b[1] + a[1] * b[5] + a[2] * b[9]  + a[3] * b[13],
        a[0] * b[2] + a[1] * b[6] + a[2] * b[10] + a[3] * b[14],
        a[0] * b[3] + a[1] * b[7] + a[2] * b[11] + a[3] * b[15],

        // Second row
        a[4] * b[0] + a[5] * b[4] + a[6] * b[8]  + a[7] * b[12],
        a[4] * b[1] + a[5] * b[5] + a[6] * b[9]  + a[7] * b[13],
        a[4] * b[2] + a[5] * b[6] + a[6] * b[10] + a[7] * b[14],
        a[4] * b[3] + a[5] * b[7] + a[6] * b[11] + a[7] * b[15],

        // Third row
        a[8] * b[0] + a[9] * b[4] + a[10] * b[8]  + a[11] * b[12],
        a[8] * b[1] + a[9] * b[5] + a[10] * b[9]  + a[11] * b[13],
        a[8] * b[2] + a[9] * b[6] + a[10] * b[10] + a[11] * b[14],
        a[8] * b[3] + a[9] * b[7] + a[10] * b[11] + a[11] * b[15],

        // Fourth row
        a[12] * b[0] + a[13] * b[4] + a[14] * b[8]  + a[15] * b[12],
        a[12] * b[1] + a[13] * b[5] + a[14] * b[9]  + a[15] * b[13],
        a[12] * b[2] + a[13] * b[6] + a[14] * b[10] + a[15] * b[14],
        a[12] * b[3] + a[13] * b[7] + a[14] * b[11] + a[15] * b[15],
    };
}

// Returns the X rotation matrix: A matrix that rotates an arbitrary point
// theta radians around the positive X axis.
array<double, 16> xRotate(double theta) {
    array<double, 16> matrix = {
        1, 0,          0,           0,
        0, cos(theta), -sin(theta), 0,
        0, sin(theta), cos(theta),  0,
        0, 0,          0,           1
    };
    return matrix;
}

// Returns the Y rotation matrix: A matrix that rotates an arbitrary point
// theta radians around the positive Y axis.
array<double, 16> yRotate(double theta) {
    array<double, 16> matrix = {
        cos(theta),  0, sin(theta), 0,
        0,           1, 0,          0,
        -sin(theta), 0, cos(theta), 0,
        0,           0, 0,          1
    };
    return matrix;
}

// Returns the Z rotation matrix: A matrix that rotates an arbitrary point
// theta radians around the positive Z axis.
array<double, 16> zRotate(double theta) {
    array<double, 16> matrix = {
        cos(theta), -sin(theta), 0, 0,
        sin(theta), cos(theta),  0, 0,
        0,          0,           1, 0,
        0,          0,           0, 1
    };
    return matrix;
}

void polyhedron(SDL_Surface* surface, const vector<Point>& vertices, const vector<vector<int>>& faces)  {
    int scrWidth = 0;
    int scrHeight = 0;
    getmaxyx(stdscr, scrHeight, scrWidth);


    const double d = 30; // Distance from camera to eyeball before projection

    for (unsigned int i = 0; i < faces.size(); i++) {
        vector<Point> currentFace;
        for (unsigned int j = 0; j < faces[i].size(); j++) {
            currentFace.push_back(vertices[faces[i][j]]);
        }
        for (unsigned int j = 0; j < currentFace.size(); j++) {
            Point currentVertex = currentFace[j];
            Point nextVertex = currentFace[(j + 1) % currentFace.size()];

            //use projection formula
            double x1 = (d * currentVertex.x) / (d + currentVertex.z);
            double y1 = (d * currentVertex.y) / (d + currentVertex.z);
            double x2 = (d * nextVertex.x) / (d + nextVertex.z);
            double y2 = (d * nextVertex.y) / (d + nextVertex.z);

            //centering
            x1 += scrWidth / 2;
            y1 += scrHeight / 2;
            x2 += scrWidth / 2;
            y2 += scrHeight / 2;

            line(surface, x1, y1, x2, y2, '*', j + 1);
            //line(121,30,121,5,'@',2);
            //cout << "Line: " << x1 << ", " << y1 << ", " << x2 << ", " << y2 << ".\n";

        }
    }
    //line(121,15,121,5,'@',2);
}

int main(int argc, const char* argv[]) {

    double thetaX, thetaY, thetaZ = 0; // In degrees per frame
    const double pi = 3.141593;
    const double DEGREES_TO_RADIANS = pi / 180.0;

    switch(argc) {
        case 4:
            // Passed in x, y, and z
            thetaZ = stod(argv[3]) * DEGREES_TO_RADIANS;
            // FALL THROUGH
        case 3:
            // Passed in x and y
            thetaY = stod(argv[2]) * DEGREES_TO_RADIANS;
            // FALL THROUGH
        case 2:
            // Passed in X
            thetaX = stod(argv[1]) * DEGREES_TO_RADIANS;
            // FALL THROUGH
        case 1:
            // No args were passed in.
            // Not an error.
            break;
    }

    // Initialize SDL.
    if (SDL_Init(SDL_INIT_EVERYTHING) != 0) {
        SDL_Log("Unable to initialize SDL: %s", SDL_GetError());
        return 1;
    }

    // Initialize an SDL window.
    SDL_Window* window = nullptr;
    SDL_Renderer* renderer = nullptr;
    if (SDL_CreateWindowAndRenderer(1280,                   // width, in pixels
                                    960,                    // height, in pixels
                                    SDL_WINDOW_RESIZABLE,
                                    &window,
                                    &renderer) != 0) {

    }

    if (window == nullptr) {
        // The window could not be spawned.
        SDL_LogError(SDL_LOG_CATEGORY_APPLICATION, "Could not create window: %s\n", SDL_GetError());
        return 1;
    }

    // Initialize the surface we will render on.
    SDL_Surface* screen = SDL_GetWindowSurface(window);
    if (screen == nullptr) {
        // Could not obtain a surface from the window (this is unusual.)
        SDL_LogError(SDL_LOG_CATEGORY_APPLICATION, "Could not create surface: %s\n", SDL_GetError());
        return 1;
    }

    //line(1,1,14,12,'*',3);

    //pset(1,1,'^');
    //xpset(14,12,'^');

    // polyhedron(vertices, faces);
    // wrefresh(stdscr);

    bool done = false;
    while (!done) {

        SDL_Event event;
        while (SDL_PollEvent(&event)) {
            // Keep eating the events until there's nothing left to munch on
            // in the event queue.
            switch (event.type) {
                case SDL_KEYDOWN:
                    if (event.keysym.sym == "Q") {
                        done = true;
                    }
                    break;
                case SDL_KEYUP:
                    break;
                case SDL_Quit:
                    // Window was closed.
                    done = true;
                    break;
                default:
                    // Ignore all other events...for now.
                    break;
            }
        }

        for (unsigned int i = 0; i < vertices.size(); i++) {
            array<double, 16> xmatrix = xRotate(thetaX);
            array<double, 16> ymatrix = yRotate(thetaY);
            array<double, 16> zmatrix = zRotate(thetaZ);
            array<double, 16> matrix = mult(zmatrix, ymatrix);
            matrix                   = mult(matrix,  xmatrix);
            array<double, 4> v = {vertices[i].x, vertices[i].y, vertices[i].z, 1};
            array<double, 4> result = mult(matrix, v);

            Point newPos = {result[0], result[1], result[2]};
            vertices[i] = newPos;
        }

        // We're about to draw.  Lock the surface so we have access to the
        // surface's pixel buffer.
        if (SDL_MUSTLOCK(screen)) {
            int error = SDL_LockSurface(screen);
            if (error != 0) {
                SDL_Log("Could not lock surface: %s\n", SDL_GetError());
                return 1;
            }
        }

        // Draw!
        polyhedron(screen, vertices, faces);

        // Make the changes to the surface visible.
        SDL_UpdateWindowSurface(window);
        if (SDL_MUSTLOCK(screen)) {
            int error = SDL_UnlockSurface(screen);
            if (error != 0) {
                SDL_Log("Could not unlock surface: %s\n", SDL_GetError());
                return 1;
            }
        }

        //need to sleep
        std::this_thread::sleep_for(std::chrono::milliseconds(10));
    }

    // Do not destroy the window's surface (screen); SDL_DestroyWindow will
    // accomplish that on its own.
    SDL_DestroyWindow(window);
    SDL_Quit();

    for (unsigned i = 0; i < vertices.size(); i++) {
        cout << "Vertices[" << i << "] = " << vertices[i].x << ", " << vertices[i].y << ", " << vertices[i].z << "\n";
    }

}
