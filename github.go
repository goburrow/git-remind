package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	githubURL = "https://api.github.com"
)

type GitHubRepository struct {
	name string

	url    string
	client http.Client
}

func NewGitHubRepository(name string) *GitHubRepository {
	if name == "" {
		log.Fatal("empty git repository name")
	}
	return &GitHubRepository{
		url:  githubURL,
		name: name,
	}
}

func (r *GitHubRepository) SetURL(url string, insecure bool) {
	r.url = url
	if insecure {
		r.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: insecure,
			},
		}
	} else {
		r.client.Transport = nil
	}
}

func (r *GitHubRepository) PullRequests() []*PullRequest {
	url := fmt.Sprintf("%s/repos/%s/pulls", r.url, r.name)
	log.Println("getting PRs at", url)
	resp, err := r.client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected github status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	pulls := make([]GitHubPulls, 0, 5)
	err = decoder.Decode(&pulls)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("%+v\n", pulls)

	pr := make([]*PullRequest, len(pulls))
	for i := range pulls {
		p := &pulls[i]
		pr[i] = &PullRequest{
			URL:         p.HTMLURL,
			Title:       p.Title,
			Author:      p.User.Login,
			CreatedTime: p.CreatedAt,
		}
	}
	return pr
}

type GitHubPulls struct {
	HTMLURL string `json:"html_url"`
	Title   string `json:"title"`
	User    struct {
		Login string `json:"login"`
	} `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
