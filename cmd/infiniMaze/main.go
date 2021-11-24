package main

import (
	"os"
)

var (
	name        = "infiniMaze"
	version     = "0.0.1"
	description = "InfiniMaze is an infinite, persistent, procedurally generated, explorable maze"
	author      = "Static-Flow"
)

func main() {
	os.Exit(run(os.Args))
}
