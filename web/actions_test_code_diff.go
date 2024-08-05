package web

import (
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"sort"
	"strings"
)

func FromToCode(newValue, oldValue any, path []string) Element {
	var class = "unkown"
	if newValue == oldValue {
		class = "unmodified"
	} else if oldValue == nil {
		class = "added"
	} else if newValue == nil {
		class = "removed"
	} else {
		class = "modified"
	}
	return Span(
		Class("kip-from-to", Fmt("kip-%s", class)),
		IfElse(
			newValue == oldValue,
			Span(Class("kip-from"), IfElse(oldValue == nil, L("&#8709;"), F(Fmt("%v", oldValue)))),
			F(
				Span(Class("kip-from"), IfElse(oldValue == nil, L("&#8709;"), F(Fmt("%v", oldValue)))),
				" 🡒 ",
				Span(Class("kip-to"), IfElse(newValue == nil, L("&#8709;"), F(Fmt("%v", newValue)))),
			),
		),
	)
}

func KeyPropsCode(c Context, key string, path []string) Element {

	router := UseRouter(c)
	qp := queryPath(c)
	routePath := path
	selected := fieldIsSelected(path, qp)

	return F(
		strings.Repeat(" ", (len(path)+1)*2),
		Span(
			Class("kip-key-props"),
			Span(Class("kip-key"),
				A(
					Href(PathWithQuery(router.CurrentPath(), map[string][]string{
						"field":    routePath,
						"action":   []string{""},
						"dataItem": router.Query()["dataItem"],
					})),
					key,
					If(selected, Id("selectedKey")),
					If(selected, Class("kip-key-selected")),
				),
			),
			If(selected, F(
				Script(`setTimeout(function(){selectedKey.scrollIntoView()}, 100)`),
			)),
		),
		": ",
	)
}

func CodeSliceDiff(c Context, newValue, oldValue []any, validators []forms.Validator, path []string) Element {

	qp := queryPath(c)
	items := []Element{}
	var listValidators []forms.Validator

	for _, validator := range validators {
		if listValidator, ok := validator.(*forms.IsList); ok {
			listValidators = listValidator.Validators
			break
		}
	}

	for i, ov := range oldValue {
		var nv any

		if i < len(newValue) {
			nv = newValue[i]
		}

		key := Fmt("%d", i)
		newPath := append(path, key)
		fullMatch := dataPathMatches(newPath, qp)

		if !fullMatch {
			// uncomment to show full values
			/* if len(path) < len(qp) {
				continue
			}*/
			items = append(items, F(
				KeyPropsCode(c, key, newPath),
				PlaceholderCode(c, nv, ov),
			))
			continue
		}

		extraContent := AnyDiffCode(c, nv, ov, listValidators, newPath)

		var item Element

		if extraContent == nil {
			extraContent = FromTo(nv, ov)
		}

		item = F(
			KeyPropsCode(c, key, newPath),
			//			Span(
			//				Class("kip-validator-list"),
			//				Span(Validators(c, listValidators, path, nil, MakeExtra("class", "kip-validators-code"))),
			//			),
			extraContent,
		)

		items = append(items, item)
	}

	prefix := strings.Repeat(" ", (len(path)+1)*2)

	return Span(
		Class("kip-slice-diff"),
		"[\n",
		items,
		prefix,
		If(len(items) == 0, F("  (empty list)\n", prefix)),
		"]\n",
	)
}

func AnyDiffCode(c Context, newValue, oldValue any, validators []forms.Validator, path []string) Element {

	if isStringMap(oldValue) || isStringMap(newValue) {

		var form *forms.Form

		for _, validator := range validators {
			if mapValidator, ok := validator.(*forms.IsStringMap); ok {
				form = mapValidator.Form
			}
		}

		return CodeMapDiff(c, toStringMap(newValue), toStringMap(oldValue), form, path)
	}

	if isSlice(oldValue) || isSlice(newValue) {
		return CodeSliceDiff(c, toSlice(newValue), toSlice(oldValue), validators, path)
	}

	// we don't return anything
	return nil
}

func MapCodeValue(c Context, key string, newValue, oldValue any, validators []forms.Validator, path []string) Element {

	qp := queryPath(c)
	fullMatch := dataPathMatches(path, qp)
	diffable := diffableValue(newValue, oldValue)

	if !diffable {

		// uncomment to show full values
		/* if !fullMatch && len(path) <= len(qp) {
			return nil
		}*/

		return Span(
			Class("kip-map-field"),
			Span(
				KeyPropsCode(c, key, path),
				Class("kip-field-and-validators"),
				FromToCode(newValue, oldValue, path),
				//				Span(
				//					Class("kip-validator-list"),
				//					Span(Validators(c, validators, path, nil, MakeExtra("class", "kip-validators-code"))),
				//				),
			),
			"\n",
		)
	}

	// we check if this field is selected
	if !fullMatch {
		// uncomment to show all values
		/* if len(path) <= len(qp) {
			return nil
		} */
		return F(
			Span(
				Class("kip-field-and-validators"),
				KeyPropsCode(c, key, path),
				//				Span(
				//					Class("kip-validator-list"),
				//					Span(Validators(c, validators, path, nil, MakeExtra("class", "kip-validators-code"))),
				//				),
			),
			PlaceholderCode(c, newValue, oldValue),
		)
	}

	extraContent := AnyDiffCode(c, newValue, oldValue, validators, path)

	return Span(
		Class("kip-map-field"),
		Span(
			Class("kip-field-and-validators"),
			KeyPropsCode(c, key, path),
			//			Span(
			//				Class("kip-validator-list"),
			//				Span(Validators(c, validators, path, nil, MakeExtra("class", "kip-validators-code"))),
			//			),
		),
		extraContent,
	)
}

func PlaceholderCode(c Context, newValue, oldValue any) Element {
	if isStringMap(oldValue) || isStringMap(newValue) {
		return F(
			"{…}\n",
		)
	}

	if isSlice(oldValue) || isSlice(newValue) {
		return F(
			"[…]\n",
		)
	}

	return F("…\n")

}

func CodeMapDiff(c Context, newMap, oldMap map[string]any, form *forms.Form, path []string) Element {

	var extraContent Element

	values := []Element{}
	allKeys := []string{}
	qp := queryPath(c)

	for key, _ := range oldMap {
		allKeys = append(allKeys, key)
	}

	for key, _ := range newMap {
		if _, ok := oldMap[key]; ok {
			// we already have this key
			continue
		}
		allKeys = append(allKeys, key)
	}

	// we always sort keys
	sort.Strings(allKeys)
	for _, key := range allKeys {

		newValue := newMap[key]
		oldValue := oldMap[key]

		var validators []forms.Validator
		var field *forms.Field

		if form != nil {
			for _, f := range form.Fields {
				if f.Name == key {
					field = &f
					validators = f.Validators
					break
				}
			}
		}

		if field != nil && fieldIsSelected(append(path, key), qp) {
			if queryAction(c) != "" {
				extraContent = FieldEditorModal(c, field, sanitizeQueryPath(append(path, key)), nil)
			}
		}

		previousPath := path

		if len(previousPath) >= 1 {
			previousPath = previousPath[:len(previousPath)-1]
		}

		values = append(values, F(
			MapCodeValue(c, key, newValue, oldValue, validators, append(path, key)),
		))
	}

	return F(
		"{\n",
		values,
		strings.Repeat(" ", (len(path)+1)*2),
		"}\n",
		extraContent,
	)
}

func CodeDiff(c Context, newItem, oldItem *kodex.Item, form *forms.Form) Element {

	query := searchQuery(c)

	if query != "" {
		return Search(c, query, oldItem, newItem)
	}

	return Div(
		Class("kip-code-diff"),
		CodeMapDiff(c, newItem.All(), oldItem.All(), form, []string{}),
	)
}
