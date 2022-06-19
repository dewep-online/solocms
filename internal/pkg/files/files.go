package files

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	ExtPage       = ".page"
	ExtPagePublic = ".html"
	ExtPageMeta   = ".meta"
	ExtTmpl       = ".tmpl"
)

type (
	Meta struct {
		Name     string            `json:"name"`
		URI      string            `json:"uri"`
		Template string            `json:"tmpl"`
		Titles   map[string]string `json:"titles"`
	}
)

func WalkDir(dir, ext string, call func(filename, name string) error) error {
	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ext {
			return nil
		}
		name := info.Name()
		return call(path, name[:len(name)-len(ext)])
	})
}

func JSONDec(filename string, model interface{}) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, model)
}
