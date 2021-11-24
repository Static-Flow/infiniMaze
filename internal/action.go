package internal

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/nsf/termbox-go"
	"github.com/urfave/cli"
)

func action(ctx *cli.Context) error {
	config, errors := makeConfig(ctx)
	if errors != nil {
		hasErr := false
		for _, err := range errors {
			if err.Error() != "" {
				_, _ = fmt.Fprintf(os.Stderr, err.Error()+"\n")
				hasErr = true
			}
		}
		if hasErr {
			_, _ = fmt.Fprintf(os.Stderr, "\n")
		}
		_ = cli.ShowAppHelp(ctx)
		return nil
	}
	rand.Seed(config.Seed)

	maze := NewInfiniMaze(config)
	err := termbox.Init()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		return nil
	}
	defer termbox.Close()
	interactive(maze, config.Format)
	return nil
}
