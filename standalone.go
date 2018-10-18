package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func run() error {
	profile := os.Getenv("CHROME_PROFILE")
	flag.StringVar(&profile, "profile", profile, "Chrome profile directory")
	flag.Parse()

	query := strings.Join(flag.Args(), " ")
	entries, err := queryHistory(profile, query, query)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
	return nil
}
