package controllers

import (
	nethttp "net/http"

	"github.com/dewep-online/goppy/plugins/http"
)

type AdminCtrl struct {
	conf  *Config
	route http.Router
	data  *Content
}

func NewAdminCtrl(conf *Config, data *Content, pool http.RouterPool) *AdminCtrl {
	return &AdminCtrl{
		conf:  conf,
		route: pool.Get("admin"),
		data:  data,
	}
}

func (v *AdminCtrl) Up() error {

	v.route.Get("/publish/all", v.PublishAll)

	return nil
}

func (v *AdminCtrl) Down() error {
	return nil
}

func (v *AdminCtrl) PublishAll(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "text/plain; charset=utf-8")

	if err := v.data.PublishAll(); err != nil {
		ctx.SetBody(nethttp.StatusInternalServerError).Raw([]byte(err.Error()))
		return
	}

	if err := v.data.Sync(); err != nil {
		ctx.SetBody(nethttp.StatusInternalServerError).Raw([]byte(err.Error()))
		return
	}

	ctx.SetBody(nethttp.StatusOK).Raw([]byte("ok"))
}
