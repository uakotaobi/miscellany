package main
import (
	"fmt"
	"math/rand"
)

type Maze struct {
	width int
	height int
	thickness int
	cells []rune
	intersection rune
	vertical rune
	horizontal rune
	floor rune
	fill rune
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

func NewMaze(width, height int) Maze {
	m := Maze{thickness: 1, intersection: '+', horizontal: '-', vertical: '|', floor: ' ', fill: '.' }
	m.setSize(width, height)

	// Clear the buffer.
	m.cells = make([]rune, m.width * m.height)
	for index := range(m.cells) {
		m.cells[index] = m.floor
	}

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

// Alters the maze's dimensions.  The maze will need re-rendering after the
// call.
func (m *Maze) setSize(newWidth, newHeight int) {
	// Moved these checks into generateMaze() itself.
	// if newWidth < 3 {
	//	newWidth = 3
	// }
	// if newWidth % 2 == 0 {
	//	// Even width, round down.
	//	newWidth -= 1
	// }
	// if newHeight < 3 {
	//	newHeight = 3
	// }
	// if newHeight % 2 == 0 {
	//	newHeight -= 1
	// }
	m.width = newWidth;
	m.height = newHeight;
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

// This utility function determines if the given coordinate has orthogonally
// adjacent neighbors whose cells all contain the given value.
//
// The neighbor positions are specified by passing the appropriate bitwise
// combination of the left, right, up and down constants into the mask
// argument.
//
// Note that a mask of 0 will always return true.
func (m *Maze) hasNeighbor(x, y int, mask mask, neighborValue rune) bool {
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

// Generates a maze by drawing it on top of the given set of cells.  Only the
// blank cells (with a value equal to m.floor) will be overwritten.
//
// This is the main generation function.
func (m *Maze) generateMaze(existingCells []rune) {

	// Start by duplicating the existing buffer.
	m.cells = make([]rune, m.width * m.height)
	copy(m.cells, existingCells)

	// Ensure that the borders of the buffer are lined with walls.

	// The dimensions must be such that they can fit in an odd number of
	// cells of size TxT, where T is the thickness.
	unitWidth, unitHeight := m.width, m.height
	if m.thickness > 1 {
		unitWidth = (m.width - 1) / (m.thickness - 1)
		unitHeight = (m.height - 1) / (m.thickness - 1)
	}
	if unitWidth % 2 == 0 {
		unitWidth--
	}
	if unitHeight % 2 == 0 {
		unitHeight--
	}

	// Helpful conversion routines.
	//
	// Converts a "unit coordinate" (a subdivision of the maze grid into
	// slightly-overlapping boxes) into a rectangle.
	unitCoordinatesToRect := func(unitColumn, unitRow int) (x, y, width, height int) {
		if m.thickness == 1 {
			return unitColumn, unitRow, 1, 1
		}
		x = unitColumn * (m.thickness - 1)
		y = unitRow * (m.thickness - 1)
		return x, y, m.thickness, m.thickness
	}

	// Draw a ring of walls around the maze.
	if m.thickness == 1 {
		// For a thickness of 1, this looks better than a bunch of
		// intersections.
		m.drawRect(0, 0, unitWidth, unitHeight, m.floor)
	} else {
		// Iterate over the outer perimeter of the maze.
		for i := 0; i < 2 * (unitWidth + unitHeight); i++ {
			var x, y int
			switch {
			case i < unitWidth:
				// Top edge.
				x, y, _, _ = unitCoordinatesToRect(i, 0)
			case i >= unitWidth && i < unitWidth + unitHeight:
				// Right edge.
				x, y, _, _ = unitCoordinatesToRect(unitWidth - 1, i - unitWidth)
			case i >= unitWidth + unitHeight && i < (2 * unitWidth + unitHeight):
				// Bottom edge (moving left.)
				x, y, _, _ = unitCoordinatesToRect(unitWidth - (i - unitWidth - unitHeight) - 1, unitHeight - 1)
			default:
				// Left edge (moving up.)
				x, y, _, _ = unitCoordinatesToRect(0, unitHeight - (i - (2 * unitWidth) - unitHeight) - 1)
			}

			m.drawRect(x, y, m.thickness, m.thickness, m.fill)
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
	// This 2D data structure -- a slice of slices -- contains a boolean
	// for each of the odd unit of the maze (where a "unit" can take up
	// more than one character cell if m.thickness > 1.)
	occupiedUnits := [][]bool{}
	numberOfUnoccupiedUnits := 0
	for unitRow := 0; unitRow < unitHeight; unitRow += 2 {
		occupiedUnitsCurrentRow := []bool{}
		for unitColumn := 0; unitColumn < unitWidth; unitColumn += 2 {
			x, y, _, _ := unitCoordinatesToRect(unitColumn, unitRow)
			unitUnoccupied := m.rectContains(x, y, m.thickness, m.thickness, m.floor)
			occupiedUnitsCurrentRow = append(occupiedUnitsCurrentRow, unitUnoccupied)
			if unitUnoccupied {
				numberOfUnoccupiedUnits += 1
			}
		}
		occupiedUnits = append(occupiedUnits, occupiedUnitsCurrentRow)
	}

	// Another helper routine.
	//
	// Prints a minimal representation of the status of occupied and
	// unoccupied unit rectangles in the maze.  I only use this as a
	// sanity check for myself.
	printOccupiedUnitsDebug := func() {
		for unitRow := 0; unitRow < unitHeight; unitRow += 2 {
			for unitColumn := 0; unitColumn < unitWidth; unitColumn += 2 {
				x, y, _, _ := unitCoordinatesToRect(unitColumn, unitRow)
				unitUnoccupied := m.rectContains(x, y, m.thickness, m.thickness, m.floor)
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

		const minLengthInUnits = 3
		var vx, vy, maxLengthInUnits int

		switch randomDirection := mask(1 << uint(rand.Intn(4))); randomDirection {
		case left:
			vx, vy, maxLengthInUnits = -1, 0, unitColumn
		case up:
			vx, vy, maxLengthInUnits = 0, -1, unitRow
		case right:
			vx, vy, maxLengthInUnits = 1, 0, unitWidth - unitColumn
		case down:
			vx, vy, maxLengthInUnits = 0, 1, unitHeight - unitRow
		}

		x, y, width, height := unitCoordinatesToRect(unitColumn, unitRow)
		if m.rectContains(x, y, width, height, m.floor) {
			// The randomly-selected unit was not occupied, so we
			// can't make a wall here.  Try again!
			misses++
			continue
		}

		// Will we be able to plop down at least three units (counting
		// the starting position)?
		//
		// Yes, the potential wall length starts at -1.  Remember that
		// it increments in units of two, and that wall lengths must
		// always be odd.
		potentialWallLength := -1
		currentUnitRow, currentUnitColumn := unitRow, unitColumn
		for {
			x, y, width, height := unitCoordinatesToRect(currentUnitColumn, currentUnitRow)
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
		if potentialWallLength < minLengthInUnits {
			// Can't draw a wall in this direction from this
			// position.
			continue
		}

		// Now that we know how long a wall we can draw, we choose the
		// length at random.
		// (min and max wall lengths should be integers, not
		// floating-point scaling factors.)
		fmt.Printf("I was able to draw a wall from (%v, %v) to (%v, %v) -- %v units long\n",
			unitColumn, unitRow,
			unitColumn + vx * (potentialWallLength - 1), unitRow + vy * (potentialWallLength - 1),
			potentialWallLength)

		_, _, _ = vx, vy, maxLengthInUnits
		break

	} // end (while the maze is not full)
	fmt.Printf("Number of misses: %v\n", misses)
}

// Draws a new maze.
func (m *Maze) Generate() {

	// Note that the existing maze may or may not be empty.
	m.generateMaze(m.cells)
}

// Renders a maze in the empty areas of an existing maze.
//
// The idea here is that we can generate a large maze using thickness>1, then
// generate a small maze on top of that in order to create a sort of
// "super-maze."
//
//  To generate a traditional maze, use Generate().
func (m *Maze) GenerateOverExisting(n *Maze) {
	// Adopt the existing maze's parameters.
	//
	// TODO: We could use our own parameters instead, but then we'd need
	// to know which cells in the existing maze to treat like floors or
	// walls.
	m = &Maze{
		width: n.width,
		height: n.height,
		thickness: n.thickness,
		intersection: n.intersection,
		horizontal: n.horizontal,
		vertical: n.vertical,
		floor: n.floor,
		fill: n.fill,
	}
	m.generateMaze(n.cells)
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

	m := NewMaze(60, 20)
	m.drawRect(0, 0, m.width, m.height, m.floor)
	m.thickness = 4

	rand.Seed(1234567)
	m.Generate()
	// foo := [][]int { {59, 10}, {5, 19} }
	// for _, value := range foo {
	//	x, y := value[0], value[1]
	//	m.Set(x, y, 'X')
	//	l := m.hasNeighbor(x, y, left, m.floor)
	//	r := m.hasNeighbor(x, y, right, m.floor)
	//	u := m.hasNeighbor(x, y, up, m.floor)
	//	d := m.hasNeighbor(x, y, down, m.floor)
	//	fmt.Printf("has blank neighbor: (%v, %v): left=%v up=%v right=%v down=%v\n", x, y, l, u, r, d)
	// }
	m.Print()
}