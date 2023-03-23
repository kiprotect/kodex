package web

import (
	. "github.com/kiprotect/gospel"
)

func Kodex(c Context) Element {
	return F(
		H1(Class("title"), "Hi, world!"),
		P("foo"),
	)
}
