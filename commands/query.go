package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	aw "github.com/deanishe/awgo"
	"github.com/pasela/alfred-chrome-history/history"
	"github.com/pasela/alfred-chrome-history/utils"
	"github.com/spf13/cobra"
)

type queryOptions struct {
	AlfredFlag bool
}

var queryOpts queryOptions

func init() {
	queryCmd.Flags().BoolVar(&queryOpts.AlfredFlag, "alfred", false, "Alfred script filter mode")
	rootCmd.AddCommand(queryCmd)
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query Chrome history",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if queryOpts.AlfredFlag {
			return runQueryCmdAlfred(cmd, args)
		}
		return runQueryCmd(cmd, args)
	},
}

func runQueryCmd(cmd *cobra.Command, args []string) error {
	query := args[0]
	entries, err := runQuery(query, query)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fmt.Println(entry)
	}
	return nil
}

func runQueryCmdAlfred(cmd *cobra.Command, args []string) error {
	wf := aw.New()
	wf.Run(func() {
		query := args[0]
		entries, err := runQuery(query, query)
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
	})
	return nil
}

func runQuery(url, title string) ([]history.Entry, error) {
	histFile := historyFile{
		Profile: globalOpts.Profile,
		Clone:   false,
	}
	defer histFile.Close()

	filePath, err := histFile.GetPath()
	if err != nil {
		return nil, err
	}

	his, err := history.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer his.Close()

	return his.Query(url, title)
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
