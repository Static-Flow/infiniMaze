package main

import (
	"errors"
	"io"
	"time"

	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

// Config is the command configuration
type Config struct {
	Width       int
	Height      int
	Interactive bool
	Format      *Format
	Seed        int64
	Output      io.Writer
}

func makeConfig(ctx *cli.Context) (*Config, []error) {
	var errs []error

	if ctx.GlobalBool("help") {
		errs = append(errs, errors.New(""))
		return nil, errs
	}

	width := ctx.GlobalInt("width")
	height := ctx.GlobalInt("height")
	if width <= 0 || height <= 0 {
		termWidth, termHeight, err := terminal.GetSize(0)
		if err != nil {
			return nil, []error{err}
		}
		if width <= 0 {
			width = (termWidth - 4) / 4
		}
		if height <= 0 {
			height = (termHeight - 5) / 2
		}
	}

	interactive := ctx.GlobalBool("interactive")

	output := ctx.App.Writer

	format := Color
	if ctx.GlobalString("format") == "ascii" {
		format = Ascii
	}

	seed := int64(ctx.GlobalInt("seed"))
	if !ctx.IsSet("seed") {
		seed = time.Now().UnixNano()
	}

	return &Config{
		Width:       width,
		Height:      height,
		Interactive: interactive,
		Format:      format,
		Seed:        seed,
		Output:      output,
	}, nil
}
