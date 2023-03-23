package web

import (
	"github.com/kiprotect/gospel"
)

func AppServer() *gospel.Server {
	return gospel.MakeServer(&gospel.App{
		Root:         Root,
		StaticFiles:  StaticFiles,
		StaticPrefix: "/static",
	})
}
