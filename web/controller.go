package web

import (
	"github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex/api"
)

func SetController(c gospel.Context, controller api.Controller) {
	gospel.GlobalVar(c, "controller", controller)
}

func UseController(c gospel.Context) api.Controller {
	return gospel.UseGlobal[api.Controller](c, "controller")
}
