package main

import (
	"flag"
	"time"
)

// PullRequest is data structure returned by Repository
type PullRequest struct {
	URL         string
	Title       string
	Author      string
	CreatedTime time.Time
}

// Repository returns list of PullRequest.
type Repository interface {
	PullRequests() []*PullRequest
}

// Reminder sends a reminder with given pull requests.
type Reminder interface {
	Remind([]*PullRequest)
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "git-remind-config.json", "Config file path")
	flag.Parse()

	if configPath == "" {
		flag.Usage()
		return
	}
	config := LoadConfig(configPath)
	var repository Repository = &config.GitHub
	var reminder Reminder = &config.HipChat
	reminder.Remind(repository.PullRequests())
}
