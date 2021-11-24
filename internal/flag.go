package internal

import "github.com/urfave/cli/v2"

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "width",
		Usage: "The width of the infiniMaze",
	},
	&cli.StringFlag{
		Name:  "height",
		Usage: "The height of the infiniMaze",
	},
	&cli.StringFlag{
		Name:  "format",
		Usage: "Output format, `default` or `ascii`",
	},
	&cli.StringFlag{
		Name:  "seed",
		Usage: "The random seed",
	},
	&cli.BoolFlag{
		Name:  "help, h",
		Usage: "Shows the help of the command",
	},
}
