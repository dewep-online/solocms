package controllers

import (
	"fmt"
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

func (v *Content) PublishAll() error {
	return v.EachPageMeta(v.publish)
}

func (v *Content) PublishByName(name string) error {
	meta, err := v.GetPageMeta(name)
	if err != nil {
		return err
	}
	return v.publish(meta)
}

func (v *Content) publish(meta *files.Meta) error {
	tmpl, err := v.GetTemplate(meta.Template)
	if err != nil {
		return err
	}
	model := render.TemplateModel{CDN: v.conf.CDN.Domain}
	for lang, title := range meta.Titles {
		model.Title = title
		model.Lang = lang
		if err = render.Build(meta.PathPage(lang), meta.PathPublic(lang), model, tmpl); err != nil {
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
	return files.WalkDir(v.conf.Content.Path, files.ExtPageMeta, func(filename, id string) error {
		model := files.NewMeta(id, v.conf.Content.Path)
		if err := model.Open(); err != nil {
			return err
		}
		v.muxP.Lock()
		v.pages[id] = model
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

func (v *Content) EachPageMeta(call func(*files.Meta) error) error {
	v.muxP.RLock()
	defer v.muxP.RUnlock()

	for _, meta := range v.pages {
		if err := call(meta); err != nil {
			return err
		}
	}
	return nil
}

func (v *Content) loadRoutes() error {
	v.muxR.Lock()
	defer v.muxR.Unlock()

	return v.EachPageMeta(func(meta *files.Meta) error {
		for lang, _ := range meta.Titles {
			b, err := meta.ReadPublic(lang)
			if err != nil {
				return fmt.Errorf("page `%s` load err: %w", meta.Name, err)
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
