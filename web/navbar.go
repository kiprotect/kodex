package web

import (
	. "github.com/gospel-dev/gospel"
)

func AppNavbar(c Context) Element {
	return Div(
		Class("bulma-navbar-end"),
		Div(
			Class("bulma-navbar-dropdown-menu", "bulma-navbar-item", "bulma-has-dropdown"),
		),
		Div(
			Id("app-dropdown"),
			Class("bulma-navbar-dropdown-menu", "bulma-navbar-item", "bulma-has-dropdown"),
			A(
				Attrib("aria-has-popup")("true"),
				Attrib("aria-expanded")("true"),
				Class("bulma-navbar-link"),
				OnClick("toggleUserMenu(event)"),
				Div(
					Class("kip-nowrap"),
					Span(
						Class("icon", "is-small"),
						I(
							Class("fas", "fa-user-circle"),
						),
					),
				),
			),
			Div(
				Class("kip-navbar-dropdown", "bulma-navbar-dropdown", "bulma-is-right"),
				Div(
					Class("bulma-dropdown-item"),
					Span(
						Class("kip-overflow-ellipsis"),
						"test",
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
			),
		),
	)
}

func UserNavbar(c Context) Element {

	// get the logged in user
	user := UseExternalUser(c)

	return Div(
		Class("bulma-navbar-end"),
		Div(
			Class("bulma-navbar-dropdown-menu", "bulma-navbar-item", "bulma-has-dropdown"),
		),
		Div(
			Id("user-dropdown"),
			Class("bulma-navbar-dropdown-menu", "bulma-navbar-item", "bulma-has-dropdown"),
			A(
				Attrib("aria-has-popup")("true"),
				Attrib("aria-expanded")("true"),
				Class("bulma-navbar-link"),
				OnClick("toggleUserMenu(event)"),
				Div(
					Class("kip-nowrap"),
					Span(
						Class("icon", "is-small"),
						I(
							Class("fas", "fa-user-circle"),
						),
					),
				),
			),
			Div(
				Class("kip-navbar-dropdown", "bulma-navbar-dropdown", "bulma-is-right"),
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
			),
		),
	)
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
							"latest",
						),
					),
				),
			),
		),
		Div(
			Class("bulma-navbar-meanu"),
			AppNavbar(c),
			UserNavbar(c),
		),
		Script(`

let dropdown = document.getElementById('user-dropdown');

function closeMenu(e){
	if (!dropdown.contains(e.target)){
		dropdown.classList.remove('bulma-is-active');
	}
}

window.addEventListener("click", closeMenu, false);

function toggleUserMenu(e){
	dropdown.classList.toggle('bulma-is-active');
}
		`),
	)
}
