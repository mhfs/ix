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

const (
	dateFormat = "2006-01-02"
)

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[options] {{end}}
"
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
			Name:  "since, s",
			Value: utcBeginningOfWeekFromLocal().In(time.Local).Format(dateFormat),
			Usage: "list issues since given date, inclusive",
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

		since, err := time.ParseInLocation(dateFormat, ctx.String("since"), time.Local)
		if err != nil {
			panic("invalid date provided")
		}

		fmt.Printf("Starting work for repo '%s' since '%s'\n", repo, since.Format(dateFormat))

		parts := strings.Split(repo, "/")
		issues, err := fetchIssues(parts[0], parts[1], since)

		if err != nil {
			panic(err) // FIXME be smarter. handle 404s, 403s, ...
		}

		for _, issue := range issues {
			printIssue(&issue)
		}
	}

	app.Run(os.Args)
}

func fetchIssues(owner string, repo string, since time.Time) ([]github.Issue, error) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: os.Getenv("GH_TOKEN")},
	}

	client := github.NewClient(t.Client())

	options := github.IssueListByRepoOptions{State: "closed", Sort: "updated", Since: since}
	issues, _, err := client.Issues.ListByRepo(owner, repo, &options)

	return issues, err
}

func printIssue(issue *github.Issue) {
	// TODO handle nils and missing assignee
	// fmt.Printf("#%d %s %s (%s)", issue.Number, issue.ClosedAt, issue.Title, issue.Assignee.Login)
	fmt.Printf("#%d - %s - %s\n", *issue.Number, issue.ClosedAt.Format(dateFormat), *issue.Title)
}

func utcBeginningOfWeekFromLocal() time.Time {
	now := time.Now()
	_, offset := now.Zone()
	beginningOfDay := now.UTC().Truncate(24 * time.Hour).Add(-1 * time.Duration(offset) * time.Second)
	weekFirstDay := beginningOfDay.Add(-time.Duration(now.Weekday()) * 24 * time.Hour)
	return weekFirstDay
}
