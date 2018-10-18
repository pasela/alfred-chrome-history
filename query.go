package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pasela/alfred-chrome-history/history"
	"github.com/pasela/alfred-chrome-history/utils"
)

func queryHistory(url, title string) ([]history.Entry, error) {
	histFile := historyFile{
		Profile: profile,
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
