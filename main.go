package main

import (
	"flag"
	"time"
)

type PullRequest struct {
	URL         string
	Title       string
	Author      string
	CreatedTime time.Time
}

type Repository interface {
	PullRequests() []*PullRequest
}

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
	if c.HipChat.Token != "" {
		hipchat.SetToken(c.HipChat.Token)
	}
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
