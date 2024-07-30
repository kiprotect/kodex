package web

import (
	"fmt"
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/web/ui"
	"sort"
	"strings"
)

func FromTo(newValue, oldValue any) Element {
	return F(
		Span(Class("kip-from"), IfElse(oldValue == nil, L("&#8709;"), F(Fmt("%v", oldValue)))),
		Span(Class("kip-to"), IfElse(newValue == nil, L("&#8709;"), F(Fmt("%v", newValue)))),
	)

}

func KeyProps(c Context, key string, path []string, up bool) Element {

	router := UseRouter(c)
	icon := "fa-caret-down"
	routePath := path

	if up {
		if len(routePath) > 0 {
			routePath = routePath[:len(routePath)-1]
			icon = "fa-caret-up"
		}
	}

	return Span(
		Class("kip-key-props"),
		Style(BorderLeftWidth(Px(float64(len(path))*4.0))),
		A(
			Href(PathWithQuery(router.CurrentPath(), map[string][]string{
				"field":    routePath,
				"dataItem": router.Query()["dataItem"],
			})),
			I(
				Class("fa", icon),
			),
		),
		Span(Class("kip-key"),
			A(
				Href(PathWithQuery(router.CurrentPath(), map[string][]string{
					"field":    routePath,
					"action":   []string{"view"},
					"dataItem": router.Query()["dataItem"],
				})), key),
		),
	)
}

func SliceDiff(c Context, newValue, oldValue []any, validators []forms.Validator, path []string) Element {

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
			items = append(items, Li(
				KeyProps(c, key, newPath, false),
			))
			continue
		}

		extraContent := AnyDiff(c, nv, ov, listValidators, newPath)

		var item Element

		if extraContent == nil {
			extraContent = FromTo(nv, ov)
		}

		item = Li(
			KeyProps(c, key, newPath, true),
			Span(
				Class("kip-validator-list"),
				Span(Validators(c, listValidators, path, nil)),
			),
			Div(
				Class("kip-extra-content"),
				extraContent,
			),
		)

		items = append(items, item)
	}

	return Ul(
		Class("kip-slice-diff"),
		items,
	)
}

func isStringMap(value any) bool {
	if _, ok := value.(map[string]any); ok {
		return true
	} else {
		return false
	}
}

func toStringMap(value any) map[string]any {
	if v, ok := value.(map[string]any); ok {
		return v
	} else {
		return nil
	}
}

func toSlice(value any) []any {
	if v, ok := value.([]any); ok {
		return v
	} else {
		return nil
	}
}

func isSlice(value any) bool {
	if _, ok := value.([]any); ok {
		return true
	} else {
		return false
	}
}

func diffableValue(newValue, oldValue any) bool {
	return isStringMap(newValue) || isStringMap(oldValue) || isSlice(oldValue) || isSlice(newValue)
}

func AnyDiff(c Context, newValue, oldValue any, validators []forms.Validator, path []string) Element {

	if isStringMap(oldValue) || isStringMap(newValue) {

		var form *forms.Form

		for _, validator := range validators {
			if mapValidator, ok := validator.(*forms.IsStringMap); ok {
				form = mapValidator.Form
			}
		}

		return MapDiff(c, toStringMap(newValue), toStringMap(oldValue), form, path)
	}

	if isSlice(oldValue) || isSlice(newValue) {
		return SliceDiff(c, toSlice(newValue), toSlice(oldValue), validators, path)
	}

	// we don't return anything
	return nil
}

// checks if a given field is selected through the active query path
func fieldIsSelected(path []string, queryPath []string) bool {
	for i, pe := range queryPath {
		if i >= len(path) {
			if strings.HasPrefix(pe, "validator-") {
				// this is a validator, we ignore it
				continue
			}
			// there are non-validator segments beyond this field
			return false
		} else if path[i] != pe {
			// the query path doesn't match
			return false
		}
	}

	return true
}

func dataPathMatches(path []string, queryPath []string) bool {
	for i, pe := range path {
		if i >= len(queryPath) {
			// there are segments beyond this field
			return false
		} else if queryPath[i] != pe {
			return false
		}
	}

	return true
}

func MapValue(c Context, key string, newValue, oldValue any, validators []forms.Validator, path []string) Element {

	qp := queryPath(c)
	fullMatch := dataPathMatches(path, qp)
	diffable := diffableValue(newValue, oldValue)

	if !diffable {

		// uncomment to show full values
		/* if !fullMatch && len(path) <= len(qp) {
			return nil
		}*/

		return Li(
			KeyProps(c, key, path, fullMatch),
			Span(
				Class("kip-validator-list"),
				Span(Validators(c, validators, path, nil)),
			),
			Span(
				Class("kip-value-map"),
				FromTo(newValue, oldValue),
			),
		)
	}

	// we check if this field is selected
	if !fullMatch {
		// uncomment to show all values
		if len(path) <= len(qp) {
			return nil
		}
		return Li(
			KeyProps(c, key, path, false),
			Span(
				Class("kip-validator-list"),
				Span(Validators(c, validators, path, nil)),
			),
		)
	}

	extraContent := AnyDiff(c, newValue, oldValue, validators, path)

	return Li(
		KeyProps(c, key, path, true),
		Span(
			Class("kip-validator-list"),
			Span(Validators(c, validators, path, nil)),
		),
		If(
			extraContent != nil,
			Div(
				Class("kip-extra-content"),
				extraContent,
			),
		),
	)
}

func fieldName(path []string) string {
	ncs := []string{}

	for _, element := range path {
		if strings.HasPrefix(element, "validator-") {
			continue
		}
		ncs = append(ncs, element)
	}
	return strings.Join(ncs, "/")
}

func FieldEditorModal(c Context, field *forms.Field, path []string, onUpdate func(ChangeInfo, string)) Element {

	router := UseRouter(c)
	queryAction := queryAction(c)

	if queryAction == "" {
		return nil
	}

	newPath := PathWithQuery(router.CurrentPath(), map[string][]string{
		"field": path,
	})

	validatorsView := ValidatorsActions(c, field.Validators,
		func(newValidator forms.Validator) int {
			field.Validators = append(field.Validators, newValidator)
			return len(field.Validators) - 1
		},
		func(index int, newValidator forms.Validator) error {

			if index >= len(field.Validators) {
				return fmt.Errorf("out of bounds: %d", index)
			}

			if newValidator == nil {
				// we delete the validator
				cv := field.Validators
				field.Validators = append(cv[:index], cv[index+1:]...)
				return nil
			}

			field.Validators[index] = newValidator
			return nil

		}, func(fromIndex, toIndex int) error {
			field.Validators[fromIndex], field.Validators[toIndex] = field.Validators[toIndex], field.Validators[fromIndex]
			return nil
		}, path, nil)

	content := F(
		Validators(c, field.Validators, path, nil),
		Hr(),
		validatorsView,
	)

	return ui.Modal(c, Fmt("%s", fieldName(path)), content, nil, newPath)
}

func ItemDiff(c Context, newItem, oldItem *kodex.Item, form *forms.Form) Element {
	//	return DataEditor(c, oldItem.All(), nil, []string{})

	query := searchQuery(c)

	if query != "" {
		return Search(c, query, newItem)
	}

	return MapDiff(c, newItem.All(), oldItem.All(), form, []string{})
}

func MapDiff(c Context, newMap, oldMap map[string]any, form *forms.Form, path []string) Element {

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

		var extraContent Element

		if field != nil && fieldIsSelected(append(path, key), qp) {
			extraContent = FieldEditorModal(c, field, append(path, key), nil)
		}

		previousPath := path

		if len(previousPath) >= 1 {
			previousPath = previousPath[:len(previousPath)-1]
		}

		values = append(values, F(
			MapValue(c, key, newValue, oldValue, validators, append(path, key)),
			extraContent,
		))
	}

	return Ul(
		Class("kip-map-diff", If(len(path) == 0, "kip-top-level")),
		If(
			len(path) == 0,
			Li(
				Class("kip-header"),
				Span(
					Class("kip-key-props"),
					"key",
				),
				Span(
					Class("kip-validator-list"),
					"validators",
				),
				Span(
					Class("kip-value-map"),
					Span(
						Class("kip-from"),
						"source",
					),
					Span(
						Class("kip-to"),
						"transformed",
					),
				),
			),
		),
		values,
	)

}

func searchForKey(query string, data any, path []string) [][]string {

	results := [][]string{}

	switch dt := data.(type) {
	case map[string]any:
		for k, v := range dt {
			if strings.Contains(k, query) {
				results = append(results, append(path, k))
			}
			results = append(results, searchForKey(query, v, append(path, k))...)
		}
	case []any:
		for i, v := range dt {
			results = append(results, searchForKey(query, v, append(path, Fmt("%d", i)))...)
		}
	default:
		return nil
	}

	return results
}

func Search(c Context, query string, newItem *kodex.Item) Element {
	router := UseRouter(c)
	matches := searchForKey(query, newItem.All(), []string{})

	matchesItems := []Element{}

	for _, match := range matches {
		matchesItems = append(matchesItems, Li(
			A(
				Href(
					PathWithQuery(router.CurrentPath(), map[string][]string{
						"field":  match,
						"action": []string{"view"},
					}),
				),
				fieldName(match),
			),
		))
	}

	return Ul(
		matchesItems,
	)

}
