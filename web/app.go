package web

import (
	"github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex/api"
)

func AppServer(controller api.Controller) (*gospel.Server, error) {

	root, err := Root(controller)

	if err != nil {
		return nil, err
	}

	return gospel.MakeServer(&gospel.App{
		Root:         root,
		StaticFiles:  StaticFiles,
		StaticPrefix: "/static",
	}), nil
}
