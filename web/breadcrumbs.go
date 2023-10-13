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
