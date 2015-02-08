package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

// Issue represents a simplified GitHub issue
type Issue struct {
	Number   int
	ClosedAt *time.Time
	Title    string
	Assignee string
	Labels   []string
	State    string
}

func newIssueFromEvent(event *github.IssueEvent) Issue {
	var assignee string
	if event.Issue.Assignee != nil {
		assignee = *event.Issue.Assignee.Login
	}

	var labels []string
	for _, label := range event.Issue.Labels {
		labels = append(labels, *label.Name)
	}

	return Issue{
		Number:   *event.Issue.Number,
		Title:    *event.Issue.Title,
		Assignee: assignee,
		ClosedAt: event.Issue.ClosedAt,
		Labels:   labels,
		State:    *event.Issue.State,
	}
}

func (issue *Issue) String() string {
	var closedAt string

	if issue.ClosedAt != nil {
		closedAt = issue.ClosedAt.In(time.Local).Format(dateFormat)
	}

	var assignee string
	if issue.Assignee != "" {
		assignee = " by @" + issue.Assignee
	}

	labelsString := strings.Join(issue.Labels, ", ")

	if len(labelsString) > 0 {
		labelsString = " (" + labelsString + ")"
	}

	return fmt.Sprintf("#%d - %s - %s%s%s", issue.Number, closedAt, issue.Title, assignee, labelsString)
}

func (issue *Issue) MatchLabels(targetLabels []string) bool {
	for _, issueLabel := range issue.Labels {
		for _, targetLabel := range targetLabels {
			if issueLabel == targetLabel {
				return true
			}
		}
	}
	return false
}
