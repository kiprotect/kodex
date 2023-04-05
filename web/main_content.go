package web

import (
	. "github.com/kiprotect/gospel"
)

func MainContent(c Context) Element {

	// get the router
	router := UseRouter(c)

	return Div(
		Class("bulma-container"),
		c.Element("breadcrumbs", Breadcrumbs),
		router.Match(
			c,
			Route("/projects/new", c.Element("newProject", NewProject())),
			Route("/projects/(?P<projectId>[^/]+)(?:/(?P<tab>actions|changes|settings))?", ProjectDetails),
			Route("/projects", c.Element("projects", Projects)),
		),
	)
}
