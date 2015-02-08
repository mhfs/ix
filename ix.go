package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

const (
	dateFormat = "2006-01-02"
)

func main() {
	app := cli.NewApp()
	app.Name = "ix"
	app.Usage = "Issues Explorer - CLI tool to explore GitHub issues by repository, time frame, labels and assignee."
	app.Version = "0.0.1"
	app.Author = "Marcelo Silveira"
	app.Email = "marcelo@mhfs.com.br"

	app.Commands = []cli.Command{
		{
			Name:      "closed",
			ShortName: "c",
			Usage:     "lists closed issues",
			Action:    closedCommand,
			Flags: []cli.Flag{
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
			},
		},
	}

	app.Run(os.Args)
}

func closedCommand(ctx *cli.Context) {
	repoPath := ctx.String("repo")
	assignee := ctx.String("assignee")
	labels := ctx.StringSlice("label")
	token := ctx.String("token")
	since, err := time.ParseInLocation(dateFormat, ctx.String("since"), time.Local)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	repo, err := NewRepoFromPath(repoPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}
	client := t.Client()

	fmt.Printf("Analyzing repo '%s' since '%s'\n\n", repo.String(), since.Format(dateFormat))

	var done bool
	var issues []Issue

	for page := 1; !done; page++ {
		issues, done = fetchIssuesFromEvents(client, repo, since, "closed", page)

		if err != nil {
			fmt.Printf("Couldn't fetch issue from GitHub. Error: '%s'\n", err)
			os.Exit(1)
		}

		for _, issue := range issues {
			// do not show if issue is not closed
			if issue.State != "closed" {
				continue
			}

			if len(assignee) > 0 && issue.Assignee != assignee {
				continue
			}

			// if label filtering set, skip labels we're not interested at
			if len(labels) > 0 && !issue.MatchLabels(labels) {
				continue
			}

			fmt.Println(issue.String())
		}
	}
}

func fetchEvents(client *http.Client, repo Repo, page int) []github.IssueEvent {
	gh := github.NewClient(client)

	options := github.ListOptions{Page: page, PerPage: 100}
	events, _, err := gh.Issues.ListRepositoryEvents(repo.Owner, repo.Name, &options)

	if err != nil {
		fmt.Printf("Couldn't fetch events from GitHub. Error: '%s'\n", err)
		os.Exit(1)
	}

	return events
}

func fetchIssuesFromEvents(httpClient *http.Client, repo Repo, since time.Time, event string, page int) ([]Issue, bool) {
	done := false
	events := fetchEvents(httpClient, repo, page)

	var issues []Issue
	for _, event := range events {
		// events are ordered by created at desc. stop if got all we wanted.
		if event.CreatedAt.Before(since) {
			done = true
			break
		}

		// filter by desired state
		if *event.Event != "closed" {
			continue
		}

		issue := newIssueFromEvent(&event)
		issues = append(issues, issue)
	}

	return issues, done
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
