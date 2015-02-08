package main

import (
	"fmt"
	"strings"
)

// Repo represents a GitHub repo
type Repo struct {
	Owner string
	Name  string
}

// NewRepoFromPath initializes a Repo by a path in the form "owner/name"
func NewRepoFromPath(path string) (Repo, error) {
	// TODO validate and return error
	parts := strings.Split(path, "/")

	return Repo{Owner: parts[0], Name: parts[1]}, nil
}

func (repo *Repo) String() string {
	return fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
}
