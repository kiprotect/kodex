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
	items := []Element{
		AppItem("/flows", "flows", "Data Flows"),
		AppItem("/admin", "admin", "Administration"),
	}

	for _, plugin := range plugins {
		if appLinkPlugin, ok := plugin.(AppLinkPlugin); ok {
			name, icon, path := appLinkPlugin.AppLink()
			items = append(items, AppItem(path, icon, name))
		}
	}

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
