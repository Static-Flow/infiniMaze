package main

import "fmt"

/*
This struct holds state for the entire application
*/
type InfiniMaze struct {
	CurrentMaze *Maze            //A Maze Struct containing the state of the currently visible Maze
	mazeHeights int              //Global height for all mazes
	mazeWidths  int              //Global width for all mazes
	globalMazes map[string]*Maze //Map of all generated Maze structs. This allows us to connect adjoining mazes as they are made
}

/*
This handles the users movement first so we can detect whether they entered a "door" to an adjoining maze.
If they do, we swap the current maze, place the user on the right spot of the new maze, and generate any new mazes as needed.
If they aren't entering a "door" we pass along the movement info to the current displayed maze.
*/
func (infiniMaze *InfiniMaze) Move(direction int) {
	//Get the users location
	point := infiniMaze.CurrentMaze.Cursor
	//Get the point they are attempting to move to
	next := point.Advance(direction)
	// If the attempted move places them on a "door" to an adjoining maze
	if infiniMaze.CurrentMaze.PointIsExit(next) {
		//We detect which "door" they entered and switch to that maze while also placing them on the other side of the "door" relative to the new maze
		switch direction {
		case Up:
			infiniMaze.CurrentMaze = infiniMaze.CurrentMaze.NorthMaze
			infiniMaze.CurrentMaze.Cursor = infiniMaze.CurrentMaze.Exits[2].Advance(direction)
		case Left:
			infiniMaze.CurrentMaze = infiniMaze.CurrentMaze.EastMaze
			infiniMaze.CurrentMaze.Cursor = infiniMaze.CurrentMaze.Exits[3].Advance(direction)
		case Down:
			infiniMaze.CurrentMaze = infiniMaze.CurrentMaze.SouthMaze
			infiniMaze.CurrentMaze.Cursor = infiniMaze.CurrentMaze.Exits[1].Advance(direction)
		case Right:
			infiniMaze.CurrentMaze = infiniMaze.CurrentMaze.WestMaze
			infiniMaze.CurrentMaze.Cursor = infiniMaze.CurrentMaze.Exits[0].Advance(direction)
		}
		//Once the new maze is in place, we need to generate that maze's 4 adjoining mazes
		infiniMaze.GenNewMazes(infiniMaze.CurrentMaze)
	} else {
		//The user did not attempt to move through a "door" so we pass the movement to the current maze to validate it
		infiniMaze.CurrentMaze.Move(direction)
	}
}

/**
This function checks if the current Maze has all of its adjoining mazes connected.
If a direction does not have a Maze, first the delta is calculated and checked against the global maze registry.
If a Maze is found, the two mazes are connected.
If one is not found, the new Maze is created, then they are both connected.
*/
func (infiniMaze *InfiniMaze) GenNewMazes(maze *Maze) {
	if maze.NorthMaze == nil {
		//calculate global maze id by performing the movement delta against the current Maze
		stringLocationId := fmt.Sprintf("%d%d", maze.XLocation, maze.YLocation+1)
		//Check if the maze given by the calculated id exists globally
		if possibleMaze := infiniMaze.globalMazes[stringLocationId]; possibleMaze != nil {
			//attach the two mazes to each other
			maze.NorthMaze = possibleMaze
			possibleMaze.SouthMaze = infiniMaze.CurrentMaze
		} else {
			//The maze does not exist so we create a new one and associate them together
			NorthMaze := NewMaze(infiniMaze.mazeHeights, infiniMaze.mazeWidths, maze.XLocation, maze.YLocation+1)
			infiniMaze.globalMazes[stringLocationId] = NorthMaze
			NorthMaze.SouthMaze = maze
			NorthMaze.Generate()
			maze.NorthMaze = NorthMaze
		}
	}
	if maze.SouthMaze == nil {
		stringLocationId := fmt.Sprintf("%d%d", maze.XLocation, maze.YLocation-1)
		if possibleMaze := infiniMaze.globalMazes[stringLocationId]; possibleMaze != nil {
			maze.SouthMaze = possibleMaze
			possibleMaze.NorthMaze = infiniMaze.CurrentMaze

		} else {
			SouthMaze := NewMaze(infiniMaze.mazeHeights, infiniMaze.mazeWidths, maze.XLocation, maze.YLocation-1)
			infiniMaze.globalMazes[stringLocationId] = SouthMaze
			SouthMaze.NorthMaze = maze
			SouthMaze.Generate()
			maze.SouthMaze = SouthMaze
		}
	}
	if maze.WestMaze == nil {
		stringLocationId := fmt.Sprintf("%d%d", maze.XLocation-1, maze.YLocation)
		if possibleMaze := infiniMaze.globalMazes[stringLocationId]; possibleMaze != nil {
			maze.WestMaze = possibleMaze
			possibleMaze.EastMaze = infiniMaze.CurrentMaze

		} else {
			WestMaze := NewMaze(infiniMaze.mazeHeights, infiniMaze.mazeWidths, maze.XLocation-1, maze.YLocation)
			infiniMaze.globalMazes[stringLocationId] = WestMaze
			WestMaze.EastMaze = maze
			WestMaze.Generate()
			maze.WestMaze = WestMaze
		}
	}
	if maze.EastMaze == nil {
		stringLocationId := fmt.Sprintf("%d%d", maze.XLocation+1, maze.YLocation)
		if possibleMaze := infiniMaze.globalMazes[stringLocationId]; possibleMaze != nil {
			maze.EastMaze = possibleMaze
			possibleMaze.WestMaze = infiniMaze.CurrentMaze
		} else {
			EastMaze := NewMaze(infiniMaze.mazeHeights, infiniMaze.mazeWidths, maze.XLocation+1, maze.YLocation)
			infiniMaze.globalMazes[stringLocationId] = EastMaze
			EastMaze.WestMaze = maze
			EastMaze.Generate()
			maze.EastMaze = EastMaze
		}
	}
}

func NewInfiniMaze(config *Config) *InfiniMaze {
	infiniMaze := &InfiniMaze{
		mazeHeights: config.Height,
		mazeWidths:  config.Width,
		globalMazes: map[string]*Maze{},
	}
	//Initial infiniMaze with a starting maze based on user supplied size or default (terminal size), and centered at 0,0
	maze := NewMaze(infiniMaze.mazeHeights, infiniMaze.mazeWidths, 0, 0)
	//Maze keys are their "coordinates" converted to a string: i.e. "00", "-1-1", etc
	infiniMaze.globalMazes["00"] = maze
	infiniMaze.CurrentMaze = maze

	//For the initial maze we generate the 4 mazes to the North, South, East, and West
	infiniMaze.GenNewMazes(maze)

	return infiniMaze
}
