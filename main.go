package main

import (
	"os"

	"github.com/pasela/alfred-chrome-history/commands"
)

func main() {
	code := commands.Execute()
	if code != 0 {
		os.Exit(code)
	}
}
