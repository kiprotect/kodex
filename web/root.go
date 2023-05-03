package web

import (
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex/api"
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
		return Div("errr")
	}

	SetExternalUser(c, externalUser)
	SetApiUser(c, apiUser)

	return F(
		c.Element("navHeader", Navbar),
		c.Element("contentWithSidebar", WithSidebar(
			c.Element("sidebar", Sidebar),
			c.Element("mainContent", MainContent),
		)),
	)

}

func Logout(c Context) Element {
	return Div()
}

func Login(c Context) Element {
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
						H2("Login"),
					),
				),
				Div(
					Class("kip-card-content", "kip-card-centered"),
					Div(
						Class("kip-login"),
						Div(
							Class("kip-provider-list"),
							Ul(
								Li(
									A(
										Href("/api/v1/login"),
										Button(
											Class("bulma-button", "bulma-is-success", "bulma-is-flex"),
											"Log in via SSO",
										),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}

func AppContent(c Context) Element {

	router := UseRouter(c)

	return router.Match(
		c,
		Route("/login", Login),
		Route("/logout", Logout),
		Route("", AuthorizedContent),
	)

}

func Root(controller api.Controller) (func(c Context) Element, error) {

	userProvider, err := controller.UserProvider()

	if err != nil {
		return nil, err
	}

	return func(c Context) Element {

		// we set the user provider
		SetUserProvider(c, userProvider)

		// we set the controller
		SetController(c, controller)

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
