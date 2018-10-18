package main

import (
	"strconv"
	"strings"

	aw "github.com/deanishe/awgo"
)

func runWithAlfred(wf *aw.Workflow) {
	args := wf.Args()
	query := strings.Join(args, " ")
	entries, err := queryHistory(query, query)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		item := wf.NewItem(entry.Title)
		item.UID(strconv.Itoa(entry.ID))
		item.Subtitle(entry.URL)
		item.Arg(entry.URL)
		item.Valid(true)
	}
	wf.SendFeedback()
}
