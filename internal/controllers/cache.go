package controllers

import (
	"fmt"
	"os"
	"sync"
	"text/template"

	"github.com/dewep-online/solocms/internal/pkg/files"
	"github.com/dewep-online/solocms/internal/pkg/render"
)

type Content struct {
	conf   *Config
	tmpl   map[string]*template.Template
	pages  map[string]*files.Meta
	routes map[string][]byte

	muxT sync.RWMutex
	muxP sync.RWMutex
	muxR sync.RWMutex
}

func NewContent(conf *Config) *Content {
	return &Content{
		conf:   conf,
		tmpl:   make(map[string]*template.Template),
		pages:  make(map[string]*files.Meta),
		routes: make(map[string][]byte),
	}
}

func (v *Content) Up() error {
	return v.Sync()
}

func (v *Content) Down() error {
	return nil
}

func (v *Content) Sync() error {
	if err := v.loadTemplate(); err != nil {
		return err
	}
	if err := v.loadPagesMeta(); err != nil {
		return err
	}
	if err := v.loadRoutes(); err != nil {
		return err
	}
	return nil
}

func (v *Content) getPathPage(lang, name string) string {
	return fmt.Sprintf("%s/%s_%s%s", v.conf.Content.Path, name, lang, files.ExtPage)
}

func (v *Content) getPathPageMeta(name string) string {
	return fmt.Sprintf("%s/%s%s", v.conf.Content.Path, name, files.ExtPageMeta)
}

func (v *Content) getPathPagePublic(lang, name string) string {
	return fmt.Sprintf("%s/%s_%s%s", v.conf.Content.Path, name, lang, files.ExtPagePublic)
}

func (v *Content) PublishAll() error {
	return v.EachPageMeta(v.publish)
}

func (v *Content) PublishByName(name string) error {
	meta, err := v.GetPageMeta(name)
	if err != nil {
		return err
	}
	return v.publish(name, meta)
}

func (v *Content) publish(name string, meta *files.Meta) error {
	tmpl, err := v.GetTemplate(meta.Template)
	if err != nil {
		return err
	}
	model := render.TemplateModel{CDN: v.conf.CDN.Domain}
	for lang, title := range meta.Titles {
		model.Title = title
		model.Lang = lang
		if err = render.Build(v.getPathPage(lang, name), v.getPathPagePublic(lang, name), model, tmpl); err != nil {
			return err
		}
	}
	return nil
}

func (v *Content) loadTemplate() error {
	return files.WalkDir(v.conf.Content.Path, files.ExtTmpl, func(filename, name string) error {
		t, err := render.NewTemplate(filename)
		if err != nil {
			return err
		}
		v.muxT.Lock()
		v.tmpl[name] = t
		v.muxT.Unlock()
		return nil
	})
}

func (v *Content) GetTemplate(name string) (*template.Template, error) {
	v.muxT.RLock()
	defer v.muxT.RUnlock()

	if t, ok := v.tmpl[name]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("template `%s` is not found", name)
}

func (v *Content) loadPagesMeta() error {
	return files.WalkDir(v.conf.Content.Path, files.ExtPageMeta, func(filename, name string) error {
		model := &files.Meta{}

		if err := files.JSONDec(filename, model); err != nil {
			return err
		}
		v.muxP.Lock()
		v.pages[name] = model
		v.muxP.Unlock()
		return nil
	})
}

func (v *Content) GetPageMeta(name string) (*files.Meta, error) {
	v.muxP.RLock()
	defer v.muxP.RUnlock()

	if meta, ok := v.pages[name]; ok {
		return meta, nil
	}
	return nil, fmt.Errorf("page `%s` is not found", name)
}

func (v *Content) EachPageMeta(call func(string, *files.Meta) error) error {
	v.muxP.RLock()
	defer v.muxP.RUnlock()

	for name, meta := range v.pages {
		if err := call(name, meta); err != nil {
			return err
		}
	}
	return nil
}

func (v *Content) loadRoutes() error {
	v.muxR.Lock()
	defer v.muxR.Unlock()

	return v.EachPageMeta(func(name string, meta *files.Meta) error {
		for lang, _ := range meta.Titles {
			filename := v.getPathPagePublic(lang, name)
			b, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("page `%s` load err: %w", name, err)
			}
			v.routes[fmt.Sprintf("/%s%s", lang, meta.URI)] = b
		}
		return nil
	})
}

func (v *Content) GetRoute(uri string) ([]byte, error) {
	v.muxR.RLock()
	defer v.muxR.RUnlock()

	if b, ok := v.routes[uri]; ok {
		return b, nil
	}
	return nil, fmt.Errorf("route `%s` is not found", uri)
}
