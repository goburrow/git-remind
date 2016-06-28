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

// GitHubRepository retrieves pull requests for given repositories.
type GitHubRepository struct {
	URL   string
	Token string

	MinAge         time.Duration
	IgnoreAssigned bool

	repos  []string
	client http.Client
}

// NewGitHubRepository initializes a new GitHubRepository with default GitHub URL.
func NewGitHubRepository() *GitHubRepository {
	return &GitHubRepository{
		URL: githubURL,
	}
}

// AddRepo registers repository for checking pull requests.
func (r *GitHubRepository) AddRepo(name string) {
	if name == "" {
		log.Fatal("empty git repository name")
	}
	r.repos = append(r.repos, name)
}

// Insecure skips verify GitHub server certificate.
func (r *GitHubRepository) Insecure() {
	r.client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

// PullRequests returns all pull requests of added repositories.
func (r *GitHubRepository) PullRequests() []*PullRequest {
	var pulls []*PullRequest
	for _, n := range r.repos {
		pulls = append(pulls, r.pullRequests(n)...)
	}
	return pulls
}

func (r *GitHubRepository) pullRequests(name string) []*PullRequest {
	url := r.URL
	if url == "" {
		url = githubURL
	}
	url = fmt.Sprintf("%s/repos/%s/pulls", url, name)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	if r.Token != "" {
		req.Header.Set("Authorization", "token "+r.Token)
	}

	log.Println("getting pull requests from", url)
	resp, err := r.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected github status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	pulls := make([]gitHubPulls, 0, 5)
	err = decoder.Decode(&pulls)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("%+v\n", pulls)

	pr := make([]*PullRequest, 0, len(pulls))
	for i := range pulls {
		p := &pulls[i]
		if r.shouldIgnore(p) {
			log.Printf("skipped pull request %s", p.HTMLURL)
			continue
		}
		pr = append(pr, &PullRequest{
			URL:         p.HTMLURL,
			Title:       p.Title,
			Author:      p.User.Login,
			CreatedTime: p.CreatedAt,
		})
	}
	return pr
}

func (r *GitHubRepository) shouldIgnore(p *gitHubPulls) bool {
	return (r.MinAge > 0 && time.Now().Sub(p.CreatedAt) < r.MinAge) || (r.IgnoreAssigned && p.Assignee.Login != "")
}

type gitHubPulls struct {
	HTMLURL string `json:"html_url"`
	Title   string `json:"title"`
	User    struct {
		Login string `json:"login"`
	} `json:"user"`
	Assignee struct {
		Login string `json:"login"`
	} `json:"assignee"`
	CreatedAt time.Time `json:"created_at"`
}
