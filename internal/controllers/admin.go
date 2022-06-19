package controllers

import (
	"fmt"
	nethttp "net/http"

	"github.com/dewep-online/solocms-admin-ui/guic"

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
	v.route.Use(
		BasicAuthMiddleware(v.conf.AdminAuth),
	)

	api := v.route.Collection("/api")
	api.Get(`/publish/{name:\w+}`, v.PublishPage)
	api.Get(`/page/{name:\w+}/{lang:[a-z]+}`, v.GetPage)
	api.Post(`/page/{name:\w+}/{lang:[a-z]+}`, v.SetPage)
	api.Get(`/meta/{name:\w+}`, v.GetMeta)
	api.Post(`/meta/{name:\w+}`, v.SetMeta)

	v.route.NotFoundHandler(v.UI)

	return nil
}

func (v *AdminCtrl) Down() error {
	return nil
}

func (v *AdminCtrl) UI(ctx http.Ctx) {
	guic.ResponseWrite(ctx.Response(), ctx.Request()) //nolint:errcheck
}

func (v *AdminCtrl) PublishPage(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "text/plain; charset=utf-8")

	name, _ := ctx.Param("name").String()

	switch name {
	case "all":
		if err := v.data.PublishAll(); err != nil {
			ctx.SetBody(nethttp.StatusInternalServerError).Error(err)
			return
		}
	default:
		if err := v.data.PublishByName(name); err != nil {
			ctx.SetBody(nethttp.StatusInternalServerError).Error(err)
			return
		}
	}

	if err := v.data.Sync(); err != nil {
		ctx.SetBody(nethttp.StatusInternalServerError).Error(err)
		return
	}

	ctx.SetBody(nethttp.StatusOK).Raw([]byte("ok"))
}

func (v *AdminCtrl) GetPage(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "text/plain; charset=utf-8")

	name, _ := ctx.Param("name").String()
	lang, _ := ctx.Param("lang").String()

	if !v.conf.HasLang(lang) {
		ctx.SetBody(nethttp.StatusBadRequest).Error(
			fmt.Errorf("lang `%s` is not supported", lang),
		)
		return
	}

	meta, err := v.data.GetPageMeta(name)
	if err != nil {
		ctx.SetBody(nethttp.StatusBadRequest).Error(err)
		return
	}

	b, err := meta.ReadPage(lang)
	if err != nil {
		ctx.SetBody(nethttp.StatusBadRequest).Error(err)
		return
	}

	ctx.SetBody(nethttp.StatusOK).Raw(b)
}

func (v *AdminCtrl) SetPage(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "text/plain; charset=utf-8")

	name, _ := ctx.Param("name").String()
	lang, _ := ctx.Param("lang").String()

	if !v.conf.HasLang(lang) {
		ctx.SetBody(nethttp.StatusBadRequest).Error(
			fmt.Errorf("lang `%s` is not supported", lang),
		)
		return
	}

	meta, err := v.data.GetPageMeta(name)
	if err != nil {
		ctx.SetBody(nethttp.StatusBadRequest).Error(err)
		return
	}

	if err = meta.SavePage(lang, ctx.GetBody().Raw()); err != nil {
		ctx.SetBody(nethttp.StatusBadRequest).Error(err)
		return
	}

	ctx.SetBody(nethttp.StatusOK).String("ok")
}

func (v *AdminCtrl) GetMeta(ctx http.Ctx) {
	name, _ := ctx.Param("name").String()

	meta, err := v.data.GetPageMeta(name)
	if err != nil {
		ctx.SetBody(nethttp.StatusBadRequest).Error(err)
		return
	}

	ctx.SetBody(nethttp.StatusOK).JSON(meta)
}

func (v *AdminCtrl) SetMeta(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "text/plain; charset=utf-8")

	name, _ := ctx.Param("name").String()

	meta, err := v.data.GetPageMeta(name)
	if err != nil {
		ctx.SetBody(nethttp.StatusBadRequest).Error(err)
		return
	}

	if err = meta.Update(ctx.GetBody().Raw()); err != nil {
		ctx.SetBody(nethttp.StatusBadRequest).Error(err)
		return
	}

	if err = meta.Save(); err != nil {
		ctx.SetBody(nethttp.StatusBadRequest).Error(err)
		return
	}

	ctx.SetBody(nethttp.StatusOK).String("ok")
}
