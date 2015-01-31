package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
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

		parts := strings.Split(repo, "/")
		issues, err := loadIssues(parts[0], parts[1])

		if err != nil {
			panic(err) // FIXME be smarter. handle 404s, 403s, ...
		}

		for _, issue := range issues {
			printIssue(&issue)
		}
	}

	app.Run(os.Args)
}

func loadIssues(owner string, repo string) ([]github.Issue, error) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: os.Getenv("GH_TOKEN")},
	}

	client := github.NewClient(t.Client())

	from, err := time.Parse("2006-01-02", "2015-01-25")

	if err != nil {
		panic("invalid date provided")
	}

	options := github.IssueListByRepoOptions{State: "closed", Sort: "updated", Since: from}
	issues, _, err := client.Issues.ListByRepo(owner, repo, &options)

	return issues, err
}

func printIssue(issue *github.Issue) {
	// TODO handle nils and missing assignee
	// fmt.Printf("#%d %s %s (%s)", issue.Number, issue.ClosedAt, issue.Title, issue.Assignee.Login)
	fmt.Printf("#%d - %s - %s\n", *issue.Number, issue.ClosedAt, *issue.Title)
}
