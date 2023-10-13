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
			Span("cm-tabs-more", "&or;"),
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
