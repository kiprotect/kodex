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

func Tabs(args ...any) Element {
	return Div(
		Class("bulma-tabs"),
		Class("active"),
		Span(
			Class("bulma-more"),
			Span(L("&or;")),
		),
		Ul(
			args,
		),
	)
}

func ActiveTab(active bool) Attribute {
	if active {
		return Class("bulma-is-active")
	}
	return nil
}

func Tab(args ...any) Element {
	return Li(
		args,
	)
}
