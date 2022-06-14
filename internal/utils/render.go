package utils

import (
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown"
	"html/template"
	"os"
)

type RenderModel struct {
	CDN     string
	Content string
}

func GetTemplate(filename string) (*template.Template, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file `%s`: %w", filename, err)
	}
	t, err := template.New("t").Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("parse template `%s`: %w", filename, err)
	}
	return t, nil
}

func Build(filename string, m *RenderModel, t *template.Template) ([]byte, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file `%s`: %w", filename, err)
	}

	m.Content = string(markdown.ToHTML(b, nil, nil))

	buf := &bytes.Buffer{}
	err = t.Execute(buf, m)
	if err != nil {
		return nil, fmt.Errorf("render template `%s`: %w", filename, err)
	}

	return buf.Bytes(), nil
}
