package web

import (
	. "github.com/kiprotect/gospel"
)

func MainContent(c Context) Element {

	// get the router
	router := UseRouter(c)

	return Div(
		Class("bulma-container"),
		router.Match(
			Route("/projects/(?P<projectId>[^/]+)", ProjectDetails),
			Route("/projects", Projects),
		),
		Div(
		// c.Element("kodex", Kodex),
		// c.Element("userForm", UserForm),
		),
	)
}
