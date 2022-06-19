package main

import (
	"github.com/dewep-online/solocms/internal"
	"github.com/deweppro/go-app/console"
)

func main() {
	root := console.New("solocms", "help")
	root.AddCommand(appRun())
	root.Exec()
}

func appRun() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("run", "run application")
		setter.Example("run --config=./config.yaml")
		setter.Flag(func(f console.FlagsSetter) {
			f.StringVar("config", "./config.yaml", "path to config file")
		})
		setter.ExecFunc(func(_ []string, config string) {
			internal.InitApp(config)
		})
	})
}
