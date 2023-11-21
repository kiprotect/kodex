package web

import (
	. "github.com/gospel-sh/gospel"
)

func NotFoundRedirect(c Context) Element {
	router := UseRouter(c)
	router.RedirectTo("/404")
	return nil
}

func Flows(c Context) Element {

	router := UseRouter(c)

	AddSidebarItem(c, &SidebarItem{Title: "Projects", Path: "/flows/projects", Icon: "bars"})

	return router.Match(
		c,
		Route("/projects/new$", c.ElementFunction("newProject", NewProject())),
		Route("/projects/(?P<projectId>[^/]+)(?:/(?P<tab>actions|streams|changes|settings))?", ProjectDetails),
		Route("(?:/(?:projects)?)?$", c.ElementFunction("projects", Projects)),
	)
}

func MainContent(c Context) Element {

	// get the router
	router := UseRouter(c)
	plugins := UsePlugins(c)

	// we initialize the sidebar
	InitSidebar(c)

	routes := []*RouteConfig{
		Route("/flows", c.ElementFunction("flows", Flows)),
		Route("/admin", c.ElementFunction("admin", Admin)),
	}

	// we add the main plugin routes
	for _, plugin := range plugins {
		routes = append(routes, plugin.MainRoutes(c)...)
	}

	// we add a "not found" catch-all route...
	routes = append(routes, Route("", c.ElementFunction("notFound", NotFoundRedirect)))

	return WithSidebar(
		c.DeferElement("sidebar", Sidebar),
		F(
			Div(
				Class("bulma-container"),
				router.Match(
					c,
					routes...,
				),
			),
		),
	)(c)
}
