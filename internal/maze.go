package internal

import (
	"bytes"
	"fmt"
	"github.com/nsf/termbox-go"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"strings"
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
	NorthMaze  *Maze     `json:"-"` //Reference to the Maze North of the current maze
	SouthMaze  *Maze     `json:"-"` //Reference to the Maze South of the current maze
	WestMaze   *Maze     `json:"-"` //Reference to the Maze West of the current maze
	EastMaze   *Maze     `json:"-"` //Reference to the Maze East of the current maze
}

// NewMaze creates a new Maze
func NewMaze(height int, width int, xLoc int, yLoc int) *Maze {
	var directions [][]int
	for y := 0; y < height; y++ {
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
		&Point{width / 2, height / 2},
		nil,
		nil,
		nil,
		nil,
	}
	maze.Generate()
	return maze
}

// Neighbors gathers the nearest undecided points
func (maze *Maze) GetUnvisitedDirectionsFromPoint(point *Point) (neighbors []int) {
	//Loop over the 4 Cardinal directions
	for _, direction := range Directions {
		//Move the Point we are collecting neighbors for in the current direction
		next := point.Advance(direction)
		//If the advanced point is in the maze, and hasn't been visited add it to the neighbors list
		if maze.Contains(next) && maze.Directions[next.Y][next.X] == 0 {
			neighbors = append(neighbors, direction)
		}
	}
	return neighbors
}

// Connected judges whether the two points is connected by a path on the maze
func (maze *Maze) Connected(point *Point, target *Point) bool {
	dir := maze.Directions[point.Y][point.X]
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
	unvisitedDirectionsFromPoint := maze.GetUnvisitedDirectionsFromPoint(point)
	//If there are no unvisited neighbors return nil
	if len(unvisitedDirectionsFromPoint) == 0 {
		return nil
	}
	//pick a random neighbor from the list of available neighbors
	randomDirection := unvisitedDirectionsFromPoint[rand.Int()%len(unvisitedDirectionsFromPoint)]
	//Mark the direction you can move from this point to its neighbor by updating the points bitmap and it's connected neighbor's bitmap
	maze.Directions[point.Y][point.X] |= randomDirection
	next := point.Advance(randomDirection)
	maze.Directions[next.Y][next.X] |= Opposite[randomDirection]
	return next
}

// Generate the Maze using this algorithm https://en.wikipedia.org/wiki/Maze_generation_algorithm#Iterative_implementation
func (maze *Maze) Generate() {
	var point *Point
	//Step 1: select a starting Point
	stack := []*Point{maze.Cursor}
	//Step 2: while the stack is not empty
	for len(stack) > 0 {
		//2.1.A: Select a random Point index from the stack
		i := rand.Int() % ((len(stack) + 1) / 2)
		//2.1.B: Retrieve Point at index from stack
		point = stack[i]
		//2.1.C: Pop selected Point from stack
		stack = append(stack[:i], stack[i+1:]...)
		//2.2.A: Choose an unvisited neighbor
		for {
			point = maze.Next(point)
			if point == nil {
				break
			}
			stack = append(stack, point)
		}
	}
	//We ensure that we don't block off the exit "doors" here after Maze generation is done
	//exitSquare := maze.Exits[0].Advance(Right)
	//maze.Next(exitSquare)
	//
	//exitSquare = maze.Exits[1].Advance(Down)
	//maze.Next(exitSquare)
	//
	//exitSquare = maze.Exits[2].Advance(Up)
	//maze.Next(exitSquare)
	//
	//exitSquare = maze.Exits[3].Advance(Left)
	//maze.Next(exitSquare)

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
	return &Point{point.Y + dy[direction], point.X + dx[direction]}
}

// Contains judges whether the argument point is inside Maze or not
func (maze *Maze) Contains(point *Point) bool {
	return 0 <= point.X && point.X < maze.Width && 0 <= point.Y && point.Y < maze.Height
}

// Move the cursor
func (maze *Maze) Move(direction int) {
	point := maze.Cursor
	next := point.Advance(direction)
	// If there's a path on the Maze, we can move the cursor
	if maze.Contains(next) && maze.Directions[point.Y][point.X]&direction == direction {
		//maze.Directions[point.X][point.Y] ^= direction << VisitedOffset
		//maze.Directions[next.X][next.Y] ^= Opposite[direction] << VisitedOffset
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

func plot(img *image.RGBA, x, y, scale int, c color.Color) {
	for dy := 0; dy < scale; dy++ {
		for dx := 0; dx < scale; dx++ {
			img.Set(x*scale+dx, y*scale+dy, c)
		}
	}
}

// PrintImage outputs the maze to the IO writer as PNG image
func (maze *Maze) PrintImage(writer io.Writer, format *Format, scale int) {
	var buf bytes.Buffer
	maze.Print(&buf, format)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	width := len(lines[0]) / 2
	height := len(lines)
	img := image.NewRGBA(image.Rect(0, 0, width*scale, height*scale))
	green := color.RGBA{0, 255, 0, 255}
	floorColor := color.RGBA{119, 136, 153, 255}
	for y := 0; y < height; y++ {
		if y >= len(lines) {
			continue
		}
		for x := 0; x < width; x++ {
			if x*2 >= len(lines[y]) {
				continue
			}
			switch lines[y][x*2 : x*2+2] {
			case "##":
				plot(img, x, y, scale, color.Black)
			case "<<", "VV", "^^", ">>":
				plot(img, x, y, scale, green)
			default:
				plot(img, x, y, scale, floorColor)
			}
		}
	}
	_ = png.Encode(writer, img)
}

// Write out the Maze to the writer channel
//It walks the 2D array of cells and follows this algorithm:
/*
1. For each column, the first cell is the top row. If the column is the middle column and the "Up" direction check this cell is the North exit
2. For the rest of the column:
	2.a
*/
func (maze *Maze) Write(writer chan string, format *Format) {
	//Print global maze location
	maze.printLocation()
	writer <- "\n"
	for y, row := range maze.Directions {
		// There are two lines printed for each Maze line
		for _, direction := range []int{Up, Right} {
			// The left wall
			if y == maze.Height/2 && direction == Right {
				writer <- format.ExitLeft
			} else {
				writer <- format.Wall
			}
			for x, directions := range row {
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
				} else if direction == Up && x == maze.Width/2 && y == 0 {
					writer <- format.ExitUp
				} else if direction == Right && x == maze.Width-1 && y == maze.Height/2 {
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
	for x := 0; x < maze.Width; x++ {
		if x == maze.Width/2 {
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
