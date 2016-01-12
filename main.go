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

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "~/.git/remind", "Config file path")
	flag.Parse()

	config := LoadConfig(configPath)

	github := NewGitHubRepository()
	github.AddRepos(config.GitHub.Repositories)
	if config.GitHub.URL != "" {
		github.SetURL(config.GitHub.URL, config.GitHub.Insecure)
	}
	pr := github.PullRequests()

	hipchat := NewHipChatReminder(config.HipChat.Room)
	if config.HipChat.Token != "" {
		hipchat.SetToken(config.HipChat.Token)
	}
	hipchat.Remind(pr)
}
