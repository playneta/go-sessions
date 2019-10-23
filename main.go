package main

import (
	"os"

	"github.com/playneta/go-sessions/src"
	"github.com/urfave/cli"
)

var (
	Version = "development"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-lessons"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name: "serve",
			Action: func(ctx *cli.Context) {
				src.Run()
			},
		},
		{
			Name: "migrate:up",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dir",
					Value: "./migrations",
				},
			},
			Action: func(ctx *cli.Context) {
				dir := ctx.String("dir")
				src.Migrate(dir)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
