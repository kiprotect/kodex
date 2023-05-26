package web

import (
	. "github.com/gospel-dev/gospel"
)

func NotFoundRedirect(c Context) Element {
	router := UseRouter(c)
	router.RedirectTo("/404")
	return nil
}

func MainContent(c Context) Element {

	// get the router
	router := UseRouter(c)

	plugins := UsePlugins(c)

	routes := []*RouteConfig{
		Route("/projects/new", c.ElementFunction("newProject", NewProject())),
		Route("/projects/(?P<projectId>[^/]+)(?:/(?P<tab>actions|streams|changes|settings))?", ProjectDetails),
		Route("/projects|^/$", c.ElementFunction("projects", Projects)),
	}

	// we add the main plugin routes
	for _, plugin := range plugins {
		routes = append(routes, plugin.MainRoutes(c)...)
	}

	// we add a "not found" catch-all route...
	routes = append(routes, Route("", c.ElementFunction("notFound", NotFoundRedirect)))

	return Div(
		Class("bulma-container"),
		c.Element("breadcrumbs", Breadcrumbs),
		router.Match(
			c,
			routes...,
		),
	)
}
