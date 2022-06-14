package admin

import (
	"github.com/dewep-online/goppy/plugins/database"
	"github.com/dewep-online/goppy/plugins/http"
)

type Controller struct {
	conf  *Config
	orm   database.SQLite
	route http.Router
}

func New(conf *Config, orm database.SQLite, pool http.RouterPool) *Controller {
	return &Controller{
		conf:  conf,
		orm:   orm,
		route: pool.Get("admin"),
	}
}

func (v *Controller) Up() error {
	return nil
}

func (v *Controller) Down() error {
	return nil
}
