package main

import (
	"log"
	"os"

	aw "github.com/deanishe/awgo"
)

var isAlfred bool
var profile string

func init() {
	isAlfred = os.Getenv("alfred_workflow_bundleid") != ""
}

func main() {
	if isAlfred {
		wf := aw.New()
		wf.Run(func() {
			runWithAlfred(wf)
		})
	} else {
		if err := run(); err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}
}
