package render

import (
	"fmt"
	"os"
	"text/template"

	"github.com/gomarkdown/markdown"
)

type TemplateModel struct {
	CDN     string
	Content string
	Lang    string
	Title   string
}

func NewTemplate(filename string) (*template.Template, error) {
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

func Build(from, to string, m TemplateModel, t *template.Template) error {
	b, err := os.ReadFile(from)
	if err != nil {
		return fmt.Errorf("read file `%s`: %w", from, err)
	}

	m.Content = string(markdown.ToHTML(b, nil, nil))

	buf, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("create file `%s`: %w", to, err)
	}
	defer buf.Close()

	err = t.Execute(buf, m)
	if err != nil {
		return fmt.Errorf("render template `%s`: %w", from, err)
	}

	return nil
}
