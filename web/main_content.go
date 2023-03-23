package web

import (
	. "github.com/kiprotect/gospel"
)

func Projects(c Context) Element {
	return Div("heydo")
}

func MainContent(c Context) Element {

	// get the router
	router := UseRouter(c)
	// get the logged in user
	user := UseUser(c)

	return Div(
		Class("bulma-container"),
		Span("You are logged in as user ", Strong(user.Email()), ", nice!"),
		router.Match("/projects", Projects),
		Div(
		// c.Element("kodex", Kodex),
		// c.Element("userForm", UserForm),
		),
	)
}
