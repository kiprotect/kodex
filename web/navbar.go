package web

import (
	. "github.com/kiprotect/gospel"
)

func Navbar(c Context) Element {

	return Header(
		Class("kip-navbar", "bulma-navbar", "bulma-is-fixed-top"),
		Div(
			Class("centered"),
			Div(
				Class("bulma-navbar-brand"),
				Div(
					Class("kip-logo-wrapper"),
					A(
						Href("/#doNotRedirect"),
						Img(
							Class("kip-logo", Alt("projects")),
							Src("/static/images/kodexlogo-blue.png"),
						),
						Img(
							Class("kip-small-logo", Alt("projects")),
							Src("/static/images/kiprotect-k.png"),
						),
						Span(
							Class("kip-version"),
							"latest",
						),
					),
				),
			),
		),
		H1(Class("bulma-navbar-item", "bulma-navbar-title"), "test"),
	)
}
