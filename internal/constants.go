package internal

var (
	name        = "InfiniMaze"
	version     = "0.0.2"
	description = "InfiniMaze is an infinite, persistent, procedurally generated, explorable maze"
	author      = "Static-Flow"
)

// The differences in the x-y coordinate
var (
	dx = map[int]int{Up: 0, Down: 0, Left: -1, Right: 1}
	dy = map[int]int{Up: -1, Down: 1, Left: 0, Right: 0}
)

// Maze cell configurations
// The paths of the Maze is represented in the binary representation.
const (
	Up = 1 << iota
	Down
	Left
	Right
)

// Directions is the set of all the directions
var Directions = []int{Up, Down, Left, Right}

// Opposite directions
var Opposite = map[int]int{Up: Down, Down: Up, Left: Right, Right: Left}
