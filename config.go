package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
	GitHub struct {
		URL          string
		Insecure     bool
		Repositories []string

		Filter struct {
			MinAge Duration
		}
	}
	HipChat struct {
		URL   string
		Token string
		Room  string
	}
}

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

type Duration struct {
	time.Duration
}

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
