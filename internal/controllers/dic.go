package controllers

import "github.com/dewep-online/goppy/plugins"

var Modules = []plugins.Plugin{
	{Config: &Config{}, Resolve: NewContent},
	{Resolve: NewStaticCtrl},
	{Resolve: NewAdminCtrl},
}
