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
	"github.com/kiprotect/kodex/api"
	"net/http"
	"time"
)

func AuthorizedContent(c Context) Element {

	userProvider := UseUserProvider(c)
	controller := UseController(c)
	router := UseRouter(c)

	// we get the user from the provider...
	externalUser, _ := userProvider.Get(controller, c.Request())

	// we redirect to the login page
	if externalUser == nil {
		router.RedirectTo("/login")
		return nil
	}

	apiUser, err := externalUser.ApiUser(controller)

	if err != nil {
		return Div(Fmt("Cannot get user: %v", err))
	}

	SetExternalUser(c, externalUser)
	SetApiUser(c, apiUser)

	return router.Match(
		c,
		Route("^/?$", c.ElementFunction("chooser", AppChooser)),
		Route("", c.ElementFunction("main", func(c Context) Element {
			return F(
				c.Element("navHeader", Navbar),
				c.Element("mainContent", MainContent),
			)
		})),
	)
}

func Logout(c Context) Element {

	userProvider := UseUserProvider(c)
	controller := UseController(c)

	userProvider.Logout(controller, c.Request(), c.ResponseWriter())

	// we clear the context
	c.Clear()

	return Section(
		Class("kip-centered-card", "kip-is-info", "kip-is-fullheight"),
		Div(
			Class("kip-card", "kip-is-centered", "kip-account"),
			Div(
				Class("kip-card-header"),
				Div(
					Class("kip-card-title"),
					H2("Logout"),
				),
			),
			Div(
				Class("kip-card-content", "kip-card-centered"),
				P(
					"You have been logged out. ",
					A(Href("/login"), "Log back in."),
				),
			),
		),
	)
}

func TokenLogin(c Context) Element {

	token := Var(c, "")
	error := Var(c, "")
	router := UseRouter(c)
	onSubmit := Func[any](c, func() {

		if token.Get() == "" {
			error.Set("Please enter a token value")
			return
		}

		w := c.ResponseWriter()

		http.SetCookie(w, &http.Cookie{Path: "/", Name: "kodex-auth", Value: token.Get(), Secure: false, HttpOnly: true, Expires: time.Now().Add(365 * 24 * 7 * time.Hour)})

		router.RedirectTo("/")

	})

	var errorNotice Element

	if error.Get() != "" {
		errorNotice = P(
			Class("bulma-help", "bulma-is-danger"),
			error.Get(),
		)
	}

	return Section(
		Class("kip-centered-card", "kip-is-info", "kip-is-fullheight"),
		Div(
			Class("kip-card", "kip-is-centered", "kip-account"),
			Div(
				Class("kip-card-header"),
				Div(
					Class("kip-card-title"),
					H2("Access Token Login"),
				),
			),
			Div(
				Class("kip-card-content", "kip-card-centered", "kip-provider-list"),
				Form(
					Method("POST"),
					OnSubmit(onSubmit),
					Div(
						Class("bulma-field"),
						errorNotice,
						Label(
							Class("bulma-label", "Token"),
							Input(
								Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
								Value(token),
								Placeholder("Token value"),
							),
						),
					),
					Div(
						Class("bulma-field"),
						P(
							Class("bulma-control"),
							Button(
								Class("bulma-button", "bulma-is-success"),
								Type("submit"),
								"Log in via token",
							),
						),
					),
				),
			),
		),
	)
}

func NotFound(c Context) Element {

	c.SetStatusCode(404)

	return Div(
		Class("kip-with-app-selector"),
		A(
			Class("kip-with-app-selector-link"),
			Href("/#"),
			Div(
				Class("kip-logo-wrapper"),
				Img(
					Class("kip-logo", Alt("projects")),
					Src("/static/images/kodexlogo-white.png"),
				),
			),
		),
		Section(
			Class("kip-centered-card", "kip-is-info", "kip-is-fullheight"),
			Div(
				Class("kip-card", "kip-is-centered", "kip-account"),
				Div(
					Class("kip-card-header"),
					Div(
						Class("kip-card-title"),
						H2("404 - Page Not Found"),
					),
				),
				Div(
					Class("kip-card-content", "kip-card-centered"),
					Div(
						"Sorry, there's nothing here...",
					),
				),
			),
		),
	)
}

func AppContent(c Context) Element {

	router := UseRouter(c)
	plugins := UsePlugins(c)

	routes := []*RouteConfig{
		Route("/token-login", TokenLogin),
		Route("/logout", Logout),
		Route("/404", NotFound),
	}

	// we add the main plugin routes
	for _, plugin := range plugins {
		routes = append(routes, plugin.Routes(c).Main...)
	}

	// we serve the authorized content
	routes = append(routes, Route("", AuthorizedContent))

	return router.Match(
		c,
		routes...,
	)

}

func Root(controller api.Controller, plugins []WebPlugin) (func(c Context) Element, error) {

	return func(c Context) Element {

		// we create a fresh clone of the controller so that e.g. transactions
		// remain isolated to this goroutine...
		controller, err := controller.ApiClone()

		if err != nil {
			return Div("cannot create controller")
		}

		userProvider, err := controller.UserProvider()

		if err != nil {
			// should never happen
			return Div("cannot create user provider")
		}

		// we set the user provider
		SetUserProvider(c, userProvider)

		// we set the controller
		SetController(c, controller)

		SetPlugins(c, plugins)

		return F(
			Doctype("html"),
			Html(
				Lang("en"),
				Head(
					Meta(Charset("utf-8")),
					Title(MainTitle(c)),
					Link(Rel("preload"), Href("/static/main.css"), As("style")),
					// Link(Rel("apple-touch-icon"), Sizes("180x180"), Href("/icons/apple-touch-icon.png")),
					// Link(Rel("icon"), Type("image/png"), Sizes("32x32"), Href("/icons/favicon-32x32.png")),
					Link(Rel("stylesheet"), Href("/static/main.css")),
					Script(Defer(), Src("/static/gospel.js"), Type("module")),
				),
				Body(
					Class("kip-fonts", "bulma-body"),
					Div(
						Class("kip"),
						c.Element("appContent", AppContent),
					),
				),
			),
		)

	}, nil
}
