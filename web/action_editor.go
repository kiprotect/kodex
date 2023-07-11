package web

import (
	"encoding/json"
	"fmt"
	. "github.com/gospel-dev/gospel"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions"
	"github.com/kiprotect/kodex/web/ui"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func ActionEditor(action kodex.ActionConfig, onUpdate func(ChangeInfo, string)) ElementFunction {
	return func(c Context) Element {

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

func SetFormAction(c Context, action *actions.FormAction) {
	GlobalVar(c, "formAction", action)
}

func UseFormAction(c Context) *actions.FormAction {
	return UseGlobal[*actions.FormAction](c, "formAction")
}

func FormEditor(c Context, actionConfig kodex.ActionConfig, onUpdate func(ChangeInfo, string)) Element {

	action, err := actionConfig.Action()

	if err != nil {
		return Div("err")
	}

	formAction, ok := action.(*actions.FormAction)

	if !ok {
		return Div("Error")
	}

	SetFormAction(c, formAction)

	form := formAction.Form()

	onActionUpdate := func(change ChangeInfo, path string) {

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

func NewField(c Context, form *forms.Form, path []string, onUpdate func(ChangeInfo, string)) Element {

	name := Var(c, "")
	error := Var(c, "")
	router := UseRouter(c)

	onSubmit := Func[any](c, func() {

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

		onUpdate(ChangeInfo{
			Description: Fmt("Create a new field with name '%s' at path '%s'.", name.Get(), strings.Join(path, ".")),
		},
			router.CurrentPathWithQuery())
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

func Validators(c Context, validators []forms.Validator, path []string, onUpdate func(ChangeInfo, string)) Element {

	router := UseRouter(c)

	elements := make([]Element, 0)

	for i, validator := range validators {
		elements = append(elements,
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
		elements,
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

func DeleteFieldNotice(c Context, form *forms.Form, field *forms.Field, path []string, onUpdate func(ChangeInfo, string)) Element {

	router := UseRouter(c)

	onSubmit := Func[any](c, func() {

		newFields := []forms.Field{}

		for _, existingField := range form.Fields {

			if existingField.Name == field.Name {
				continue
			}

			newFields = append(newFields, existingField)

		}

		form.Fields = newFields

		onUpdate(ChangeInfo{}, router.CurrentPathWithQuery())
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

func NewValidator(c Context, create func(validator forms.Validator) int, path []string, onUpdate func(ChangeInfo, string)) Element {

	router := UseRouter(c)

	validatorType := Var(c, router.Query().Get("validatorType"))

	action := UseFormAction(c)

	onSubmit := Func[any](c, func() {

		var validator forms.Validator

		validator, err := forms.ValidatorFromDescription(&forms.ValidatorDescription{
			Type:   validatorType.Get(),
			Config: map[string]any{},
		}, action.Context())

		if err != nil {
			kodex.Log.Error(err)
			return
		}

		if validator != nil {
			index := create(validator)
			onUpdate(ChangeInfo{}, router.CurrentPathWithQuery())
			router.RedirectTo(PathWithQuery(router.CurrentPath(), map[string][]string{
				"field": append(path, Fmt("%d", index)),
			}))
		}
	})

	values := []any{}

	vts := []string{}

	for vt, _ := range action.Context().Validators {
		vts = append(vts, vt)
	}

	sort.Strings(vts)

	for _, vt := range vts {
		values = append(values, Option(Value(vt), vt))
	}

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
							values,
							Value(validatorType),
						),
					),
				),
				Div(
					Class("bulma-control"),
					Button(
						Class("bulma-button", "bulma-is-primary"),
						"add validator",
					),
				),
			),
		),
	)
}

func ValidatorEditor(c Context, update func(validator forms.Validator) error, validator forms.Validator, path []string, onUpdate func(ChangeInfo, string)) Element {

	configJson, err := json.MarshalIndent(validator, "", "  ")

	if err != nil {
		return Div("Error serializing validator config")
	}

	config := Var(c, string(configJson))
	error := Var(c, "")
	router := UseRouter(c)

	if onUpdate == nil {
		return Pre(
			config.Get(),
		)
	}

	onSubmit := Func[any](c, func() {

		if config.Get() == "" {
			error.Set("Please enter a config")
			return
		}

		var configMap map[string]any

		if err := json.Unmarshal([]byte(config.Get()), &configMap); err != nil {
			error.Set(Fmt("Invalid JSON: %v", err))
			return
		}

		validatorType := forms.GetType(validator)

		action := UseFormAction(c)

		newValidator, err := forms.ValidatorFromDescription(&forms.ValidatorDescription{
			Type:   validatorType,
			Config: configMap,
		}, action.Context())

		if err != nil {
			error.Set(Fmt("Cannot update validator: %v", err))
			return
		}

		update(newValidator)

		onUpdate(ChangeInfo{
			Description: Fmt("Update validator config with value '%s' at path '%s'.", config.Get(), strings.Join(path, ".")),
		},
			router.CurrentPathWithQuery())
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
				Class("bulma-field"),
				P(
					Class("bulma-control"),
					Style("flex-grow: 1"),
					Div(
						Class("bulma-control"),
						Textarea(
							Class("bulma-textarea"),
							Attrib("rows")(Fmt("%d", 5)),
							Value(config),
						),
					),
				),
				Br(),
				P(
					Class("bulma-control"),
					Button(
						Class("bulma-button", "bulma-is-success"),
						Type("submit"),
						"update config",
					),
				),
			),
		),
	)
}

func SwitchValidator(c Context, validator *forms.Switch, onUpdate func(ChangeInfo, string), path []string) Element {

	key := Var(c, validator.Key)
	newCase := Var(c, "")
	formError := Var(c, "")
	router := UseRouter(c)

	queryPath := queryPath(c)
	queryAction := queryAction(c)

	onChangeKey := Func[any](c, func() {

		if key.Get() == "" {
			formError.Set("Please enter a key")
			return
		}

		validator.Key = key.Get()

		onUpdate(ChangeInfo{
			Description: "Update switch key",
		},
			router.CurrentPathWithQuery())
	})

	onAddCase := Func[any](c, func() {

		if newCase.Get() == "" {
			formError.Set("Please enter a case value")
			return
		}

		validator.Cases[newCase.Get()] = []forms.Validator{}

		onUpdate(ChangeInfo{
			Description: "Update switch key",
		},

			router.CurrentPathWithQuery())
	})

	var errorNotice Element

	if formError.Get() != "" {
		errorNotice = P(
			Class("bulma-help", "bulma-is-danger"),
			formError.Get(),
		)
	}

	cases := []Element{}

	caseValues := []string{}

	for caseValue, _ := range validator.Cases {
		caseValues = append(caseValues, caseValue)
	}

	sort.Strings(caseValues)

	for _, caseValue := range caseValues {

		casePath := append(path, caseValue)
		_, fullMatch := pathMatches(casePath, queryPath)

		if onUpdate != nil && fullMatch && queryAction == "delete" {

			url := PathWithQuery(router.CurrentPath(), map[string][]string{
				"field": path,
			})

			onSubmit := Func[any](c, func() {

				delete(validator.Cases, caseValue)

				onUpdate(ChangeInfo{}, url)
			})

			item := Li(
				Class("kip-item", "kip-is-danger"),
				Div(
					Class("kip-col", "kip-is-lg"),
					"Do you really want to delete the case \"",
					Strong(caseValue),
					"\"?",
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
									Href(url),
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
			cases = append(cases, item)
			continue
		}

		validators := validator.Cases[caseValue]

		var extraContent = ValidatorsActions(c, validators, func(newValidator forms.Validator) int {
			validator.Cases[caseValue] = append(validator.Cases[caseValue], newValidator)
			return len(validator.Cases[caseValue]) - 1
		}, func(index int, newValidator forms.Validator) error {
			if index >= len(validator.Cases[caseValue]) {
				return fmt.Errorf("out of bounds: %d", index)
			}

			if newValidator == nil {
				// we delete the validator
				cv := validator.Cases[caseValue]
				validator.Cases[caseValue] = append(cv[:index], cv[index+1:]...)
				return nil
			}

			validator.Cases[caseValue][index] = newValidator
			return nil
		}, casePath, onUpdate)

		item := Li(
			Class("kip-item", If(extraContent != nil, "kip-with-extra-content")),
			Div(
				Class("kip-field-name", "kip-col", "kip-is-sm"),
				H3(
					caseValue,
				),
			),
			Div(
				Class("kip-col", "kip-is-md"),
				Validators(c, validators, casePath, onUpdate),
			),
			Div(
				Class("kip-col", "kip-is-icon"),
				If(onUpdate != nil,
					A(
						Href(PathWithQuery(router.CurrentPath(), map[string][]string{
							"field":  casePath,
							"action": []string{"delete"},
						})),
						I(
							Class("fa", "fa-trash"),
						),
					),
				),
			),
			If(
				extraContent != nil,
				Div(
					Class("kip-extra-content"),
					extraContent,
				),
			),
		)

		cases = append(cases, item)

	}

	return F(
		Form(
			Class("bulma-form"),
			Method("POST"),
			OnSubmit(onChangeKey),
			Fieldset(
				errorNotice,
				Div(
					Class("bulma-field", "bulma-has-addons"),
					P(
						Class("bulma-control"),
						Style("flex-grow: 1"),
						Input(
							Class("bulma-input", If(formError.Get() != "", "bulma-is-danger")),
							Value(key),
						),
					),
					P(
						Class("bulma-control"),
						Button(
							Class("bulma-button", "bulma-is-success"),
							Type("submit"),
							"update key",
						),
					),
				),
			),
		),
		Div(
			Class("kip-form-config"),
			Ul(
				Class("kip-fields", "kip-list"),
				Li(
					Class("kip-item", "kip-is-header"),
					Div(
						Class("kip-col", "kip-is-sm"),
						"Case",
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
				cases,
				DoIf(onUpdate != nil,
					func() Element {
						return Li(
							Class("kip-item"),
							Form(
								Class("bulma-form"),
								Method("POST"),
								OnSubmit(onAddCase),
								Fieldset(
									errorNotice,
									Div(
										Class("bulma-field", "bulma-has-addons"),
										P(
											Class("bulma-control"),
											Style("flex-grow: 1"),
											Input(
												Class("bulma-input", If(formError.Get() != "", "bulma-is-danger")),
												Value(newCase),
											),
										),
										P(
											Class("bulma-control"),
											Button(
												Class("bulma-button", "bulma-is-success"),
												Type("submit"),
												"add case",
											),
										),
									),
								),
							),
						)
					},
				),
			),
		),
	)

}

func IsListValidator(c Context, validator *forms.IsList, onUpdate func(ChangeInfo, string), path []string) Element {
	return Div(
		Validators(c, validator.Validators, path, onUpdate),
		ValidatorsActions(c, validator.Validators, func(newValidator forms.Validator) int {
			validator.Validators = append(validator.Validators, newValidator)
			return len(validator.Validators) - 1
		}, func(index int, newValidator forms.Validator) error {
			if index >= len(validator.Validators) {
				return fmt.Errorf("out of bounds: %d", index)
			}

			if newValidator == nil {
				// we delete the validator
				cv := validator.Validators
				validator.Validators = append(cv[:index], cv[index+1:]...)
				return nil
			}

			validator.Validators[index] = newValidator
			return nil
		}, path, onUpdate),
	)
}

func ValidatorDetails(c Context, validator forms.Validator, update func(validator forms.Validator) error, path []string, onUpdate func(ChangeInfo, string)) Element {

	var content Element

	queryPath := queryPath(c)
	queryAction := queryAction(c)

	_, fullMatch := pathMatches(path, queryPath)

	router := UseRouter(c)

	if onUpdate != nil && fullMatch && queryAction == "delete" {

		url := PathWithQuery(router.CurrentPath(), map[string][]string{
			"field": path[:len(path)-1],
		})

		onSubmit := Func[any](c, func() {
			update(nil)
			onUpdate(ChangeInfo{}, url)
		})

		return ui.Message(
			"danger",
			F(
				"Do you really want to delete this validator?",
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
									Href(url),
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
			),
		)
	}

	switch vt := validator.(type) {
	case *forms.IsList:
		content = IsListValidator(c, vt, onUpdate, path)
	case *forms.Switch:
		content = SwitchValidator(c, vt, onUpdate, path)
	case *forms.IsStringMap:

		// we always create a form
		if vt.Form == nil {
			vt.Form = &forms.Form{}
		}
		content = FormFields(c, vt.Form, onUpdate, path)
	default:
		content = ValidatorEditor(c, update, validator, path, onUpdate)
	}

	return F(
		H2(
			Class("bulma-subtitle"),
			F(
				forms.GetType(validator),
				If(onUpdate != nil,
					A(
						Style("float: right"),
						Href(
							PathWithQuery(router.CurrentPath(), map[string][]string{
								"field":  path,
								"action": []string{"delete"},
							}),
						),
						Nbsp,
						Nbsp,
						I(Class("fas fa-trash")),
					),
				),
			),
		),

		content,
	)

}

func pathMatches(path []string, queryPath []string) (bool, bool) {

	fullMatch := true
	match := true

	for i, pe := range path {
		if i >= len(queryPath) {
			// there are segments beyond this field
			fullMatch = false
			break
		} else if queryPath[i] != pe {
			fullMatch = false
			match = false
			break
		}
	}

	if len(queryPath) != len(path) {
		fullMatch = false
	}

	return match, fullMatch

}

func ValidatorsActions(c Context, validators []forms.Validator, create func(validator forms.Validator) int, update func(index int, validator forms.Validator) error, path []string, onUpdate func(ChangeInfo, string)) Element {

	queryPath := queryPath(c)
	queryAction := queryAction(c)

	partialMatch, fullMatch := pathMatches(path, queryPath)

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
		return NewValidator(c, create, path, onUpdate)
	} else if partialMatch && index >= 0 && index < len(validators) {
		return ValidatorDetails(c, validators[index], func(validator forms.Validator) error {
			return update(index, validator)
		}, append(path, Fmt("%d", index)), onUpdate)
	}

	return nil
}

func Field(c Context, form *forms.Form, field *forms.Field, path []string, onUpdate func(ChangeInfo, string)) Element {

	router := UseRouter(c)

	queryPath := queryPath(c)
	queryAction := queryAction(c)

	_, fullMatch := pathMatches(path, queryPath)

	if onUpdate != nil && fullMatch && queryAction == "delete" {
		return DeleteFieldNotice(c, form, field, path, onUpdate)
	}

	extraContent := ValidatorsActions(c, field.Validators,
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

		}, path, onUpdate)

	return Li(
		Class("kip-item", If(extraContent != nil, "kip-with-extra-content")),
		Div(
			Class("kip-field-name", "kip-col", "kip-is-sm"),
			H3(
				field.Name,
			),
		),
		Div(
			Class("kip-col", "kip-is-md"),
			Validators(c, field.Validators, path, onUpdate),
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
		If(
			extraContent != nil,
			Div(
				Class("kip-extra-content"),
				extraContent,
			),
		),
	)
}

func FormFields(c Context, form *forms.Form, onUpdate func(ChangeInfo, string), path []string) Element {

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
			Class("kip-fields", "kip-list", If(len(path) == 1, "kip-top-level")),
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
			DoIf(onUpdate != nil,
				func() Element {
					return Li(
						Class("kip-item"),
						NewField(c, form, path, onUpdate),
					)
				},
			),
		),
	)
}
