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
	URL          string
	Insecure     bool
	Token        string
	Repositories []string

	Filter struct {
		MinAge         Duration
		IgnoreAssigned bool
	}
}

// PullRequests returns all pull requests of added repositories.
func (r *GitHubRepository) PullRequests() []*PullRequest {
	baseURL := r.URL
	if baseURL == "" {
		baseURL = githubURL
	}
	client := http.Client{}
	if r.Insecure {
		// Skip verifying GitHub server certificate
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	var pulls []*PullRequest
	for _, n := range r.Repositories {
		url := fmt.Sprintf("%s/repos/%s/pulls", baseURL, n)
		pulls = append(pulls, r.pullRequests(&client, url)...)
	}
	return pulls
}

func (r *GitHubRepository) pullRequests(client *http.Client, url string) []*PullRequest {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	if r.Token != "" {
		req.Header.Set("Authorization", "token "+r.Token)
	}

	log.Println("getting pull requests from", url)
	resp, err := client.Do(req)
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
	return (r.Filter.MinAge.Duration > 0 && time.Now().Sub(p.CreatedAt) < r.Filter.MinAge.Duration) ||
		(r.Filter.IgnoreAssigned && p.Assignee.Login != "")
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
