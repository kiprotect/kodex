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

package web

import (
	. "github.com/gospel-sh/gospel"
	"strings"
)

type SidebarItem struct {
	Title   string
	Path    string
	Icon    string
	Header  bool
	Submenu []*SidebarItem
}

func InitSidebar(c Context) {
	// we get the sidebar variable
	itemsVar := GlobalVar(c, "sidebar", []*SidebarItem{})

	// we reset the sidebar items
	itemsVar.Set([]*SidebarItem{})
}

func getItem(items []*SidebarItem, path string) *SidebarItem {
	for _, item := range items {
		if item.Path == path {
			return item
		}
		if item.Submenu != nil {
			if item := getItem(item.Submenu, path); item != nil {
				return item
			}
		}
	}
	return nil
}

func GetSidebarItemByPath(c Context, path string) *SidebarItem {
	return getItem(UseGlobal[[]*SidebarItem](c, "sidebar"), path)
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

func Submenu(c Context, level int, items []*SidebarItem) Element {

	return Ul(
		Class("kip-menu-list", Fmt("kip-is-submenu-%d", level)),
		menuItems(c, level+1, items),
	)
}

func NavItem(c Context, item *SidebarItem) Element {

	router := UseRouter(c)
	active := strings.HasPrefix(router.FullPath(), item.Path)

	return Li(
		Class(
			"kip-nav-item",
			If(item.Header, "kip-is-header"),
			If(active, "kip-is-active"),
		),
		A(
			Href(item.Path),
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
	return menuItems(c, 1, UseGlobal[[]*SidebarItem](c, "sidebar"))
}

func menuItems(c Context, level int, items []*SidebarItem) Element {

	router := UseRouter(c)
	menuItems := make([]Element, 0, len(items))

	for _, item := range items {
		active := strings.HasPrefix(router.FullPath(), item.Path)
		menuItems = append(menuItems, NavItem(c, item))
		if len(item.Submenu) > 0 && active {
			menuItems = append(menuItems, Submenu(c, level, item.Submenu))
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
