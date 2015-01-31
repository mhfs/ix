package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[options] {{end}}

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
}

func main() {
	app := cli.NewApp()
	app.Name = "thisweek"
	app.Usage = "create a report of your team's activity this week (or for whenever you'd like)"
	app.Version = "0.0.1"
	app.Author = "Marcelo Silveira"
	app.Email = "marcelo@mhfs.com.br"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "repo, r",
			Value: "",
			Usage: "GitHub repository to analyze e.g. mhfs/thisweek",
		},
		cli.StringFlag{
			Name: "from, f",
			// Value: "2015-01-25", // FIXME make default beginning of current week
			Usage: "from date, inclusive",
		},
		cli.StringFlag{
			Name: "to, t",
			// Value: "2015-01-31", // FIXME make default end of current week
			Usage: "to date, inclusive",
		},
		cli.StringSliceFlag{
			Name:  "label, l",
			Value: &cli.StringSlice{},
			Usage: "label to process, defaults to all",
		},
	}

	app.Action = func(ctx *cli.Context) {
		repo := ctx.String("repo")

		if repo == "" {
			fmt.Println("\n***** Missing required flag --repo *****\n")
			cli.ShowAppHelp(ctx)
			return
		}

		fmt.Printf("Starting Processing for repo '%s'\n", repo)
	}

	app.Run(os.Args)
}
