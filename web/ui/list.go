package ui

import (
	. "github.com/gospel-dev/gospel"
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
