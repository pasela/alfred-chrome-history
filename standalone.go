package main

import (
	"fmt"
	"os"
	"strings"
)

func run() error {
	profile := os.Getenv("CHROME_PROFILE")
	query := strings.Join(os.Args[1:], " ")
	entries, err := queryHistory(profile, query, query)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
	return nil
}
