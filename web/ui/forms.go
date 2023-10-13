package ui

import (
	. "github.com/gospel-sh/gospel"
)

func Message(msgType string, content any) *HTMLElement {
	return Div(
		Class("bulma-message", Fmt("bulma-is-%s", msgType)),
		Div(
			Class("bulma-message-body"),
			content,
		),
	)
}

func MessageWithTitle(msgType string, title, content any) *HTMLElement {
	return Div(
		Class("bulma-message", Fmt("bulma-is-%s", msgType)),
		Div(
			Class("bulma-message-header"),
			P(title),
		),
		Div(
			Class("bulma-message-body"),
			content,
		),
	)
}
