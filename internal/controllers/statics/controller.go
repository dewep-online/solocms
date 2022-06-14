package statics

import (
	"fmt"
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
		route: pool.Main(),
	}
}

func (v *Controller) Up() error {

	v.route.Get("/", v.Cache)
	v.route.Get("/robots.txt", v.Robots)
	v.route.Get("/sitemap.xml", v.Sitemap)

	return nil
}

func (v *Controller) Down() error {
	return nil
}

const robotsTxt = `User-agent: *

Clean-param: utm_source&utm_content&utm_term&utm_medium&utm_campaign /

Sitemap: %s://%s/sitemap.xml
`

func (v *Controller) Robots(ctx http.Ctx) {
	uri := ctx.URL()
	if len(uri.Scheme) == 0 {
		uri.Scheme = "http"
	}
	ctx.SetHead("Content-Type", "text/plain; charset=utf-8")
	ctx.SetBody().Raw([]byte(
		fmt.Sprintf(robotsTxt, uri.Scheme, uri.Host),
	))
}

func (v *Controller) Sitemap(ctx http.Ctx) {
	ctx.SetHead("Content-Type", "application/xml")
}

func (v *Controller) Cache(ctx http.Ctx) {

}
