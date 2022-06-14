package internal

import (
	"github.com/dewep-online/goppy"
	"github.com/dewep-online/goppy/plugins/database"
	"github.com/dewep-online/goppy/plugins/http"
	"github.com/dewep-online/solocms/internal/controllers/admin"
	"github.com/dewep-online/solocms/internal/controllers/statics"
)

func IntApp(conf string) {

	app := goppy.New()
	app.WithConfig(conf)
	app.Plugins(
		http.WithHTTP(),
		database.WithSQLite(),
	)
	app.Plugins(statics.Module, admin.Module)
	app.Run()

}
