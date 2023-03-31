package web

import (
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
)

func ActionEditor(action kodex.ActionConfig) ElementFunction {
	return func(c Context) Element {

		kodex.Log.Infof("Config data: %v", action.ConfigData())

		return Div(
			"test",
		)
	}
}
