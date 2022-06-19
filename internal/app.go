package internal

import (
	"github.com/dewep-online/goppy"
	"github.com/dewep-online/goppy/plugins/http"
	"github.com/dewep-online/solocms/internal/controllers"
)

func InitApp(conf string) {
	app := goppy.New()
	app.WithConfig(conf)
	app.Plugins(
		http.WithHTTP(),
	)
	app.Plugins(controllers.Modules...)
	app.Run()
}
