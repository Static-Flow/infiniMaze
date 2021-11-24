package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func StartServer(maze *InfiniMaze) {

	r := gin.Default()
	r.LoadHTMLGlob("./internal/templates/*")
	// Define your handlers
	r.GET("/maze", func(c *gin.Context) {
		if err, mazeData := CompressMazeData(maze.CurrentMaze); err == nil {
			c.JSON(200, mazeData)
		} else {
			c.JSON(500, "")
		}
	})

	r.POST("/mazeImg", func(context *gin.Context) {
		deltaX := context.PostForm("deltaX")
		deltaY := context.PostForm("deltaY")
		if deltaX != "" && deltaY != "" {
			fmt.Println("Hit door, changing mazes")
			maze.ChangeCurrentMaze(deltaX + deltaY)
		} else {
			context.Status(500)
		}
	})

	r.GET("/mazeImg", func(context *gin.Context) {
		maze.CurrentMaze.PrintImage(context.Writer, Ascii, maze.Scale)
	})

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.tmpl", gin.H{"title": "InfiniMaze"})
	})

	// Start service
	r.Run(":3000")
}
