package controllers

import (
	"fmt"
	nethttp "net/http"
	"regexp"

	"github.com/dewep-online/goppy/middlewares"
	"github.com/dewep-online/goppy/plugins/http"
	"golang.org/x/text/language"
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

func (v *StaticCtrl) AllowDomainsMiddleware() middlewares.Middleware {
	domains := make(map[string]struct{})
	for _, domain := range v.conf.AllowDomains {
		domains[domain] = struct{}{}
	}

	return func(call func(nethttp.ResponseWriter, *nethttp.Request)) func(nethttp.ResponseWriter, *nethttp.Request) {
		return func(w nethttp.ResponseWriter, r *nethttp.Request) {
			if _, ok := domains[r.Host]; !ok {
				w.WriteHeader(nethttp.StatusForbidden)
				return
			}
			call(w, r)
		}
	}
}

func (v *StaticCtrl) LangMiddleware() middlewares.Middleware {
	langs := make(map[string]struct{})
	defaultLang := "en"
	for i, lang := range v.conf.Langs {
		if i == 0 {
			defaultLang = lang
		}
		langs[lang] = struct{}{}
	}

	rex := regexp.MustCompile(`^/([a-z]{2,3})(\/|$)`)

	return func(call func(nethttp.ResponseWriter, *nethttp.Request)) func(nethttp.ResponseWriter, *nethttp.Request) {
		return func(w nethttp.ResponseWriter, r *nethttp.Request) {

			rex.

			lang, uri := "", r.RequestURI

			if len(uri) >= 3 {
				lang, uri = uri[1:3], uri[3:]
			}

			if _, ok := langs[lang]; !ok {
				lang = defaultLang
				tags, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
				if err == nil {
					for _, tag := range tags {
						base, _ := tag.Base()
						if _, ok = langs[base.String()]; ok {
							lang = base.String()
							break
						}
					}
				}

				nethttp.Redirect(w, r, "/"+lang+uri, nethttp.StatusPermanentRedirect)
				return
			}

			call(w, r)
		}
	}
}

func (v *StaticCtrl) Up() error {

	v.route.Use(
		v.AllowDomainsMiddleware(),
		v.LangMiddleware(),
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
	ctx.SetBody(nethttp.StatusOK).Raw([]byte(
		fmt.Sprintf(robotsTxt, uri.Scheme, uri.Host),
	))
}

func (v *StaticCtrl) Sitemap(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "application/xml; charset=utf-8")
}

func (v *StaticCtrl) Cache(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "text/html; charset=utf-8")

	b, err := v.data.GetRoute(ctx.URL().Path)
	if err != nil {
		ctx.SetBody(nethttp.StatusNotFound).Raw([]byte(err.Error()))
		return
	}
	ctx.SetBody(nethttp.StatusOK).Raw(b)
}
