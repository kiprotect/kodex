package web

import (
	. "github.com/kiprotect/gospel"
)

func WithSidebar(c Context) Element {

	sidebar := c.Element("sidebar", Sidebar)
	content := c.Element("mainContent", MainContent)

	return Div(
		Class("kip-with-sidebar"),
		Div(
			Class("kip-with-sidebar__sidebar"),
			sidebar,
		),
		Div(
			Class("kip-with-sidebar__content"),
			content,
		),
	)
}

func NavItem(c Context) Element {
	return Li(
		Class("kip-nav-item"),
		A(
			Href("/projects"),
			Span(
				Span(
					Class("icon", "is-small"),
					I(
						Class("fas", "fa-chalkboard"),
					),
				),
				"Projects",
			),
		),
	)
}

func MenuItems(c Context) Element {
	return F(
		c.Element("navItem", NavItem),
	)
}

func KipMenu(c Context) Element {
	return Aside(
		Class("kip-menu-aside"),
		Ul(
			Class("kip-menu-list"),
			c.Element("menuItems", MenuItems),
		),
	)
}

func Sidebar(c Context) Element {
	return Div(
		Class("kip-sidebar", "kip-sidebar--collapsed"),
		c.Element("menu", KipMenu),
	)
}
