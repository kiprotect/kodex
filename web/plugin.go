package web

import (
	"github.com/gospel-dev/gospel"
	"github.com/kiprotect/kodex/api"
)

type WebPluginMaker interface {
	InitializeWebPlugin(controller api.Controller) (WebPlugin, error)
}

type WebPlugin interface {
	MainRoutes(gospel.Context) []*gospel.RouteConfig
}

type AppLinkPlugin interface {
	AppLink() (string, string, string)
}
