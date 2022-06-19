package files

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Meta struct {
	id, root string
	Name     string            `json:"name"`
	URI      string            `json:"uri"`
	Template string            `json:"tmpl"`
	Titles   map[string]string `json:"titles"`
}

func NewID() string {
	return uuid.NewString()
}

func NewMeta(id, root string) *Meta {
	return &Meta{
		id:   id,
		root: root,
	}
}

func (v *Meta) Root(path string) {
	v.root = path
}

func (v *Meta) ID(id string) {
	v.id = id
}

func (v *Meta) PathPage(lang string) string {
	return fmt.Sprintf("%s/%s_%s%s", v.root, v.id, lang, ExtPage)
}

func (v *Meta) PathPublic(lang string) string {
	return fmt.Sprintf("%s/%s_%s%s", v.root, v.id, lang, ExtPagePublic)
}

func (v *Meta) path() string {
	return fmt.Sprintf("%s/%s%s", v.root, v.id, ExtPageMeta)
}

func (v *Meta) ReadPublic(lang string) ([]byte, error) {
	if err := v.validate(); err != nil {
		return nil, err
	}
	return os.ReadFile(v.PathPublic(lang))
}

func (v *Meta) ReadPage(lang string) ([]byte, error) {
	if err := v.validate(); err != nil {
		return nil, err
	}
	return os.ReadFile(v.PathPage(lang))
}

func (v *Meta) SavePage(lang string, b []byte) error {
	if err := v.validate(); err != nil {
		return err
	}
	return os.WriteFile(v.PathPage(lang), b, 0666)
}

func (v *Meta) Update(b []byte) error {
	if err := v.validate(); err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func (v *Meta) Save() error {
	if err := v.validate(); err != nil {
		return err
	}
	return JSONEnc(v.path(), v)
}

func (v *Meta) Open() error {
	if err := v.validate(); err != nil {
		return err
	}
	return JSONDec(v.path(), v)
}

func (v *Meta) validate() error {
	if len(v.id) == 0 || len(v.root) == 0 {
		return fmt.Errorf("invalid meta setting: set ID and ROOT")
	}
	return nil
}
