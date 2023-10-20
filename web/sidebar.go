package web

import (
	. "github.com/gospel-sh/gospel"
)

type SidebarItem struct {
	Title   string
	Path    string
	Icon    string
	Header  bool
	Active  bool
	Submenu []*SidebarItem
}

func InitSidebar(c Context) {
	// we get the sidebar variable
	itemsVar := GlobalVar(c, "sidebar", []*SidebarItem{})

	// we reset the sidebar items
	itemsVar.Set([]*SidebarItem{})
}

func GetSidebarItemByPath(c Context, path string) *SidebarItem {

	for _, item := range UseGlobal[[]*SidebarItem](c, "sidebar") {
		if item.Path == path {
			return item
		}
	}

	return nil
}

func AddSidebarItem(c Context, item *SidebarItem) {

	itemsVar := GlobalVar(c, "sidebar", []*SidebarItem{})
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

func Submenu(c Context, items []*SidebarItem) Element {

	return Ul(
		Class("kip-menu-list", "kip-is-submenu"),
		menuItems(c, items),
	)
}

func NavItem(c Context, item *SidebarItem) Element {
	return Li(
		Class("kip-nav-item", If(item.Header, "kip-is-header")),
		A(
			Href(item.Path),
			If(item.Active, Class("kip-is-active")),
			Span(
				If(
					item.Icon != "",
					Span(
						Class("icon", "is-small"),
						I(
							Class("fas", Fmt("fa-%s", item.Icon)),
						),
					),
				),
				item.Title,
			),
		),
	)
}

func MenuItems(c Context) Element {
	return menuItems(c, UseGlobal[[]*SidebarItem](c, "sidebar"))
}

func menuItems(c Context, items []*SidebarItem) Element {

	menuItems := make([]Element, 0, len(items))

	for _, item := range items {
		menuItems = append(menuItems, NavItem(c, item))
		if len(item.Submenu) > 0 {
			menuItems = append(menuItems, Submenu(c, item.Submenu))
		}
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
		Id("sidebar"),
		Class("kip-sidebar"),
		c.Element("menu", KipMenu),
	)
}
