package utils

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

func CopyFile(src, dst string) (int64, error) {
	srcStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !srcStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()

	df, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return 0, err
	}
	defer df.Close()

	return io.Copy(df, sf)
}

func ExpandTilde(path string) (string, error) {
	if path[0] != '~' {
		return path, nil
	}

	si := strings.IndexAny(path, "/\\")
	var head, tail string
	if si == -1 {
		head, tail = path, ""
	} else {
		head, tail = path[:si], path[si:]
	}

	var usr *user.User
	var err error
	if head == "~" {
		usr, err = user.Current()
	} else {
		usr, err = user.Lookup(head[1:])
	}
	if err != nil {
		return path, err
	}

	return usr.HomeDir + tail, nil
}
