package web

import (
	"encoding/json"
	"github.com/kiprotect/go-helpers/forms"
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions"
	"github.com/kiprotect/kodex/api"
	"reflect"
	"strconv"
)

func ActionEditor(action kodex.ActionConfig, onUpdate func(api.Change, string)) ElementFunction {
	return func(c Context) Element {

		kodex.Log.Infof("Config data: %v", action.ConfigData())

		var content Element

		switch action.ActionType() {
		case "form":
			content = FormEditor(c, action, onUpdate)
		}

		return Div(
			Class("bulma-card"),
			Div(
				Class("bulma-card-content"),
				H2(Class("bulma-subtitle"), "Action Editor"),
				content,
			),
		)
	}
}

func FormEditor(c Context, actionConfig kodex.ActionConfig, onUpdate func(api.Change, string)) Element {

	action, err := actionConfig.Action()

	if err != nil {
		return Div("err")
	}

	formAction, ok := action.(*actions.FormAction)

	if !ok {
		return Div("errr")
	}

	form := formAction.Form()

	onActionUpdate := func(change api.Change, path string) {

		bytes, err := json.Marshal(form)

		if err != nil {
			return
		}

		config := map[string]interface{}{}

		if err := json.Unmarshal(bytes, &config); err != nil {
			return
		}

		actionConfig.SetConfigData(config)

		// we update the project
		onUpdate(change, path)

	}

	if onUpdate == nil {
		onActionUpdate = nil
	}

	return Div(
		Class("kip-action-form"),
		FormFields(c, form, onActionUpdate, []string{"root"}),
	)
}

func NewField(c Context, form *forms.Form, path []string, onUpdate func(api.Change, string)) Element {

	name := Var(c, "")
	error := Var(c, "")
	router := UseRouter(c)

	onSubmit := Func(c, func() {

		if name.Get() == "" {
			error.Set("Please enter a name")
			return
		}

		for _, field := range form.Fields {
			if field.Name == name.Get() {
				error.Set("A field with the same name already exists")
				return
			}
		}

		form.Fields = append(form.Fields, forms.Field{
			Name:       name.Get(),
			Validators: []forms.Validator{},
		})

		onUpdate(api.Change{}, router.CurrentPathWithQuery())
	})

	var errorNotice Element

	if error.Get() != "" {
		errorNotice = P(
			Class("bulma-help", "bulma-is-danger"),
			error.Get(),
		)
	}

	return Form(
		Class("bulma-form"),
		Method("POST"),
		OnSubmit(onSubmit),
		Fieldset(
			errorNotice,
			Div(
				Class("bulma-field", "bulma-has-addons"),
				P(
					Class("bulma-control"),
					Style("flex-grow: 1"),
					Input(
						Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
						Value(name),
					),
				),
				P(
					Class("bulma-control"),
					Button(
						Class("bulma-button", "bulma-is-success"),
						Type("submit"),
						"add field",
					),
				),
			),
		),
	)
}

func typeOf(validator forms.Validator) string {
	if t := reflect.TypeOf(validator); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func Validators(c Context, field *forms.Field, path []string, onUpdate func(api.Change, string), selected bool) Element {

	router := UseRouter(c)

	validators := make([]Element, 0)

	for i, validator := range field.Validators {
		validators = append(validators,
			Li(
				A(
					Href(
						PathWithQuery(router.CurrentPath(), map[string][]string{
							"field": append(path, Fmt("%d", i)),
						}),
					),
					typeOf(validator),
				),
			),
		)
	}

	return Ul(
		Class("kip-validators"),
		validators,
		If(onUpdate != nil,
			Li(
				A(
					Href(PathWithQuery(router.CurrentPath(), map[string][]string{
						"field":  path,
						"action": []string{"addValidator"},
					})),
					I(
						Class("fa", "fa-plus"),
					),
				),
			),
		),
	)
}

func queryPath(c Context) []string {
	router := UseRouter(c)

	if field, ok := router.Query()["field"]; ok {
		return field
	}

	return nil
}

func queryAction(c Context) string {
	router := UseRouter(c)

	if field, ok := router.Query()["action"]; ok {
		return field[0]
	}
	return ""
}

func DeleteFieldNotice(c Context, form *forms.Form, field *forms.Field, path []string, onUpdate func(api.Change, string)) Element {

	router := UseRouter(c)

	onSubmit := Func(c, func() {

		newFields := []forms.Field{}

		for _, existingField := range form.Fields {

			if existingField.Name == field.Name {
				continue
			}

			newFields = append(newFields, existingField)

		}

		form.Fields = newFields

		onUpdate(api.Change{}, router.CurrentPathWithQuery())
	})

	return Li(
		Class("kip-item", "kip-is-danger"),
		Div(
			Class("kip-col", "kip-is-lg"),
			"Do you really want to delete this field?",
		),
		Div(
			Class("kip-col", "kip-is-icon"),
			Form(
				Method("POST"),
				OnSubmit(onSubmit),
				Div(
					Class("bulma-field", "bulma-is-grouped"),
					P(
						Class("bulma-control"),
						A(
							Class("bulma-button"),
							Href(router.CurrentPath()),
							"Cancel",
						),
					),
					P(
						Class("bulma-control"),
						Button(
							Class("bulma-button", "bulma-is-danger"),
							"Delete",
						),
					),
				),
			),
		),
	)
}

func StringMapValidator(c Context, validator *forms.IsStringMap, path []string, onUpdate func(api.Change, string)) Element {

	return Div(
		Style("flex-grow: 1;"),
		H2(
			Class("bulma-subtitle"),
			"New map<string,any> validator",
		),
		Input(Class("bulma-control"), Type("checkbox")),
		FormFields(c, validator.Form, onUpdate, path),
	)
}

func NewValidator(c Context, field *forms.Field, path []string, onUpdate func(api.Change, string)) Element {

	router := UseRouter(c)

	validatorType := Var(c, router.Query().Get("validatorType"))

	onSubmit := Func(c, func() {

		switch validatorType.Get() {
		case "IsStringMap":

			validator := forms.IsStringMap{
				Form: &forms.Form{
					Strict: true,
					Fields: []forms.Field{},
				},
			}

			field.Validators = append(field.Validators, validator)

			onUpdate(api.Change{}, router.CurrentPathWithQuery())
		}

		/*router.RedirectTo(PathWithQuery(router.CurrentPath(), map[string][]string{
			"field":         append(path, field.Name),
			"validatorType": []string{validatorType.Get()},
			"action":        []string{"addValidator"},
		}))*/
	})

	return Form(
		Method("POST"),
		OnSubmit(onSubmit),
		Div(
			H2(Class("bulma-subtitle"), "New Validator"),
			Div(
				Class("bulma-field", "bulma-has-addons"),
				Div(
					Class("bulma-control", "bulma-is-expanded"),
					Div(
						Class("bulma-select", "bulma-is-fullwidth"),
						Select(
							Option(Value("IsString"), "string"),
							Option(Value("IsStringMap"), "map<string,any>"),
							Value(validatorType),
						),
					),
				),
				Div(
					Class("bulma-control"),
					Button(
						Class("bulma-button", "bulma-is-primary"),
						"Continue",
					),
				),
			),
		),
	)
}

func Field(c Context, form *forms.Form, field *forms.Field, path []string, onUpdate func(api.Change, string)) Element {

	router := UseRouter(c)

	queryPath := queryPath(c)
	queryAction := queryAction(c)

	fullMatch := true
	match := true

	for i, pe := range queryPath {
		if i >= len(path) {
			// there are segments beyond this field
			fullMatch = false
			break
		} else if path[i] != pe {
			fullMatch = false
			match = false
			break
		}
	}

	if onUpdate != nil && fullMatch && queryAction == "delete" {
		return DeleteFieldNotice(c, form, field, path, onUpdate)
	}

	var extraContent Element

	var index int = -1

	if len(queryPath) > len(path) {
		// we get the validator index from the query path
		var err error

		if index, err = strconv.Atoi(queryPath[len(path)]); err != nil {
			// invalid index, we ignore...
			index = -1
		}
	}

	if fullMatch && queryAction == "addValidator" {
		extraContent = NewValidator(c, field, path, onUpdate)
	} else if match && index >= 0 && index < len(field.Validators) {

		validator := field.Validators[index]

		switch vt := validator.(type) {
		case *forms.IsStringMap:
			extraContent = FormFields(c, vt.Form, onUpdate, append(path, Fmt("%d", index)))
		default:
			extraContent = Span(Fmt("don't know: %v", vt))
		}
	}

	return F(
		Li(
			Class("kip-item"),
			Div(
				Class("kip-field-name", "kip-col", "kip-is-sm"),
				H3(
					field.Name,
				),
			),
			Div(
				Class("kip-col", "kip-is-md"),
				Validators(c, field, path, onUpdate, match),
			),
			Div(
				Class("kip-col", "kip-is-icon"),
				If(onUpdate != nil,
					A(
						Href(PathWithQuery(router.CurrentPath(), map[string][]string{
							"field":  path,
							"action": []string{"delete"},
						})),
						I(
							Class("fa", "fa-trash"),
						),
					),
				),
			),
			If(extraContent != nil, F(Div(Class("kip-break")), extraContent)),
		),
	)
}

func FormFields(c Context, form *forms.Form, onUpdate func(api.Change, string), path []string) Element {

	fields := []Element{}

	// we copy the fields as they might change during iteration
	// since we e.g. can delete a field in an action...
	fvs := form.Fields

	for i, _ := range fvs {
		field := &fvs[i]
		fields = append(fields, Field(c, form, field, append(path, field.Name), onUpdate))
	}

	return Div(
		Class("kip-form-config"),
		Ul(
			Class("kip-fields", "kip-top-level", "kip-list"),
			Li(
				Class("kip-item", "kip-is-header"),
				Div(
					Class("kip-col", "kip-is-sm"),
					"Name",
				),
				Div(
					Class("kip-col", "kip-is-md"),
					"Validators",
				),
				Div(
					Class("kip-col", "kip-is-icon"),
					"",
				),
			),
			fields,
			If(onUpdate != nil,
				Li(
					Class("kip-item"),
					NewField(c, form, path, onUpdate),
				),
			),
		),
	)
}
