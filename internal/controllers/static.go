package controllers

import (
	nethttp "net/http"

	"github.com/dewep-online/goppy/plugins/http"
)

type StaticCtrl struct {
	conf  *Config
	route http.Router
	data  *Content
}

func NewStaticCtrl(conf *Config, data *Content, pool http.RouterPool) *StaticCtrl {
	return &StaticCtrl{
		conf:  conf,
		route: pool.Main(),
		data:  data,
	}
}

func (v *StaticCtrl) Up() error {
	v.route.Use(
		AllowDomainsMiddleware(v.conf.AllowDomains),
		LangMiddleware(v.conf.Langs),
	)

	v.route.Get("/robots.txt", v.Robots)
	v.route.Get("/sitemap.xml", v.Sitemap)

	v.route.NotFoundHandler(v.Cache)

	return nil
}

func (v *StaticCtrl) Down() error {
	return nil
}

const robotsTxt = `User-agent: *

Clean-param: utm_source&utm_content&utm_term&utm_medium&utm_campaign /

Sitemap: %s://%s/sitemap.xml
`

func (v *StaticCtrl) Robots(ctx http.Ctx) {
	uri := ctx.URL()
	if len(uri.Scheme) == 0 {
		uri.Scheme = "http"
	}
	ctx.SetHead("Content-Type", "text/plain; charset=utf-8")
	ctx.SetBody(nethttp.StatusOK).String(robotsTxt, uri.Scheme, uri.Host)
}

func (v *StaticCtrl) Sitemap(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "application/xml; charset=utf-8")
}

func (v *StaticCtrl) Cache(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "text/html; charset=utf-8")

	b, err := v.data.GetRoute(ctx.URL().Path)
	if err != nil {
		ctx.SetBody(nethttp.StatusNotFound).String(err.Error())
		return
	}
	ctx.SetBody(nethttp.StatusOK).Raw(b)
}
