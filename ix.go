package main

import (
	"fmt"
	"net/http"
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

EXAMPLES:
   {{.Name}} --repo mhfs/ix --since 2015-01-01
   {{.Name}} --repo mhfs/ix --assignee mhfs
   {{.Name}} --repo mhfs/ix --label bug

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
	app.Name = "ix"
	app.Usage = "cli to explore closed GitHub issue for a repository by time frame, labels and assignee"
	app.Version = "0.0.1"
	app.Author = "Marcelo Silveira"
	app.Email = "marcelo@mhfs.com.br"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "repo, r",
			Value: "",
			Usage: "GitHub repository to analyze e.g. mhfs/ix",
		},
		cli.StringFlag{
			Name:  "since, s",
			Value: beginningOfWeek().In(time.Local).Format(dateFormat),
			Usage: "list issues since given date, inclusive",
		},
		cli.StringSliceFlag{
			Name:  "label, l",
			Value: &cli.StringSlice{},
			Usage: "label to process, defaults to all",
		},
		cli.StringFlag{
			Name:  "assignee, a",
			Value: "",
			Usage: "filter results by assignee",
		},
		cli.StringFlag{
			Name:   "token, t",
			Value:  "",
			Usage:  "oauth token. defaults to GH_TOKEN env var.",
			EnvVar: "GH_TOKEN",
		},
	}

	app.Action = mainCommand

	app.Run(os.Args)
}

func mainCommand(ctx *cli.Context) {
	repo := ctx.String("repo")
	assignee := ctx.String("assignee")
	labels := ctx.StringSlice("label")
	token := ctx.String("token")

	if repo == "" {
		fmt.Println("\n***** Missing required flag --repo *****\n")
		cli.ShowAppHelp(ctx)
		return
	}

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}
	client := t.Client()

	since, err := time.ParseInLocation(dateFormat, ctx.String("since"), time.Local)
	if err != nil {
		panic("invalid date provided")
	}

	fmt.Printf("Analyzing repo '%s' since '%s'\n\n", repo, since.Format(dateFormat))

	parts := strings.Split(repo, "/")
	owner, repo := parts[0], parts[1]

	var done bool
	for page := 1; !done; page++ {
		events, err := fetchEvents(client, owner, repo, page)

		if err != nil {
			fmt.Printf("Couldn't fetch issue from GitHub. Error: '%s'\n", err)
			os.Exit(1)
		}

		done = printEvents(events, since, labels, assignee)
	}
}

func fetchEvents(httpClient *http.Client, owner string, repo string, page int) ([]github.IssueEvent, error) {
	client := github.NewClient(httpClient)

	options := github.ListOptions{Page: page, PerPage: 100}
	events, _, err := client.Issues.ListRepositoryEvents(owner, repo, &options)

	return events, err
}

func printEvents(events []github.IssueEvent, since time.Time, labels []string, assignee string) bool {
	for _, event := range events {
		// events are ordered by created at desc. stop if got all we wanted.
		if event.CreatedAt.Before(since) {
			return true
		}

		// if event is closed or issue didn't remain closed
		if *event.Event != "closed" || *event.Issue.State != "closed" {
			continue
		}

		if len(assignee) > 0 && (event.Issue.Assignee == nil || assignee != *event.Issue.Assignee.Login) {
			continue
		}

		// if label filtering set, skip labels we're not interested at
		if len(labels) > 0 {
			if !matchingLabels(event.Issue.Labels, labels) {
				continue
			}
		}

		printEvent(&event)
	}

	return false
}

func printEvent(event *github.IssueEvent) {
	number := event.Issue.Number
	closedAt := event.Issue.ClosedAt.In(time.Local).Format(dateFormat)
	title := event.Issue.Title

	var assignee string
	if event.Issue.Assignee != nil {
		assignee = " by @" + *event.Issue.Assignee.Login
	}

	var labelsNames []string
	if labels := event.Issue.Labels; len(labels) > 0 {
		for _, l := range labels {
			labelsNames = append(labelsNames, *l.Name)
		}
	}
	labelsString := strings.Join(labelsNames, ", ")

	if len(labelsString) > 0 {
		labelsString = " (" + labelsString + ")"
	}

	fmt.Printf("#%d - %s - %s%s%s\n", *number, closedAt, *title, assignee, labelsString)
}

func beginningOfWeek() time.Time {
	now := time.Now()
	// truncate internal HH:MM:SS to zero and compensate for local zone offset
	// 2015-01-31 10:45:54.720292964 -0800 PST > 2015-01-30 16:00:00 -0800 PST > 2015-01-31 00:00:00 -0800 PST
	_, offset := now.Zone()
	beginningOfDay := now.Truncate(24 * time.Hour).Add(-1 * time.Duration(offset) * time.Second)

	// subtract days to get to sunday
	beginningOfWeek := beginningOfDay.Add(-time.Duration(now.Weekday()) * 24 * time.Hour)

	return beginningOfWeek
}

func matchingLabels(issueLabels []github.Label, targetLabels []string) bool {
	for _, issueLabel := range issueLabels {
		for _, targetLabel := range targetLabels {
			if *issueLabel.Name == targetLabel {
				return true
			}
		}
	}
	return false
}
