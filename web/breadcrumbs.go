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

type Breadcrumb struct {
	Title string
	Path  string
}

func Breadcrumbs(c Context) Element {

	crumbs := []Element{}

	breadcrumbs := UseGlobal[[]Breadcrumb](c, "breadcrumbs")

	path := ""
	title := ""

	for _, breadcrumb := range breadcrumbs {

		path += breadcrumb.Path

		if title != "" {
			title += " :: "
		}

		title += breadcrumb.Title

		crumbs = append(crumbs, Li(
			A(Href(path), breadcrumb.Title),
		))
	}

	return Nav(
		Class("bulma-breadcrumb bulma-has-bullet-separator"),
		Ul(
			crumbs,
		),
	)

}

func AddBreadcrumb(c Context, title string, path string) {

	breadcrumbs := GlobalVar(c, "breadcrumbs", []Breadcrumb{})

	bcs := breadcrumbs.Get()

	bcs = append(bcs, Breadcrumb{
		Title: title,
		Path:  path,
	})

	breadcrumbs.Set(bcs)
}

func MainTitle(c Context) string {

	breadcrumbs := GlobalVar(c, "breadcrumbs", []Breadcrumb{})

	title := ""

	for _, breadcrumb := range breadcrumbs.Get() {

		if title != "" {
			title += " :: "
		}

		title += breadcrumb.Title
	}

	// we reset the breadcrumbs
	breadcrumbs.Set([]Breadcrumb{})

	return title

}
