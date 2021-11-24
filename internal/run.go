package internal

import "github.com/urfave/cli/v2"

func Run(args []string) int {
	app := newApp()
	if app.Run(args) != nil {
		return 1
	}
	return 0
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.HelpName = name
	app.Usage = description
	app.Version = version
	app.Authors = []*cli.Author{{
		Name: author,
	}}
	app.Flags = flags
	app.HideHelp = true
	app.Action = action
	return app
}
