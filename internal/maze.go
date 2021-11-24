package internal

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
	"math/rand"
)

/*
Maze represents the configuration of a specific Maze within InfiniMaze
*/
type Maze struct {
	Directions [][]int   //Each Point on the map is represented as an integer that defines the directions that can be traveled from that Point
	Height     int       //Height of Maze
	Width      int       //Width of Maze
	XLocation  int       //Global X position within InfiniMaze
	YLocation  int       //Global Y position within InfiniMaze
	Exits      [4]*Point //Point array for the four "doors" leading the adjoining Mazes
	Cursor     *Point    //Users location within the Maze
	NorthMaze  *Maze     //Reference to the Maze North of the current maze
	SouthMaze  *Maze     //Reference to the Maze South of the current maze
	WestMaze   *Maze     //Reference to the Maze West of the current maze
	EastMaze   *Maze     //Reference to the Maze East of the current maze
}

// NewMaze creates a new Maze
func NewMaze(height int, width int, xLoc int, yLoc int) *Maze {
	var directions [][]int
	for x := 0; x < height; x++ {
		directions = append(directions, make([]int, width))
	}
	maze := &Maze{
		directions,
		height,
		width,
		xLoc,
		yLoc,
		[4]*Point{
			{height / 2, -1},    //left wall exit
			{-1, width / 2},     //top wall exit
			{height, width / 2}, //bottom wall exit
			{height / 2, width}, //right wall exit
		},
		&Point{height / 2, width / 2},
		nil,
		nil,
		nil,
		nil,
	}
	maze.Generate()
	return maze
}

// Neighbors gathers the nearest undecided points
func (maze *Maze) Neighbors(point *Point) (neighbors []int) {
	for _, direction := range Directions {
		next := point.Advance(direction)
		if maze.Contains(next) && maze.Directions[next.X][next.Y] == 0 {
			neighbors = append(neighbors, direction)
		}
	}
	return neighbors
}

// Connected judges whether the two points is connected by a path on the maze
func (maze *Maze) Connected(point *Point, target *Point) bool {
	dir := maze.Directions[point.X][point.Y]
	for _, direction := range Directions {
		if dir&direction != 0 {
			next := point.Advance(direction)
			if next.X == target.X && next.Y == target.Y {
				return true
			}
		}
	}
	return false
}

// Next advances the Maze path randomly and returns the new point
func (maze *Maze) Next(point *Point) *Point {
	neighbors := maze.Neighbors(point)
	if len(neighbors) == 0 {
		return nil
	}
	direction := neighbors[rand.Int()%len(neighbors)]
	maze.Directions[point.X][point.Y] |= direction
	next := point.Advance(direction)
	maze.Directions[next.X][next.Y] |= Opposite[direction]
	return next
}

// Generate the Maze
func (maze *Maze) Generate() {
	point := maze.Cursor
	stack := []*Point{point}
	for len(stack) > 0 {
		for {
			point = maze.Next(point)
			if point == nil {
				break
			}
			stack = append(stack, point)
		}
		i := rand.Int() % ((len(stack) + 1) / 2)
		point = stack[i]
		stack = append(stack[:i], stack[i+1:]...)
	}
	//We ensure that we don't block off the exit "doors" here after Maze generation is done
	exitSquare := maze.Exits[0].Advance(Right)
	maze.Next(exitSquare)

	exitSquare = maze.Exits[1].Advance(Down)
	maze.Next(exitSquare)

	exitSquare = maze.Exits[2].Advance(Up)
	maze.Next(exitSquare)

	exitSquare = maze.Exits[3].Advance(Left)
	maze.Next(exitSquare)

}

//Prints the user global position in the top right
func (maze *Maze) printLocation() {
	str := fmt.Sprintf("%d : %d      ", maze.XLocation, maze.YLocation)
	fg, bg := termbox.ColorDefault, termbox.ColorDefault
	for i, c := range str {
		termbox.SetCell(4*maze.Width+i-8, 1, c, fg, bg)
	}

}

//Check for whether point on Maze is an exit "door" to another Maze
func (maze *Maze) PointIsExit(point *Point) bool {
	for _, exitPoint := range maze.Exits {
		if exitPoint.X == point.X && exitPoint.Y == point.Y {
			return true
		}
	}
	return false
}

// Advance the point forward by the argument direction
func (point *Point) Advance(direction int) *Point {
	return &Point{point.X + dx[direction], point.Y + dy[direction]}
}

// Contains judges whether the argument point is inside Maze or not
func (maze *Maze) Contains(point *Point) bool {
	return 0 <= point.X && point.X < maze.Height && 0 <= point.Y && point.Y < maze.Width
}

// Move the cursor
func (maze *Maze) Move(direction int) {
	point := maze.Cursor
	next := point.Advance(direction)
	// If there's a path on the Maze, we can move the cursor
	if maze.Contains(next) && maze.Directions[point.X][point.Y]&direction == direction {
		maze.Directions[point.X][point.Y] ^= direction << VisitedOffset
		maze.Directions[next.X][next.Y] ^= Opposite[direction] << VisitedOffset
		maze.Cursor = next
	}
}

// Print out the Maze to the IO writer
func (maze *Maze) Print(writer io.Writer, format *Format) {
	strwriter := make(chan string)
	go maze.Write(strwriter, format)
	for {
		str := <-strwriter
		switch str {
		case "\u0000":
			return
		default:
			_, _ = fmt.Fprint(writer, str)
		}
	}
}

// Write out the Maze to the writer channel
func (maze *Maze) Write(writer chan string, format *Format) {
	//Print global maze location
	maze.printLocation()
	writer <- "\n"
	for x, row := range maze.Directions {
		// There are two lines printed for each Maze line
		for _, direction := range []int{Up, Right} {
			// The left wall
			if x == maze.Height/2 && direction == Right {
				writer <- format.ExitLeft
			} else {
				writer <- format.Wall
			}
			for y, directions := range row {
				// In the `direction == Right` line, we print the path cell
				if direction == Right {
					if maze.Cursor.X == x && maze.Cursor.Y == y {
						writer <- format.Cursor
					} else {
						writer <- format.Path
					}
				}
				if directions&direction != 0 { // If there is a path in the direction (Up or Right) on the Maze
					writer <- format.Path
				} else if direction == Up && y == maze.Width/2 && x == 0 {
					writer <- format.ExitUp
				} else if direction == Right && x == maze.Height/2 && y == maze.Width-1 {
					writer <- format.ExitRight
				} else {
					writer <- format.Wall
				}
				if direction == Up {
					writer <- format.Wall
				}
			}
			writer <- "\n"
		}
	}
	// Print the bottom wall of the Maze
	writer <- format.Wall
	for y := 0; y < maze.Width; y++ {
		if y == maze.Width/2 {
			writer <- format.ExitDown
		} else {
			writer <- format.Wall
		}
		writer <- format.Wall
	}
	writer <- "\n\n"
	// Inform that we finished printing the Maze
	writer <- "\u0000"
}
