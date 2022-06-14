package admin

import "github.com/dewep-online/goppy/plugins"

var Module = plugins.Plugin{
	Config:  &Config{},
	Resolve: New,
}
