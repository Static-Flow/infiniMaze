package internal

import (
	"fmt"
	"github.com/gin-contrib/sessions"
)

/*
This struct holds state for the entire application
*/
type InfiniMaze struct {
	CurrentMaze    *Maze            //A Maze Struct containing the state of the currently visible Maze
	mazeHeights    int              //Global height for all mazes
	mazeWidths     int              //Global width for all mazes
	mazeWebHeights int              //Global height for all web mazes
	mazeWebWidths  int              //Global width for all web mazes
	scale          int              //How much to scale up the PNG image of the maze
	globalMazes    map[string]*Maze //Map of all generated Maze structs. This allows us to connect adjoining mazes as they are made
	mazeSessions   map[string]*Maze //Map of current users to the map they are on
	webPort        string           //Port for web server to listen on
}

func (infiniMaze *InfiniMaze) ChangeCurrentMaze(globalIndexId string) {
	if newMaze, exists := infiniMaze.globalMazes[globalIndexId]; exists {
		fmt.Printf("Changing Maze to: %s\n", globalIndexId)
		infiniMaze.CurrentMaze = newMaze
		infiniMaze.GenNewMazes(infiniMaze.CurrentMaze)
	} else {
		fmt.Println("Maze " + globalIndexId + " does not exist")
	}
}

/*
Checks the current users position against where the exits are on the current maze.
If the move is valid we update the users position in the session
*/
func (infiniMaze *InfiniMaze) ValidateUserIsNearMapExit(session sessions.Session) bool {
	currentPosition := session.Get("position").([]int)
	if currentPosition[0] == infiniMaze.mazeWebWidths/2 && currentPosition[1] == 20 {
		//Check North exit which should be X=WebWidth/2,Y=20
		//Move current position to above the bottom door of the North maze the user moved to
		currentPosition[1] = infiniMaze.mazeWebHeights - 20
		session.Set("position", currentPosition)
		_ = session.Save()
		return true
	} else if currentPosition[0] == infiniMaze.mazeWebWidths/2 && currentPosition[1] == infiniMaze.mazeWebHeights-20 {
		//Check South exit which should be X=WebWidth/2,Y=WebHeight-20
		//Move current position to below the top door of the South maze the user moved to
		currentPosition[1] = 20
		session.Set("position", currentPosition)
		_ = session.Save()
		return true
	} else if currentPosition[0] == infiniMaze.mazeWebWidths-20 && currentPosition[1] == infiniMaze.mazeWebHeights/2 {
		//Check East exit which should be X=WebWidth-20,Y=WebHeight/2
		//Move current position to the right of the left door of the East maze the user moved to
		currentPosition[0] = 20
		session.Set("position", currentPosition)
		_ = session.Save()
		return true
	} else if currentPosition[0] == 20 && currentPosition[1] == infiniMaze.mazeWebHeights/2 {
		//Check West exit which should be X=20,Y=WebHeight/2
		//Move current position to the left of the right door of the West maze the user moved to
		currentPosition[0] = infiniMaze.mazeWebWidths - 20
		session.Set("position", currentPosition)
		_ = session.Save()
		return true
	} else {
		fmt.Printf("User is not near the exit they say they are. They are at: %d,%d\n", currentPosition[0], currentPosition[1])
		return false
	}
}

func (infiniMaze *InfiniMaze) ChangeCurrentMazeForSession(session sessions.Session) {
	id := session.Get("id").(string)
	globalLocation := session.Get("globalIndex").(string)
	if newMaze, exists := infiniMaze.globalMazes[globalLocation]; exists {
		fmt.Printf("Changing Maze for session: %s to: %s\n", id, globalLocation)
		infiniMaze.mazeSessions[id] = newMaze
		infiniMaze.GenNewMazes(newMaze)
	} else {
		fmt.Println("Maze " + globalLocation + " does not exist")
	}
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
		stringLocationId := fmt.Sprintf("%d,%d", maze.XLocation, maze.YLocation+1)
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
		stringLocationId := fmt.Sprintf("%d,%d", maze.XLocation, maze.YLocation-1)
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
		stringLocationId := fmt.Sprintf("%d,%d", maze.XLocation-1, maze.YLocation)
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
		stringLocationId := fmt.Sprintf("%d,%d", maze.XLocation+1, maze.YLocation)
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
		mazeHeights:    config.Height,
		mazeWebHeights: config.Height * config.Scale * 2,
		mazeWebWidths:  config.Width * config.Scale * 2,
		mazeWidths:     config.Width,
		scale:          config.Scale,
		globalMazes:    map[string]*Maze{},
		webPort:        config.WebPort,
		mazeSessions:   map[string]*Maze{},
	}
	//Initial infiniMaze with a starting maze based on user supplied size or default (terminal size), and centered at 0,0
	maze := NewMaze(infiniMaze.mazeHeights, infiniMaze.mazeWidths, 0, 0)
	//Maze keys are their "coordinates" converted to a string: i.e. "00", "-1-1", etc
	infiniMaze.globalMazes["0,0"] = maze
	infiniMaze.CurrentMaze = maze

	//For the initial maze we generate the 4 mazes to the North, South, East, and West
	infiniMaze.GenNewMazes(maze)

	return infiniMaze
}
