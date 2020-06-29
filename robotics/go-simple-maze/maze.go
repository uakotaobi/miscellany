package main
import (
	"os"
	"fmt"
	"math"
	"math/rand"
	"time"
	"unicode/utf8"
	"github.com/akamensky/argparse"
)

type Maze struct {
	width int
	height int
	cells []rune
	thickness int
	intersection rune
	horizontal rune
	vertical rune
	verbosity int
	floor rune
	fill rune
	minWallLength, maxWallLength int
	entrance, exit struct{x, y, width, height int}
}

// Constants used for neighbor specification.  For instance, "every neighbor
// except the top" would be represented as (left | right | down).
type mask int
const (
	left mask = 1 << iota
	up
	right
	down
)

var directions []struct{x, y int} = []struct{x, y int} {
	{-1, 0}, // Left vector
	{0, -1}, // Up vector
	{1, 0},  // Right vector
	{0, 1},  // Down vector
}

func NewMaze(width, height int) Maze {
	m := Maze{
		thickness: 1,
		intersection: '+',
		horizontal: '-',
		vertical: '|',
		floor: ' ',
		fill: '.',
		minWallLength: 3,
		maxWallLength: math.MaxInt64,
	}
	m.setSize(width, height)

	m.Clear()
	return m
}

// Uses the empty areas of an existing maze as the basis for this new one.
//
// The idea here is that we can generate a large maze using thickness>1, then
// generate a small maze on top of that in order to create a sort of
// "super-maze."
//
//  To generate a singular maze, use NewMaze().
//
// Note: If you are going to change the display runes, you should probably
// make sure m.floor continues to match other.floor.
func NewMazeOverExisting(other Maze) Maze {

	// Copy the other maze's parameters by value: same width, height,
	// display runes....
	m := other

	// But take care of the reference member (the cells slice) using a
	// deep copy.
	m.cells = make([]rune, m.width * m.height)
	copy(m.cells, other.cells)

	return m
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Alters the maze's dimensions.  The maze will need re-rendering after the
// call.
func (m *Maze) setSize(newWidth, newHeight int) {
	m.width = newWidth;
	m.height = newHeight;

	// TODO: What are we going to do about m.cells?  The right thing would
	// be to copy the existing cells into the upper-left corner of the
	// resized slice (or to truncate, if the new size was smaller.)
}

// Erases the contents of the maze, overwriting it with the m.floor rune.
// You should call this if m.floor changes.
func (m *Maze) Clear() {
	m.cells = make([]rune, m.width * m.height)
	for index := range(m.cells) {
		m.cells[index] = m.floor
	}
}

func (m *Maze) SetWidth(newWidth int) {
	m.setSize(newWidth, m.height)
	m.Generate()
}
func (m *Maze) SetHeight(newHeight int) {
	m.setSize(m.width, newHeight)
	m.Generate()
}

// Gets the width of the maze in cells.
func (m *Maze) Width() int { return m.width; }

// Gets the height of the maze in cells.
func (m *Maze) Height() int { return m.height; }

// Utility function for calculating the offset of a cell within the cells
// array for a given coordinate.  Note that this function does not perform
// out-of-bounds checking.
func (m *Maze) offset(x, y int) int {
	if m != nil {
		return m.width * y + x
	}
	return 0
}

// A simple utility function for determining if the given coordinate is in-bounds.
func (m *Maze) valid(x, y int) bool {
	if m == nil {
		return false
	}
	return (x >= 0 && y >= 0 && x < m.width && y < m.height)
}

// Utility function for drawRect().  Determines if the given coordinate has
// orthogonally adjacent neighbors whose cells all contain the given value.
//
// The neighbor positions are specified by passing the appropriate bitwise
// combination of the left, right, up and down constants into the mask
// argument.
//
// Note that a mask of 0 will always return true.
func (m *Maze) hasNeighbor(x, y int, mask mask, neighborValue rune) bool {
	if m == nil {
		// Missing receiver.
		return false
	}
	if (mask & left) != 0 && m.valid(x - 1, y) && m.cells[m.offset(x - 1, y)] != neighborValue {
		return false
	}
	if (mask & up) != 0 && m.valid(x, y - 1) && m.cells[m.offset(x, y - 1)] != neighborValue {
		return false
	}
	if (mask & right) != 0 && m.valid(x + 1, y) && m.cells[m.offset(x + 1, y)] != neighborValue {
		return false
	}
	if (mask & down) != 0 && m.valid(x, y + 1) && m.cells[m.offset(x, y + 1)] != neighborValue {
		return false
	}
	return true
}

// Utility function for drawRect().  Determines if the rectangle with the
// given dimensions and upper-left corner are within the bounds of the maze.
//
// Returns true if at least part of the rectangle is in-bounds.  Call clip()
// in order to force the rectangle to be entirely in-bounds.
func (m *Maze) validRect(x, y, width, height int) bool {
	if m == nil {
		// Missing receiver.
		return false
	}

	if width <= 0 || height <= 0 {
		// Degenerate rectangle.
		return false
	}

	// Rectangle is off-screen?
	if x + width <= 0 || x >= m.width || y + height <= 0 || y >= m.height {
		return false
	}

	return true
}

// Utility function for drawRect().  Adjusts the boundaries of the given
// rectangle boundary so it lies entirely in-bounds.
//
// Takes the upper-left corner and dimensions of the input rectangle as
// arguments; the output rectangle is returned as a tuple (newX, newY,
// newWidth, newHeight.)
//
// Warning: If given a rectangle that is invalid (i.e., one for which
// validRect() is false), the rectangle you get out of this function will also
// be invalid.
func (m *Maze) clipRect(x, y, width, height int) (int, int, int, int) {

	if m == nil {
		// Missing receiver.  Clipping does nothing.
		return x, y, width, height
	}

	// Clip the rectangle so that we only draw the portion of it that is
	// on-screen.
	newX, newWidth := max(0, x), min(width, x + width)    // Handle x < 0
	newY, newHeight := max(0, y), min(height, y + height) // Handle y < 0
	newWidth = min(newWidth, m.width - x)                 // Handle x + width > m.width
	newHeight = min(newHeight, m.height - y)              // Handle x + height > m.height

	return newX, newY, newWidth, newHeight
}

// Utility function for generateMaze().  Draws a rectangle of the given
// dimensions in the maze, using the given coordinates as the rectangle's
// upper left corner.
//
// Edge cases:
//
//   1. Drawing a rectangle with a minimum dimension equal to 1 will draw a
//      horizontal line if the width is greater than 1, a vertical line if the
//      height is greater than 1, and an intersection character if the rectangle
//      is 1x1.
//   2. Drawing a rectangle with a minimum dimension greater than 2 will
//      utilize the given fill character argument for the rectangle interior.
//   3. A rectangle that is adjacent to unoccupied cells will be
//      drawn with solid edges; if any of the edge-adjacent cells are
//      occupied, the edges between will be drawn so as to connect the two
//      cells.
//   4. Drawing out of bounds is harmless.
func (m *Maze) drawRect(x, y, width, height int, fill rune) {

	if m == nil {
		// Missing receiver; someone called drawRect() on a nil Maze
		// instance.
		return
	}

	if !m.validRect(x, y, width, height) {
		// Rectangle is degenerate or entirely out of bounds.
		return
	}

	// Ensure that the rectangle is entirely in-bounds.
	x, y, width, height = m.clipRect(x, y, width, height)

	for row := y; row < y + height; row++ {
		for column := x; column < x + width; column++ {

			// Find the nominal value the rectangle is supposed to
			// have at this position.
			var proposedCell rune
			switch {
			case (column == x || column == x + width - 1) && (row == y || row == y + height - 1):
				proposedCell = m.intersection
				if (height == 1 || width == 1) && m.cells[m.offset(column, row)] == m.floor {
					// Minor optimization for walls of
					// thickness 1.  This end of the line
					// isn't touching anything, so replace
					// the intersection with a
					// better-looking rune.
					if height == 1 {
						proposedCell = m.horizontal
					} else {
						proposedCell = m.vertical
					}
				}
			case column == x:
				proposedCell = m.vertical
				if !m.hasNeighbor(column, row, left, m.floor) {
					// Touching a non-floor to the left.
					// proposedCell = m.intersection
				}
			case column == x + width - 1:
				proposedCell = m.vertical
				if !m.hasNeighbor(column, row, right, m.floor) {
					// Touching a non-floor to the right.
					// proposedCell = m.intersection
				}
			case row == y:
				proposedCell = m.horizontal
				if !m.hasNeighbor(column, row, up, m.floor) {
					// Touching a non-floor to the top.
					// proposedCell = m.intersection
				}
			case  row == y + height - 1:
				proposedCell = m.horizontal
				if !m.hasNeighbor(column, row, down, m.floor) {
					// Touching a non-floor to the bottom.
					// proposedCell = m.intersection
				}
			default:
				// This is the interior of the rectangle.
				proposedCell = fill
			}

			m.cells[m.offset(column, row)] = proposedCell
		}
	}
}

// Returns true if the given rectangle (or at least the portions of it that
// are on-screen) contains only cells with the given content.
//
// A rectangle that is out of bounds will return false.
func (m *Maze) rectContains(x, y, width, height int, cell rune) bool {
	if !m.validRect(x, y, width, height) {
		return false
	}
	x, y, width, height = m.clipRect(x, y, width, height)
	for row := y; row < y + height; row++ {
		for column := x; column < x + width; column++ {
			if m.cells[m.offset(column, row)] != cell {
				return false
			}
		}
	}
	return true
}

// Returns true if the given rectangle represents a passageway: its center (if
// it has one) solely consists of floor runes, and walls do not hem the
// center in on all four sides.
func (m *Maze) rectIsPassage(x, y, width, height int) bool {
	if !m.validRect(x, y, width, height) {
		return false
	}
	x, y, width, height = m.clipRect(x, y, width, height)

	wallOpening := false
	for row := y; row < y + height; row++ {
		for column := x; column < x + width; column++ {
			c := m.cells[m.offset(column, row)]

			if column > x && column < x + width - 1 && row > y && row < y + height - 1 {
				// Interior cell.
				if c != m.floor {
					return false
				}
			} else {
				// Border cell.
				if c == m.floor {
					wallOpening = true
				}
			}
		}
	}
	return wallOpening
}

// Helpful conversion routines.
//
// Converts a "unit coordinate" (a subdivision of the maze grid into
// slightly-overlapping boxes) into a square.
func (m *Maze) unitCoordinatesToRect(unitColumn, unitRow int) (x, y, width, height int) {
	switch m.thickness {
	case 1:
		return unitColumn, unitRow, 1, 1
	case 2:
		// Notice how we special-case this so that the squares *don't*
		// overlap.  (If they did, the corridors would have a width of
		// 0.)  This mandates slight adjustments in a few other parts
		// of the code.
		return unitColumn*2, unitRow*2, 2, 2
	default:
		x = unitColumn * (m.thickness - 1)
		y = unitRow * (m.thickness - 1)
		return x, y, m.thickness, m.thickness
	}
}

// Utility function for findEntranceAndExit().  Provides all the points around
// the perimeter of the given rectangle.
//
// Arguments:
// - width: The width of the rectangle.
// - height: The height of the rectangle.
//
// Return value:
// - Returns an array of points that lie along the perimeter of the rectangle,
//   starting with the upper-left corner and proceeding clockwise.
//
// TODO: This is really better done as a generator function, or perhaps as an
// iterator.
func rectPerimeter(width, height int) []struct{x, y int} {
	result := []struct{x, y int}{}

	for i := 0; i < 2 * (width - 1 + height - 1); i++ {
		switch {
		case i < width - 1:
			// Top of rect (excluding upper-right corner),
			// heading right
			result = append(result, struct{x, y int}{
				x: i,
				y: 0,
			})
		case i < (width - 1 + height - 1):
			// Right side of rect (excluding bottom
			// right), heading down
			result = append(result, struct{x, y int}{
				x: width - 1,
				y: i - (width - 1),
			})
		case i < (2 * (width - 1) + height - 1):
			// Bottom of rect (excluding lower-left
			// corner), heading left
			result = append(result, struct{x, y int}{
				x: (width - 1) - (i - (width - 1) - (height - 1)),
				y: height - 1,
			})
		default:
			// Left side of rect (excluding top-left
			// corner), heading up
			result = append(result, struct{x, y int}{
				x: 0,
				y: (height - 1) - (i - (2 * (width - 1)) - (height - 1)),
			})
		}
	}
	return result
}


// Helper function for generateMaze().
//
// Let us define a unit rectangle as part of a maze's "outer corridor" if it
// is both completely empty and it borders a (wall-filled) unit rectangle on
// the edge of the maze.  This function's purpose is to find the coordinates
// of the two outer corridor unit rectangles that are furthest apart in a maze
// walk.  These will then be used to cut the entrance and exit for the maze.
//
// A single flood fill from any starting point will eventually find all the
// points that are furthest away from it.  This algorithm repeats that flood
// fill once for HALF of the outer corridor unit rectangles in the maze (since
// the other half will produce identical results.)  This is obviously a brute
// force algorithm, and it is most likely highly inefficient.
//
// Returns the unit coordinates of the entrance and exit, such that you can simply
// punch holes in those positions to complete the maze.
func (m *Maze) findEntranceAndExit(unitWidth, unitHeight int) (entranceUnitColumn, entranceUnitRow, exitUnitColumn, exitUnitRow, solutionLength int) {

	type Point struct {
		x, y int
	}

	// A data structure that records which points we have and have not
	// visited during recursion.
	type VisitedMap struct {
		visitedUnits map[Point]bool
		furthestPoints []Point
		longestDistance int
	}

	// Returns a list of outer corridor points that are the furthest
	// distance away from the given start point (each will have the same
	// distance from the start point.)
	//
	// Arguments:
	// - unitColumn:        The X-coordinate of the starting unit rectangle.
	// - unitRow:           The Y-coordinate of the starting unit rectangle.
	// - currentDistance:	The Manhattan distance, in unit coordinates,
	//			from our originating point when the recursion
	//			started.
	// - VisitedMap:        A data structure that contains:
	//   * visitedUnits:    An associative array mapping unit rectangle
	//                      coordinates to booleans.  It is true for all
	//                      points visited so far during maze recursion,
	//                      and false for all other points.
	//
	//   * furthestPoints:  A slice of unit positions that have the
	//                      furthest distance from the starting position
	//                      of the recursion.
	//
	//   * longestDistance: The value of that furthest distance.
	//
	//                      In general, the _shorter_ (yes, shorter) the
	//                      best distance is, the higher-quality the maze
	//                      is, since that means that more of it is
	//                      devoted to misleading branches and dead ends.
	//
	//                      The longer this value is, the more
	//                      labyrinthine the maze is, up to the point of
	//                      being _unicursal_ (that is, having a single,
	//                      winding path with no branches at all.)
	//                      Unicursal mazes take a long time to traverse,
	//                      but they are not challenging.
	var findFurthestOuterCorridorPositionRecursive func(unitColumn, unitRow, currentDistance int, visitedMap *VisitedMap)
	findFurthestOuterCorridorPositionRecursive = func(unitColumn, unitRow, currentDistance int, visitedMap *VisitedMap) {

		visitedMap.visitedUnits[Point{x: unitColumn, y: unitRow}] = true

		// Test all four neighbors in turn.
		for i := 0; i < 4; i++ {
			var dx, dy int = directions[i].x, directions[i].y

			neighbor := Point{x: unitColumn + dx, y: unitRow + dy}
			if neighbor.x < 1 || neighbor.x > unitWidth - 2 || neighbor.y < 1 || neighbor.y > unitHeight - 2 {
				// Neighbor position just hit a border
				// wall: out of bounds.
				continue
			}

			if visitedMap.visitedUnits[neighbor] {
				// Neighbor position was already
				// visited via recursion.
				continue
			}

			// We have an unvisited neighboring unit cell.
			x, y, width, height := m.unitCoordinatesToRect(neighbor.x, neighbor.y)
			if !m.rectIsPassage(x, y, width, height) {
				// The neighboring cell was occupied.
				continue
			}

			// We have an empty, unvisited neighboring
			// unit cell.  Flood recursively.
			findFurthestOuterCorridorPositionRecursive(unitColumn + dx,
				unitRow + dy,
				currentDistance + 1,
				visitedMap)

		} // end (for each orthogonally neighboring unit rectangle position)

		// If control makes it here, all neighbors are
		// visited (or occupied.)
		//
		// Add this point only if it qualifies.

		if (unitColumn != 1 && unitColumn != unitWidth - 2) || (unitRow != 1 && unitRow != unitHeight - 2) {
			// Not a unit coordinate on an outer corridor.
			return
		}

		if currentDistance > visitedMap.longestDistance {
			visitedMap.furthestPoints = []Point{Point{x: unitColumn, y: unitRow}}
			visitedMap.longestDistance = currentDistance
		} else if currentDistance == visitedMap.longestDistance {
			visitedMap.furthestPoints = append(visitedMap.furthestPoints, Point{x: unitColumn, y: unitRow})
		}
	}

	// Test the entire perimeter of the outer corridor of the maze,
	// looking in each case for the outer cooridor unit coordinate that is
	// the furthest distance away.
	//
	// The winners become the maze's entrance and exit.

	longestDistance := -1
	finalCandidates := []Point{}
	for _, p := range(rectPerimeter(unitWidth - 2, unitHeight - 2)) {
		var unitColumn, unitRow int = p.x + 1, p.y + 1

		x, y, width, height := m.unitCoordinatesToRect(unitColumn, unitRow)
		if !m.rectIsPassage(x, y, width, height) {
			if unitWidth > 3 || unitHeight > 3 {
				// Can't start entrances behind a wall.
				continue
			} else {
				// Degenerate maze is only 3x3 units.
				// The middle unit at (1, 1) is hemmed in on
				// all 4 directions (and hence is not a
				// passage.)  Ignore this so that we can cut
				// an entrance and exit out.
			}
		}

		visited := VisitedMap{visitedUnits: map[Point]bool{}, furthestPoints: []Point{}, longestDistance: longestDistance}
		findFurthestOuterCorridorPositionRecursive(unitColumn, unitRow, 0, &visited)

		// Is this the best we've seen so far?
		if visited.longestDistance > longestDistance {
			entranceUnitColumn = unitColumn
			entranceUnitRow = unitRow
			longestDistance = visited.longestDistance
			finalCandidates = visited.furthestPoints
			if m.verbosity > 1 {
				fmt.Printf("[>] Distance from entrance (%v) to exit (%v): %v\n",
					Point{x: entranceUnitColumn, y: entranceUnitRow},
					finalCandidates,
					longestDistance)
			}
		} else if visited.longestDistance == longestDistance {
			// fmt.Printf("[=] Distance from entrance (%v) to exit (%v): %v\n",
			//	Point{x: unitColumn, y: unitRow},
			//	visited.furthestPoints,
			//	visited.longestDistance)
		}
	}

	i := rand.Intn(len(finalCandidates))
	exitUnitColumn, exitUnitRow = finalCandidates[i].x, finalCandidates[i].y

	for i := 0; i < 2; i++ {
		var p Point
		switch i {
		case 0:
			p = Point{x: entranceUnitColumn, y: entranceUnitRow}
		case 1:
			p = Point{x: exitUnitColumn, y: exitUnitRow}
		}

		var horizontal bool

		// The entrance and exit actually need to be on the outer walls, not
		// the outer corridors.
		switch {
		case p.x == 1:              // On left edge of maze.
			p.x--
			horizontal = false
		case p.x == unitWidth - 2:  // On right edge of maze.
			p.x++
			horizontal = false
		case p.y == 1:              // On top edge of maze.
			p.y--
			horizontal = true
		case p.y == unitHeight - 2: // On bottom edge of maze.
			p.y++
			horizontal = true
		}

		// Update the entrance and exit positions post-adjustment.
		switch i {
		case 0:
			entranceUnitColumn, entranceUnitRow = p.x, p.y
		case 1:
			exitUnitColumn, exitUnitRow = p.x, p.y
		}

		// Knock out the entrance/exit itself.
		if m.thickness == 1 {
			m.cells[m.offset(p.x, p.y)] = m.floor
		} else {

			x, y, width, height := m.unitCoordinatesToRect(p.x, p.y)

			if m.thickness != 2 {
				// Only cut the interior of the entrance and
				// exit; leave the borders on the sides.
				if horizontal {
					x, width = x + 1, width - 2
				} else {
					y, height = y + 1, height - 2
				}
			}

			for row := y; row < y + height; row++ {
				for column := x; column < x + width; column++ {
					m.cells[m.offset(column, row)] = m.floor
				}
			}
		}
	}

	return entranceUnitColumn, entranceUnitRow, exitUnitColumn, exitUnitRow, longestDistance + 2
}


// Generates a maze by drawing it on top of any existing runes in m.cells that
// are already present.  Only the blank cells (with a value equal to m.floor)
// will be overwritten.  (This is the key to getting the maze-within-a-maze
// effect to work.)
//
// This is the main generation function.
func (m *Maze) Generate() {

	// If we need to bias the walls of a square maze due to minWallLength
	// being large (so as to make it unicursal), then we choose the
	// direction at random before rendering.
	horizontalDirections := []struct {x, y int}{directions[0], directions[2]}
	verticalDirections := []struct {x, y int}{directions[1], directions[3]}
	var unicursalBiasDirections []struct {x, y int}
	if rand.Intn(2) > 0 {
		unicursalBiasDirections = horizontalDirections
	} else {
		unicursalBiasDirections = verticalDirections
	}


	// The dimensions must be such that they can fit in an odd number of
	// cells of size TxT, where T is the thickness.
	unitWidth, unitHeight := m.width, m.height
	if m.thickness == 2 {
		unitWidth = m.width / 2
		unitHeight = m.height / 2
	} else if m.thickness > 1 {
		unitWidth = (m.width - 1) / (m.thickness - 1)
		unitHeight = (m.height - 1) / (m.thickness - 1)
	}
	if unitWidth % 2 == 0 {
		unitWidth--
	}
	if unitHeight % 2 == 0 {
		unitHeight--
	}

	// Draw a ring of walls around the maze.
	//
	// TODO: We shouldn't overwrite the border where it already exists,
	// and that will necessitate different algorithms for thicknesses of 1
	// and >1.
	if m.thickness == 1 {
		// For a thickness of 1, this looks better than a bunch of
		// intersections.
		m.drawRect(0, 0, unitWidth, unitHeight, m.floor)
	} else {
		// Iterate over the outer perimeter of the maze.
		for _, p := range(rectPerimeter(unitWidth, unitHeight)) {
			x, y, width, height := m.unitCoordinatesToRect(p.x, p.y)
			m.drawRect(x, y, width, height, m.fill)

		}
	}

	// The actual maze algorithm.
	//
	// 1. Consider only the odd coordinates of the maze.  Record which are
	//    completely unoccupied and which are not.
	// 2. Select an odd coordinate at random that is NOT completely
	//    unoccupied.
	// 3. Attempt to draw a line in a randomly-chosen cardinal direction
	//    (recall that a "unit" is a TxT square of cells, where T is
	//    m.thickness.)  The line may have any odd length as long as it is
	//    greater than 3 and it does not touch any other occupied units.
	// 4. If the drawing attempt succeeds, mark the odd coordinates that
	//    the line touched as occupied.
	// 5. If there is still at least one unoccupied coordinate, return to
	//    step 2.
	// 6. Punch holes through the outside walls at two randomly-chosen points
	//    (preferably facing in different cardinal directions) to
	//    represent to entrance and exit.

	// STEP 1
	//
	// Count the number of odd coordinates in the map that are completely
	// unoccupied.  Those are our targets -- where we wish to draw the
	// maze.
	//
	// Suppose our thickness is 1 and we have an existing cells buffer
	// that looks like this (clearly the result of this maze renderer's
	// previous handiwork):
	//
	// +--+  +--+--+--+
	// |..|  |..|..|..|
	// |..|  |..|..|..|
	// +--+  +--+--+--+
	// |..|        |..|
	// |..|        |..|
	// +--+--+--+  +--+
	// |..|..|..|  |..|
	// |..|..|..|  |..|
	// +--+--+--+  +--+
	//
	// From the perspective of the present rendering algorithm, here is
	// the coordinate analysis:
	//
	//    1 2 3 4 5 6 7 8 9 10111213141516
	//  1 * * * * .   * * * * * * * * * *
	//  2 * * * *     * * * * * * * * * *
	//  3 * * * * .   * * * * * * * * * *
	//  4 * * * *     * * * * * * * * * *
	//  5 * * * * .   .   .   .   * * * *
	//  6 * * * *                 * * * *
	//  7 * * * * * * * * * * .   * * * *
	//  8 * * * * * * * * * *     * * * *
	//  9 * * * * * * * * * * .   * * * *
	// 10 * * * * * * * * * *     * * * *
	//
	// Since only odd coordinates are considered, there are 3 + 2 + 3 = 8
	// available spaces.
	//
	// But drawing a line from (say) (7, 11) to (5, 11) would seal off the
	// northwest corridor and instantly ruin the maze.  As it turns out,
	// any unit cell without at least three liberties -- three places to
	// draw walls to -- is unusable by the maze algorithm.  All eight of the
	// available cells in this maze have either two liberties or one, so
	// this maze is actually complete with 0 available spaces.

	numberOfUnoccupiedUnits := 0
	for unitRow := 0; unitRow < unitHeight; unitRow += 2 {
		for unitColumn := 0; unitColumn < unitWidth; unitColumn += 2 {
			x, y, width, height := m.unitCoordinatesToRect(unitColumn, unitRow)
			if m.rectIsPassage(x, y, width, height) {
				numberOfUnoccupiedUnits += 1
			}
		}
	}

	// Another helper routine.
	//
	// Prints a minimal representation of the status of occupied and
	// unoccupied unit rectangles in the maze.  I only use this as a
	// sanity check for myself.
	printOccupiedUnitsDebug := func() {
		if m.verbosity < 1 {
			return
		}
		for unitRow := 0; unitRow < unitHeight; unitRow += 2 {
			for unitColumn := 0; unitColumn < unitWidth; unitColumn += 2 {
				x, y, width, height := m.unitCoordinatesToRect(unitColumn, unitRow)
				unitUnoccupied := m.rectContains(x, y, width, height, m.floor)
				fmt.Printf("%5v ", unitUnoccupied)
			}
			fmt.Printf("\n")
		}
	}
	printOccupiedUnitsDebug()

	misses := 0
	for numberOfUnoccupiedUnits > 0 {

		// STEP 2
		//
		// Select a random position in which to begin.
		//
		// The unit width and height are odd, so round the range up.
		//
		// For instance, if the unit width was 5, unitRow would be one
		// of 0, 2, or 4.
		unitRow := 2 * rand.Intn(unitHeight / 2 + 1)
		unitColumn := 2 * rand.Intn(unitWidth / 2 + 1)

		// Note that the actual wall length will always be between 3
		// and potentialWallLength, regardless of what minWallLength and
		// maxWallLength are set to.  They are _guidelines_, not rigid
		// constraints.
		var (
			minWallLength int = m.minWallLength
			maxWallLength int = m.maxWallLength
		)
		minWallLength += (minWallLength + 1) % 2                        // Round even minimum lengths up to the next highest odd number
		maxWallLength -= (maxWallLength + 1) % 2                        // Round even maximum lengths down to the next lowest odd number
		if maxWallLength < minWallLength {
			maxWallLength = minWallLength // minWallLength <= maxWallLength
		}

		randomDirection := directions[rand.Intn(4)]

		// If longer walls are to be had in one orientation, bias the
		// "random" directions slightly in favor of that orientation.
		switch {
		case minWallLength > unitWidth - 2 && minWallLength > unitHeight - 2 && unitWidth == unitHeight:
			// Bias the maze in a consistent set of directions
			// that were chosen in advance.
			randomDirection = unicursalBiasDirections[rand.Intn(2)]
		case minWallLength > unitHeight - 2 && (minWallLength < unitWidth - 2 || unitWidth > unitHeight):
			// Bias the maze in the horizontal direction.
			if rand.Intn(5) >= 0 {
				randomDirection = horizontalDirections[rand.Intn(2)]
			}
		case minWallLength > unitWidth - 2 && (minWallLength < unitHeight - 2 || unitHeight > unitWidth):
			// Bias the maze in the vertical direction.
			if rand.Intn(5) >= 0 {
				randomDirection = verticalDirections[rand.Intn(2)]
			}
		}

		vx, vy :=  randomDirection.x, randomDirection.y

		x, y, width, height := m.unitCoordinatesToRect(unitColumn, unitRow)
		if m.rectContains(x, y, width, height, m.floor) {
			// The randomly-selected unit was not occupied, so we
			// can't make a wall here.  Try again!
			misses++
			continue
		}

		// STEP 3
		//
		// Will we be able to plop down at least three units (counting
		// the starting position)?
		//
		// Yes, the potential wall length starts at -1.  Remember that
		// it increments in units of two, and that wall lengths must
		// always be odd.
		potentialWallLength := -1
		currentUnitRow, currentUnitColumn := unitRow, unitColumn
		for {
			x, y, width, height := m.unitCoordinatesToRect(currentUnitColumn, currentUnitRow)
			if potentialWallLength == -1 || m.rectContains(x, y, width, height, m.floor) {
				// This spot's clear.  (Or we're in our
				// starting position, which is always valid.)
				//
				// Notice how we advance in groups of two
				// units in order to keep the coordinates
				// odd.
				potentialWallLength += 2
				currentUnitColumn += 2 * vx
				currentUnitRow += 2 * vy

				// If we just moved off-screen, we're done.
				// (This also indicates that we started on a
				// border wall and our random direction vector
				// pointed outward -- our potentialWallLength
				// right now is exactly +1.)
				if currentUnitColumn < 0 || currentUnitRow < 0 || currentUnitColumn >= unitWidth || currentUnitRow >= unitHeight {
					break
				}
			} else {
				// This spot's occupied, so we're done.
				break
			}
		}
		if potentialWallLength < 3 {
			// Can't draw a wall in this direction from this
			// position.
			continue
		}

		// Now that we know how long a wall we can draw, we choose the
		// actual length at random.
		minWallLength = min(potentialWallLength, max(3, minWallLength)) // 3 <= minWallLength <= potentialWallLength
		maxWallLength = min(potentialWallLength, max(3, maxWallLength)) // 3 <= maxWallLength <= potentialWallLength

		if potentialWallLength > maxWallLength {
			// We're drawing a wall that's too short or too long.
			// We'll allow it...some of the time.  (This
			// time-wasting strategy will make the maze provide
			// longer walls more often.)
			if rand.Intn(50) > 0 {
				misses++
				continue
			}
		}

		// This produces a random odd number between minWallLength and maxWallLength.
		wallLength := minWallLength + 2 * rand.Intn((maxWallLength - minWallLength) / 2 + 1)

		if m.verbosity > 1 {
			fmt.Printf("I was able to draw a wall from (%v, %v) to (%v, %v) -- %v units long.  Actually chose %v units (%v <= %v <= %v).\n",
				unitColumn, unitRow,
				unitColumn + vx * (potentialWallLength - 1), unitRow + vy * (potentialWallLength - 1),
				potentialWallLength,
				wallLength,
				minWallLength, wallLength, maxWallLength)
		}

		// Draw the wall.
		currentUnitRow, currentUnitColumn = unitRow, unitColumn
		if m.thickness == 1 {
			x1, y1, _, _ := m.unitCoordinatesToRect(currentUnitColumn, currentUnitRow)
			x2, y2 := x1 + vx * (wallLength - 1), y1 + vy * (wallLength - 1)

			// Force (x1, y1) to be the upper left corner of the rectangle.
			x1, y1, x2, y2 = min(x1, x2), min(y1, y2), max(x1, x2), max(y1, y2)

			m.drawRect(x1, y1, x2 - x1 + 1, y2 - y1 + 1, m.fill)
		} else {
			for i := 0; i < wallLength; i++ {
				x, y, width, height := m.unitCoordinatesToRect(currentUnitColumn, currentUnitRow)
				m.drawRect(x, y, width, height, m.fill)
				currentUnitColumn += vx
				currentUnitRow += vy
			}
		}

		// STEP 4
		//
		// We know how many empty (odd-numbered) cells we just drew over.
		numberOfUnoccupiedUnits -= (wallLength - 1)/2

	} // end (while the maze is not full) [STEP 5]

	entranceUnitColumn, entranceUnitRow, exitUnitColumn, exitUnitRow, solutionDistance := m.findEntranceAndExit(unitWidth, unitHeight)
	m.entrance.x, m.entrance.y, m.entrance.width, m.entrance.height = m.unitCoordinatesToRect(entranceUnitColumn, entranceUnitRow)
	m.exit.x, m.exit.y, m.exit.width, m.exit.height = m.unitCoordinatesToRect(exitUnitColumn, exitUnitRow)
	// m.drawRect(m.entrance.x, m.entrance.y, m.entrance.width, m.entrance.height, '1')
	// m.drawRect(m.exit.x, m.exit.y, m.exit.width, m.exit.height, '2')
	if m.verbosity > 0 {
		fmt.Printf("Maze solution distance: %v.  Number of misses: %v\n", solutionDistance, misses)
	}
}

func  (m *Maze) Set(x, y int, cell rune) {
	if m != nil && m.valid(x, y) {
		m.cells[m.offset(x, y)] = cell
	}
}

func  (m *Maze) Get(x, y int) rune {
	if m != nil && m.valid(x, y) {
		return m.cells[m.offset(x, y)]
	}
	return 0
}

func (m *Maze) Print() {
	for index, y := 0, 0; y < m.height; y++ {
		for x := 0; x < m.width; x, index = x+1, index+1 {
			fmt.Printf("%c", m.cells[index])
		}
		fmt.Printf("\n")
	}
}

func main() {

	m := NewMaze(79, 25)
	parser := argparse.NewParser("maze", "Generates a maze out of Unicode characters with an entrance and an exit.")
	var w *int = parser.Int("W", "width", &argparse.Options{
		Required: false,
		Help: "The width of the maze, in characters",
		Default: m.width,
	})
	var h *int = parser.Int("H", "height", &argparse.Options{
		Required: false,
		Help: "The height of the maze, in characters",
		Default: m.height,
	})
	var t *int = parser.Int("t", "thickness", &argparse.Options{
		Required: false,
		Help: "The thickness of the maze walls (or, equivalently, the size of the maze cells), in characters",
		Default: m.thickness,
	})
	var floor *string = parser.String("f", "floor", &argparse.Options{
		Required: false,
		Help: "The character to use for empty corridor spaces",
		Default: string(m.floor),
	})
	var fill *string = parser.String("F", "fill", &argparse.Options{
		Required: false,
		Help: "The character to use between walls when thickness > 2",
		Default: string(m.fill),
	})
	var intersection *string = parser.String("i", "intersection", &argparse.Options{
		Required: false,
		Help: "The character to use for junctions between maze walls",
		Default: string(m.intersection),
	})
	var horizontal *string = parser.String("x", "horizontal", &argparse.Options{
		Required: false,
		Help: "The character to use for horizontal maze walls",
		Default: string(m.horizontal),
	})
	var vertical *string = parser.String("y", "vertical", &argparse.Options{
		Required: false,
		Help: "The character to use for vertical maze walls",
		Default: string(m.vertical),
	})
	var verbosity *int = parser.FlagCounter("v", "verbose", &argparse.Options{
		Required: false,
		Help: "Verboseness (prints auxiliary information in addition to the maze itself.)  Repeat twice for maximum verboseness.",
	})
	var minWallLength *int = parser.Int("m", "min", &argparse.Options{
		Required: false,
		Help: "The desired minimum wall length, in cells.  This is normally a guideline rather than a constraint, but if this value exceeds the maximum horizontal or vertical wall length, all walls will have the maximum length",
		Default: m.minWallLength,
	})
	var maxWallLength *int = parser.Int("M", "max", &argparse.Options{
		Required: false,
		Help: "The desired maximum wall length, in cells.  This is a guideline, not a constraint, and will be met on a best-effort basis",
		Default: m.maxWallLength,
	})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}
	badCharacterMessage := func(charType, value string) {
		fmt.Fprintf(os.Stderr,
			"The %v argument, \"%v\", has too many characters.  String length must be 1.\n",
			charType,
			value)
		fmt.Print(parser.Usage(err))
	}
	switch {
	case utf8.RuneCountInString(*fill) > 1:
		badCharacterMessage("fill", *fill)
		return
	case utf8.RuneCountInString(*floor) > 1:
		badCharacterMessage("floor", *floor)
		return
	case utf8.RuneCountInString(*intersection) > 1:
		badCharacterMessage("intersection", *intersection)
		return
	case utf8.RuneCountInString(*horizontal) > 1:
		badCharacterMessage("horizontal", *horizontal)
		return
	case utf8.RuneCountInString(*vertical) > 1:
		badCharacterMessage("vertical", *vertical)
		return
	}

	m = NewMaze(*w, *h);
	m.thickness = *t;
	m.fill = ([]rune(*fill))[0]
	m.floor = ([]rune(*floor))[0]
	m.vertical = ([]rune(*vertical)[0])
	m.verbosity = *verbosity;
	m.horizontal = ([]rune(*horizontal))[0]
	m.intersection = ([]rune(*intersection)[0])
	m.minWallLength = *minWallLength
	m.maxWallLength = *maxWallLength
	// m.fill = '█'; m.vertical = '▒'; m.horizontal = '▒'; m.intersection = '▒'; m.floor = '░'

	// rand.Seed(12345678)
	rand.Seed(time.Now().UTC().UnixNano())

	m.Clear()
	m.Generate()

	// x, y, width, height := m.unitCoordinatesToRect(2, 2)
	// m.drawRect(x, y, (width - 1) * 2 + 1, (height - 1) * 2 + 1, m.floor)

	m.Print()
	//n := NewMazeOverExisting(m);
	//n.Generate()
	//n.Print()
}
