package web

import (
	. "github.com/gospel-dev/gospel"
	"github.com/kiprotect/kodex"
)

func Dropdown(c Context, id, icon string, items []Element) Element {
	return Div(
		Id(id),
		Attrib("data-type")("navbar-dropdown"),
		Class("kip-navbar-dropdown-menu", "bulma-navbar-item", "bulma-has-dropdown"),
		A(
			Attrib("aria-has-popup")("true"),
			Attrib("aria-expanded")("true"),
			Class("bulma-navbar-link"),
			OnClick(Fmt("toggleDropdown('%s', event)", id)),
			Div(
				Class("kip-nowrap"),
				Span(
					Class("icon", "is-small"),
					I(
						Class("fas", icon),
					),
				),
			),
		),
		Div(
			Class("kip-navbar-dropdown", "bulma-navbar-dropdown", "bulma-is-right"),
			items,
		),
	)
}

func AppNavbar(c Context) Element {

	plugins := UsePlugins(c)

	items := []Element{}

	for _, plugin := range plugins {
		if appLinkPlugin, ok := plugin.(AppLinkPlugin); ok {

			name, icon, path := appLinkPlugin.AppLink()

			items = append(items, A(
				Class("bulma-dropdown-item"),
				Href(path),
				Span(
					Class("icon", "is-small"),
					I(
						Class("fas", icon),
					),
				),
				name,
			))
		}
	}

	if len(items) == 0 {
		return nil
	}

	return Dropdown(c, "apps-dropdown", "fa-th-large", items)
}

func UserNavbar(c Context) Element {

	// get the logged in user
	user := UseExternalUser(c)

	items := []Element{
		Div(
			Class("bulma-dropdown-item"),
			Span(
				Class("kip-overflow-ellipsis"),
				user.Email,
			),
		),
		Hr(
			Class("bulma-dropdown-divider"),
		),
		A(
			Class("bulma-dropdown-item"),
			Href("/logout"),
			Span(
				Class("icon", "is-small"),
				I(
					Class("fas", "fa-sign-out-alt"),
				),
			),
			"Logout",
		),
	}

	return Dropdown(c, "user-dropdown", "fa-user-circle", items)

}

func Navbar(c Context) Element {

	return Header(
		Class("kip-navbar", "bulma-navbar", "bulma-is-fixed-top"),
		Div(
			Class("centered"),
			Div(
				Class("bulma-navbar-brand"),
				Div(
					Class("kip-logo-wrapper"),
					A(
						Href("/#doNotRedirect"),
						Img(
							Class("kip-logo", Alt("projects")),
							Src("/static/images/kodexlogo-blue.png"),
						),
						Img(
							Class("kip-small-logo", Alt("projects")),
							Src("/static/images/kiprotect-k.png"),
						),
						Span(
							Class("kip-version"),
							kodex.Version,
						),
					),
				),
			),
		),
		Div(
			Class("bulma-navbar-menu"),
			Div(
				Class("bulma-navbar-end"),
				AppNavbar(c),
				UserNavbar(c),
			),
		),
		Script(`

			function closeMenu(e){
				let dropdowns = document.querySelectorAll('[data-type="navbar-dropdown"]');
				for(const dropdown of dropdowns){
					if (!dropdown.contains(e.target)){
						dropdown.classList.remove('bulma-is-active');
					}		
				}
			}

			window.addEventListener("click", closeMenu, false);

			function toggleDropdown(id, e){
				window[id].classList.toggle('bulma-is-active');
			}`,
		),
	)
}
