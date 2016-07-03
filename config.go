package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// Config is the model of application configuration.
type Config struct {
	GitHub  GitHubRepository
	HipChat HipChatReminder
}

// LoadConfig returns a new Config read from given file path.
func LoadConfig(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	decoder := json.NewDecoder(f)

	c := &Config{}
	err = decoder.Decode(c)
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// Duration wraps time.Duration to support unmarshalling from JSON number and string.
type Duration struct {
	time.Duration
}

// UnmarshalJSON parses duration from either number or string.
func (d *Duration) UnmarshalJSON(b []byte) error {
	// Check if it's a string
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		val, err := json.Number(string(b)).Int64()
		if err != nil {
			return err
		}
		d.Duration = time.Duration(val)
		return nil
	}
	d.Duration, err = time.ParseDuration(s)
	return err
}
