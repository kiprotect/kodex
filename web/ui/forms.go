package ui

import (
	. "github.com/gospel-dev/gospel"
)

func Message(msgType string, content any) Element {
	return Div(
		Class("bulma-message", Fmt("bulma-is-%s", msgType)),
		Div(
			Class("bulma-message-body"),
			content,
		),
	)
}

func MessageWithTitle(msgType string, title, content any) Element {
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
