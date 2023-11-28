package web

import (
	. "github.com/gospel-sh/gospel"
)

func AppItem(path, icon, name string) Element {
	return Li(A(
		Href(path),
		Span(
			Class("icon", "is-small"),
			I(
				Class("fas", icon),
			),
		),
		name,
	))
}

func AppChooser(c Context) Element {

	plugins := UsePlugins(c)
	router := UseRouter(c)
	user := UseExternalUser(c)

	items := []Element{}

	superuser := user.HasRole(nil, "superuser")

	if superuser {
		items = append(items, AppItem("/admin", "admin", "Administration"))
	}

	for _, plugin := range plugins {
		if appLinkPlugin, ok := plugin.(AppLinkPlugin); ok {
			appLink := appLinkPlugin.AppLink()
			if !appLink.Superuser || superuser {
				items = append(items, AppItem(appLink.Path, appLink.Icon, appLink.Name))
			}
		}
	}

	if len(items) == 0 {
		// there's only one choice, so we redirect directly
		router.RedirectTo("/flows")
		return nil
	}

	// we prepend the flows app
	items = append([]Element{AppItem("/flows", "flows", "Data Flows")}, items...)

	return Section(
		Class("kip-centered-card", "kip-is-info", "kip-is-fullheight"),
		Div(
			Class("kip-card", "kip-is-centered", "kip-account"),
			Div(
				Class("kip-card-header"),
				Div(
					Class("kip-card-title"),
					H2("Kodex - App Selector"),
				),
			),
			Div(
				Class("kip-card-content", "kip-card-centered", "kip-provider-list"),
				Aside(
					Class("bulma-menu"),
					Ul(
						Class("bulma-menu-list"),
						items,
					),
				),
			),
		),
	)

}
