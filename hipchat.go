package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

const (
	hipchatURL = "https://api.hipchat.com/v2"
)

// HipChatReminder reminds pull requests on a HipChat room.
type HipChatReminder struct {
	Token string

	url  string
	room string

	client http.Client
}

// NewHipChatReminder initializes a new HipChatReminder with given room.
func NewHipChatReminder(room string) *HipChatReminder {
	if room == "" {
		log.Fatal("empty hipchat room")
	}
	return &HipChatReminder{
		url:  hipchatURL,
		room: room,
	}
}

// Remind sends a notification for given pull requests.
func (r *HipChatReminder) Remind(pr []*PullRequest) {
	if len(pr) == 0 {
		return
	}

	now := time.Now()
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Reminders (%d)", len(pr))
	for _, p := range pr {
		fmt.Fprintf(&buf, " <a href=%q>%s</a> (%s %s ago)",
			p.URL, html.EscapeString(p.Title), p.Author,
			(now.Sub(p.CreatedTime)/time.Second)*time.Second)
	}

	notif := make(map[string]string)
	notif["from"] = "GitReminder"
	notif["message"] = buf.String()

	buf.Reset()
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(notif)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/room/%s/notification", r.url, r.room)
	log.Println("sending notification to", url)
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if r.Token != "" {
		req.Header.Set("Authorization", "Bearer "+r.Token)
	}
	res, err := r.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusNoContent {
		log.Fatalf("unexpected hipchat status code: %d", res.StatusCode)
	}
}
