package utils

import (
	"io/fs"
	"path/filepath"
	"strings"
)

const (
	ExtPage = ".md"
	ExtTmpl = ".tmpl"
)

func WalkDir(dir, ext string, call func(filename, uri string) error) error {
	root := len(strings.TrimRight(dir, "/"))
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ext {
			return nil
		}
		return call(path, path[root:len(path)-3])
	})
}
