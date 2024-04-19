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
