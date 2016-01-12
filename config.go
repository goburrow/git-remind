package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	GitHub struct {
		URL          string
		Insecure     bool
		Repositories []string
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
