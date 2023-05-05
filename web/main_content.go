package web

import (
	. "github.com/kiprotect/gospel"
)

func NotFoundRedirect(c Context) Element {
	router := UseRouter(c)
	router.RedirectTo("/404")
	return nil
}

func MainContent(c Context) Element {

	// get the router
	router := UseRouter(c)

	return Div(
		Class("bulma-container"),
		c.Element("breadcrumbs", Breadcrumbs),
		router.Match(
			c,
			Route("/projects/new", c.ElementFunction("newProject", NewProject())),
			Route("/projects/(?P<projectId>[^/]+)(?:/(?P<tab>actions|changes|settings))?", ProjectDetails),
			Route("/projects|^/$", c.ElementFunction("projects", Projects)),
			Route("", c.ElementFunction("notFound", NotFoundRedirect)),
		),
	)
}
