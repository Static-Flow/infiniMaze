package main

import (
	"github.com/itchyny/maze/internal"
	"os"
)

func main() {
	os.Exit(internal.Run(os.Args))
}
