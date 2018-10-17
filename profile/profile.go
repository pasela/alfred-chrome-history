package profile

import (
	"path/filepath"
	"strings"

	"github.com/pasela/alfred-chrome-history/utils"
)

const (
	chromeConfigPath = "~/Library/Application Support/Google/Chrome"
)

func GetProfilePath(profile string) (string, error) {
	var path string
	if strings.Contains(profile, "/") {
		path = filepath.Clean(profile)
	} else {
		if profile == "" {
			profile = "Default"
		}
		path = filepath.Join(chromeConfigPath, profile)
	}

	return utils.ExpandTilde(path)
}
