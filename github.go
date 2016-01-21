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

	minAge time.Duration
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

func (r *GitHubRepository) SetMinAge(minAge time.Duration) {
	r.minAge = minAge
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

	pr := make([]*PullRequest, 0, len(pulls))
	for i := range pulls {
		p := &pulls[i]
		if !r.shouldRemind(p) {
			log.Printf("skipped repo %s", p.HTMLURL)
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

func (r *GitHubRepository) shouldRemind(p *GitHubPulls) bool {
	return r.minAge == 0 || time.Now().Sub(p.CreatedAt) > r.minAge
}

type GitHubPulls struct {
	HTMLURL string `json:"html_url"`
	Title   string `json:"title"`
	User    struct {
		Login string `json:"login"`
	} `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
