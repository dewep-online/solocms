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

func JSONEnc(filename string, model interface{}) error {
	b, err := json.Marshal(model)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0666)
}
