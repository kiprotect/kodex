package web

import (
	"github.com/gospel-sh/gospel"
)

func SetPlugins(c gospel.Context, plugins []WebPlugin) {
	gospel.GlobalVar(c, "plugins", plugins)
}

func UsePlugins(c gospel.Context) []WebPlugin {
	return gospel.UseGlobal[[]WebPlugin](c, "plugins")
}
