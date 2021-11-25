package internal

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

func StartServer(infiniMaze *InfiniMaze) error {
	r := gin.Default()
	gob.Register(Maze{})
	store := cookie.NewStore([]byte("Thisissupersecret456!"))
	r.Use(sessions.Sessions("mysession", store))
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)
	// Define your handlers
	r.GET("/maze", func(c *gin.Context) {
		if err, mazeData := CompressMazeData(infiniMaze.CurrentMaze); err == nil {
			c.JSON(200, mazeData)
		} else {
			c.JSON(500, "")
		}
	})

	//Method for updating position of user within maze
	r.POST("/move", func(context *gin.Context) {
		session := sessions.Default(context)
		v := session.Get("id")
		if v == nil {
			context.Redirect(302, "/")
		} else {
			currentPosition := session.Get("position").([]int)
			deltaForm := &MazeMoveDeltaChange{}
			err := context.ShouldBind(deltaForm)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "bad delta"})
			} else if deltaForm.MazeDeltaChangeValidation(currentPosition) {
				session.Set("position", []int{*deltaForm.DeltaX, *deltaForm.DeltaY})
				_ = session.Save()
				fmt.Printf("Updated user position in maze to %d,%d\n", *deltaForm.DeltaX, *deltaForm.DeltaY)
				context.Status(204)
			} else {
				context.JSON(http.StatusBadRequest, gin.H{"error": "bad delta"})
			}
		}
	})

	//Method for retrieving new maze in the direction the user exited their current maze
	r.POST("/mazeImg", func(context *gin.Context) {
		session := sessions.Default(context)
		v := session.Get("id")
		if v == nil {
			context.Redirect(302, "/")
		} else {
			deltaForm := &MazeGlobalDeltaChange{}
			err := context.ShouldBind(deltaForm)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "bad delta"})
			}
			/*
				Here we do two checks:
				1. Was the global delta the user sent us valid, i.e. are both values present and the delta is only 1 move in 1 Up,Down,Left,Right direction
				2. Are the currently near an exit. This is to validate that they didn't "warp" themselves via client side edits to be near one
			*/
			if deltaForm.MazeDeltaChangeValidation(infiniMaze.mazeSessions[v.(string)]) &&
				infiniMaze.ValidateUserIsNearMapExit(session) {
				mazeId := fmt.Sprintf("%d,%d", *deltaForm.DeltaX, *deltaForm.DeltaY)
				fmt.Printf("Hit door, changing mazes to %s\n", mazeId)
				session.Set("globalIndex", mazeId)
				_ = session.Save()
				infiniMaze.ChangeCurrentMazeForSession(session)
				context.Status(204)
			} else {
				context.JSON(http.StatusBadRequest, gin.H{"error": "bad delta"})
			}
		}
	})

	r.GET("/mazeImg", func(context *gin.Context) {
		session := sessions.Default(context)
		v := session.Get("id")
		if v == nil {
			context.Redirect(302, "/")
		} else {
			infiniMaze.mazeSessions[v.(string)].PrintImage(context.Writer, Ascii, infiniMaze.scale)
		}
	})

	r.GET("/", func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		id := session.Get("id")
		globalLocation := session.Get("globalIndex")
		position := session.Get("position")
		//If the user doesn't have a session make one
		if id == nil {
			newId := uuid.New()
			session.Set("id", newId.String())
			_ = session.Save()
		}
		if globalLocation == nil {
			fmt.Println("adding global index to session")
			globalLocation = "0,0"
			session.Set("globalIndex", globalLocation)
			_ = session.Save()
		}
		if position == nil {
			fmt.Println("adding position to session")
			position = []int{infiniMaze.mazeWebWidths / 2, infiniMaze.mazeWebHeights / 2}
			session.Set("position", position)
			_ = session.Save()
		}
		infiniMaze.ChangeCurrentMazeForSession(session)
		globalPosition := strings.Split(globalLocation.(string), ",")
		data := gin.H{
			"title":           "InfiniMaze",
			"globalLocationX": globalPosition[0],
			"globalLocationY": globalPosition[1],
			"position":        position,
		}
		fmt.Printf("%+v\n", data)
		ctx.HTML(
			200,
			"/internal/templates/index.tmpl",
			data)
	})

	// Start service
	return r.Run(":" + infiniMaze.webPort)
}
