package internal

// Point on a Maze
type Point struct {
	X, Y int
}

// Equal judges the equality of the two points
func (point *Point) Equal(target *Point) bool {
	return point.X == target.X && point.Y == target.Y
}
