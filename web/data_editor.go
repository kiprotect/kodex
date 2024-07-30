package web

import (
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/go-helpers/forms"
	"sort"
)

/*
- Each map shows the key, the value (if it is a scalar on the right,
  if it is an array or a map below), the associated validators (if any) and
  on the right side the transformed value.

*/

func SliceEditor(c Context, value []any, path []string) Element {

	items := []Element{}

	for i, nv := range value {

		extraContent := AnyEditor(c, nv, append(path, Fmt("%d", i)))

		var item Element

		if extraContent == nil {
			item = Li(
				Span(Class("kip-key"), Fmt("%d", i)),
			)
		} else {
			item = Li(
				Div(
					Class("kip-extra-content"),
					extraContent,
				),
			)
		}

		items = append(items, item)
	}

	return Ul(
		Class("kip-slice-diff"),
		items,
	)
}

func AnyEditor(c Context, value any, path []string) Element {

	switch nv := value.(type) {
	case map[string]any:
		return DataEditor(c, nv, nil, path)
	case []any:
		return SliceEditor(c, nv, path)
	}
	// we don't return anything
	return nil
}

func MapValueEditor(c Context, key string, value any, path []string) Element {

	extraContent := AnyEditor(c, value, path)

	var fromTo Element

	return Li(
		Span(Class("kip-key"), key),
		fromTo,
		If(
			extraContent != nil,
			Div(
				Class("kip-extra-content"),
				extraContent,
			),
		),
	)
}

func DataEditor(c Context, data map[string]any, form *forms.Form, path []string) Element {
	values := []Element{}
	keys := []string{}

	for key, _ := range data {
		keys = append(keys, key)
	}

	// we always sort keys
	sort.Strings(keys)

	for _, key := range keys {
		value := data[key]
		values = append(values, MapValueEditor(c, key, value, append(path, key)))
	}

	return Ul(
		Class("kip-map-diff", If(len(path) == 0, "kip-top-level")),
		values,
	)

}
