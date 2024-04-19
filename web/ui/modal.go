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

package ui

import (
	. "github.com/gospel-sh/gospel"
)

func Modal(c Context, title string, content, footer any, closeUrl string) Element {
	return Div(
		Class("bulma-modal", "bulma-is-active"),
		Div(
			Class("bulma-modal-background"),
			OnClick(Fmt("location.href = '%s'", closeUrl)),
		),
		Div(
			Class("bulma-modal-card"),
			Header(
				Class("bulma-modal-card-head"),
				P(
					Class("bulma-modal-card-title"),
					title,
				),
			),
			Section(
				Class("bulma-modal-card-body"),
				content,
			),
			Footer(
				Class("bulma-modal-card-foot"),
				footer,
			),
		),
	)
}

// A Modal that embeds a form
func FormModal(c Context, headers, content, footer any, title, closeUrl string) Element {
	return Div(
		Class("bulma-modal", "bulma-is-active"),
		Div(
			Class("bulma-modal-background"),
			OnClick(Fmt("location.href = '%s'", closeUrl)),
		),
		Form(
			headers,
			Div(
				Class("bulma-modal-card"),
				Header(
					Class("bulma-modal-card-head"),
					P(
						Class("bulma-modal-card-title"),
						title,
					),
				),
				Section(
					Class("bulma-modal-card-body"),
					content,
				),
				Footer(
					Class("bulma-modal-card-foot"),
					footer,
				),
			),
		),
	)
}
