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
	repos []string

	url    string
	client http.Client
}

func NewGitHubRepository() *GitHubRepository {
	return &GitHubRepository{
		url: githubURL,
	}
}

func (r *GitHubRepository) AddRepo(name string) {
	if name == "" {
		log.Fatal("empty git repository name")
	}
	r.repos = append(r.repos, name)
}

func (r *GitHubRepository) AddRepos(names []string) {
	if len(names) == 0 {
		log.Fatal("empty git repository name")
	}
	for _, n := range names {
		r.AddRepo(n)
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
	var pulls []*PullRequest
	for _, n := range r.repos {
		pulls = append(pulls, r.pullRequests(n)...)
	}
	return pulls
}

func (r *GitHubRepository) pullRequests(name string) []*PullRequest {
	url := fmt.Sprintf("%s/repos/%s/pulls", r.url, name)
	log.Println("getting pull requests from", url)
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
