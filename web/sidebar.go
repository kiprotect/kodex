package web

import (
	. "github.com/gospel-sh/gospel"
)

type SidebarItem struct {
	Title string
	Path  string
	Icon  string
}

func InitSidebar(c Context) {
	// we get the sidebar variable
	itemsVar := GlobalVar(c, "sidebar", []SidebarItem{})

	// we reset the sidebar items
	itemsVar.Set([]SidebarItem{})
}

func AddSidebarItem(c Context, item SidebarItem) {

	itemsVar := GlobalVar(c, "sidebar", []SidebarItem{})
	itemsVar.Set(append(itemsVar.Get(), item))
}

func WithSidebar(sidebar any, content any) ElementFunction {

	return func(c Context) Element {

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

}

func NavItem(c Context, item SidebarItem) Element {
	return Li(
		Class("kip-nav-item"),
		A(
			Href(item.Path),
			Span(
				Span(
					Class("icon", "is-small"),
					I(
						Class("fas", Fmt("fa-%s", item.Icon)),
					),
				),
				item.Title,
			),
		),
	)
}

func MenuItems(c Context) Element {

	items := UseGlobal[[]SidebarItem](c, "sidebar")
	menuItems := make([]Element, 0, len(items))

	for _, item := range items {

		menuItems = append(menuItems, NavItem(c, item))
	}
	return F(
		menuItems,
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
		Class("kip-sidebar"),
		c.Element("menu", KipMenu),
	)
}
