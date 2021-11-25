package internal

import (
	"errors"
	"io"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

// Config is the command configuration
type Config struct {
	Width   int
	Height  int
	Format  *Format
	Seed    int64
	Output  io.Writer
	IsWeb   bool
	Scale   int
	WebPort string
}

func makeConfig(ctx *cli.Context) (*Config, []error) {
	var errs []error

	if ctx.Bool("help") {
		errs = append(errs, errors.New(""))
		return nil, errs
	}

	width := ctx.Int("width")
	height := ctx.Int("height")
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

	output := ctx.App.Writer

	format := Color
	if ctx.String("format") == "ascii" {
		format = Ascii
	}

	scale := ctx.Int("scale")

	seed := int64(ctx.Int("seed"))
	if !ctx.IsSet("seed") {
		seed = time.Now().UnixNano()
	}
	web := ctx.Bool("web")
	if web {
		format = Ascii
		//I don't really understand how to generalize this yet so I hardcode it for the web version
		//TODO: understand how to get image sizes based on these 3 values
		scale = 20
		width = 25
		height = 25
	}

	return &Config{
		Width:   width,
		Height:  height,
		Format:  format,
		Seed:    seed,
		Output:  output,
		Scale:   scale,
		IsWeb:   web,
		WebPort: ctx.String("webPort"),
	}, nil
}
