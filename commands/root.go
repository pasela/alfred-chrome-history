package commands

import (
	"os"

	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	Profile string
}

var globalOpts GlobalOptions

func init() {
	rootCmd.PersistentFlags().StringVar(&globalOpts.Profile, "profile", os.Getenv("CHROME_PROFILE"), "Chrome profile path")
}

var rootCmd = &cobra.Command{
	Use:   "alfred-chrome-history",
	Short: "alfred-chrome-history is a query tool for Chrome history",
}

func Execute() int {
	if err := rootCmd.Execute(); err != nil {
		return 1
	}

	return 0
}
