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

func List(args ...any) *HTMLElement {
	return Div(
		Class("kip-list"),
		args,
	)
}

func ListItem(args ...any) *HTMLElement {
	return Div(
		Class("kip-item", "kip-is-card"),
		args,
	)
}

func ListHeader(args ...any) *HTMLElement {
	return Div(
		Class("kip-item", "kip-is-header"),
		args,
	)
}

func ListColumn(size string, args ...any) *HTMLElement {
	return Div(
		Class("kip-col", Fmt("kip-is-%s", size)),
		args,
	)
}
