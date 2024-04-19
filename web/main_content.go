// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
		Route("/user", c.ElementFunction("user", UserProfile)),
	}

	// we add the main plugin routes
	for _, plugin := range plugins {
		routes = append(routes, plugin.Routes(c).Authorized...)
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
