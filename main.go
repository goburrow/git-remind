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

func newRepository(c *Config) Repository {
	github := NewGitHubRepository()
	for _, url := range c.GitHub.Repositories {
		github.AddRepo(url)
	}
	if c.GitHub.URL != "" {
		github.URL = c.GitHub.URL
	}
	if c.GitHub.Insecure {
		github.Insecure()
	}
	github.Token = c.GitHub.Token
	github.MinAge = c.GitHub.Filter.MinAge.Duration
	return github
}

func newReminder(c *Config) Reminder {
	hipchat := NewHipChatReminder(c.HipChat.Room)
	hipchat.Token = c.HipChat.Token
	return hipchat
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
	repository := newRepository(config)
	reminder := newReminder(config)
	reminder.Remind(repository.PullRequests())
}
