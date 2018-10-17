package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pasela/alfred-chrome-history/history"
	"github.com/pasela/alfred-chrome-history/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query Chrome history",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runQueryCmd,
}

func runQueryCmd(cmd *cobra.Command, args []string) error {
	histFile := historyFile{
		Profile: globalOpts.Profile,
		Clone:   false,
	}
	defer histFile.Close()

	filePath, err := histFile.GetPath()
	if err != nil {
		return err
	}

	his, err := history.Open(filePath)
	if err != nil {
		return err
	}
	defer his.Close()

	query := args[0]
	entries, err := his.Query(query, query)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}

	return nil
}

type historyFile struct {
	Profile string
	Clone   bool

	tempDir    string
	autoRemove bool
}

func (h *historyFile) initTempDir() error {
	if !h.Clone || h.tempDir != "" {
		return nil
	}

	tempDir, err := ioutil.TempDir("", "alfred-chrome-history")
	if err != nil {
		return err
	}
	h.tempDir = tempDir
	h.autoRemove = true
	return nil
}

func (h *historyFile) Close() error {
	if h.autoRemove && h.tempDir != "" {
		return os.RemoveAll(h.tempDir)
	}
	return nil
}

func (h *historyFile) GetPath() (string, error) {
	origFile, err := history.GetHistoryPath(h.Profile)
	if err != nil {
		return "", err
	}
	if !h.Clone {
		return origFile, nil
	}

	if err := h.initTempDir(); err != nil {
		return "", err
	}

	clonedFile := filepath.Join(h.tempDir, filepath.Base(origFile))
	if _, err := utils.CopyFile(origFile, clonedFile); err != nil {
		return "", err
	}
	return clonedFile, nil
}
