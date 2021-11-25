package internal

import (
	"fmt"
	"math"
)

/*
This type is for mapping POST request data from the /mazeImg endpoint
which contains the global delta movement between mazes
*/
type MazeGlobalDeltaChange struct {
	DeltaX *int `form:"deltaX"`
	DeltaY *int `form:"deltaY"`
}

//Validates that the user supplied delta contains both values and is only 1 distance away
func (mgdc *MazeGlobalDeltaChange) MazeDeltaChangeValidation(maze *Maze) bool {
	if mgdc.DeltaX == nil || mgdc.DeltaY == nil {
		fmt.Println("Missing Delta field")
		return false
	}
	if math.Abs(float64(maze.XLocation-*mgdc.DeltaX)) == 0 && math.Abs(float64(maze.YLocation-*mgdc.DeltaY)) == 0 {
		//If we haven't moved in the X plane or the Y plane its a useless move

		return false
	} else if math.Abs(float64(maze.XLocation-*mgdc.DeltaX)) == 0 && math.Abs(float64(maze.YLocation-*mgdc.DeltaY)) == 1 {
		//If we haven't moved in the X plane and only moved 1 direction in the Y plane it is valid
		return true
	} else if math.Abs(float64(maze.XLocation-*mgdc.DeltaX)) == 1 && math.Abs(float64(maze.YLocation-*mgdc.DeltaY)) == 0 {
		//If we haven't moved in the Y plane and only moved 1 direction in the X plane it is valid
		return true
	} else {
		//Catch all for any other malformed inputs
		return false
	}
}

/*
This type is for mapping POST request data from the /move endpoint
which contains the global delta movement within the users maze
*/
type MazeMoveDeltaChange struct {
	DeltaX *int `form:"deltaX"`
	DeltaY *int `form:"deltaY"`
}

//Validates that the user supplied delta contains both values and is only moving 1 distance away
func (mmdc *MazeMoveDeltaChange) MazeDeltaChangeValidation(currentPosition []int) bool {
	fmt.Printf("Validating user movement.They are at: %d,%d. Moving to %d,%d\n", currentPosition[0], currentPosition[1], *mmdc.DeltaX, *mmdc.DeltaY)
	if mmdc.DeltaX == nil || mmdc.DeltaY == nil {
		fmt.Println("Missing Delta field")
		return false
	}
	if math.Abs(float64(currentPosition[0]-*mmdc.DeltaX)) == 0 && math.Abs(float64(currentPosition[1]-*mmdc.DeltaY)) == 0 {
		//If we haven't moved in the X plane or the Y plane its a useless move

		return false
	} else if math.Abs(float64(currentPosition[0]-*mmdc.DeltaX)) == 0 && math.Abs(float64(currentPosition[1]-*mmdc.DeltaY)) == 20 {
		//If we haven't moved in the X plane and only moved 1 direction in the Y plane it is valid
		return true
	} else if math.Abs(float64(currentPosition[0]-*mmdc.DeltaX)) == 20 && math.Abs(float64(currentPosition[1]-*mmdc.DeltaY)) == 0 {
		//If we haven't moved in the Y plane and only moved 1 direction in the X plane it is valid
		return true
	} else {
		//Catch all for any other malformed inputs
		return false
	}
}
