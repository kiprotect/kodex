package ui

import (
	. "github.com/gospel-sh/gospel"
)

func Modal(c Context, closeUrl string) Element {
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
					"test",
				),
			),
			Section(
				Class("bulma-modal-card-body"),
				Div("test"),
			),
			Footer(
				Class("bulma-modal-card-foot"),
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
