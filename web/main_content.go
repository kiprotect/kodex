package web

import (
	. "github.com/kiprotect/gospel"
)

func Projects(c Context) Element {
	return Div("heydo")
}

func MainContent(c Context) Element {

	router := UseRouter(c)

	return Div(
		Class("bulma-container"),
		router.Match("/projects", Projects),
		Div(
		// c.Element("kodex", Kodex),
		// c.Element("userForm", UserForm),
		),
	)
}
