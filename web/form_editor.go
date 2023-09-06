package web

import (
	. "github.com/gospel-dev/gospel"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"strings"
)

func fieldElement(input Element, errorNotice Element) Element {
	return Div(
		Class("bulma-field"),
		errorNotice,
		Label(
			Class("bulma-label", "Name"),
			input,
		),
	)
}

func name(path []string, name string) string {
	return strings.Join(append(path, name), ".")
}

func validatorInput(
	c Context,
	field forms.Field,
	data *FormData,
	path []string,
	validators []forms.Validator,
	values map[string]any,
	update func(string, any),
) Element {

	optional := false
	var defaultValue any

	var vv any
	var ok bool

	vv, ok = values[field.Name]

	if !ok {
		vv = ""
	}

	for _, validator := range validators {
		switch vt := validator.(type) {
		case forms.IsOptional:
			optional = true
			defaultValue = vt.Default
		case forms.IsString:

			// the variable value should be tied to the form so that it's
			// identifiable, i.e.

			vs, _ := vv.(string)

			value := data.Var(name(path, field.Name), vs)

			return fieldElement(
				Input(
					Class("bulma-input", If(false, "bulma-is-danger")),
					Type("text"),
					Value(value),
					Placeholder(field.Description),
				),
				nil,
			)
		case forms.IsStringMap:

			mapValue, ok := values[field.Name].(map[string]any)

			if !ok {
				if optional && defaultValue != nil {
					mapValue, ok = defaultValue.(map[string]any)

					if !ok {
						return Div("default value is not a map")
					}

				}

				if mapValue == nil {

					if optional {
						return Div("optional")
					}

					return Div(Fmt("not a map: %v - %s", values, field.Name))
				}

			}

			if vt.Form != nil {
				return Div(
					Style("border: 1px solid #eee; padding: 10px; margin-top: 10px; margin-bottom: 10px;"),
					formAutoEditor(c, *vt.Form, mapValue, data, append(path, field.Name), func(key string, value any) {

						// we build a new string map
						newValue := make(map[string]any)

						for k, v := range mapValue {
							newValue[k] = v
						}

						// we update the specific key in the string map
						newValue[key] = value

						// we update the entire string map value
						update(field.Name, newValue)
					}, false),
				)
			}
			return Div("map")
		case forms.IsInteger:
			return Div("integer")
		case forms.IsFloat:
			return Div("float")
		case forms.IsBytes:
			return Div("bytes")
		case forms.IsIn:

			options := make([]Element, 0)

			value := data.Var(name(path, field.Name), "0")

			value.OnUpdate(func() {
				kodex.Log.Infof("Selected value updated: %s: %v", field.Name, value.Get())
			})

			for i, choice := range vt.Choices {

				d := Fmt("%d", i)

				selected := false

				if value.Get() == d {
					// we update the value
					values[field.Name] = choice
					// we mark this option as selected
					selected = true
				}

				options = append(options,
					Option(
						Value(d),
						If(selected, BooleanAttrib("selected")()),
						choice,
					),
				)
			}

			return fieldElement(
				Div(
					Class("bulma-control", "bulma-is-expanded"),
					Div(
						Class("bulma-select", If(false, "bulma-is-danger")),
						Select(
							options,
							Value(value),
							OnChange("this.form.requestSubmit()"),
						),
					),
				),
				nil,
			)

		case forms.Switch:

			keyValue, ok := values[vt.Key].(string)

			if !ok {
				return Div("error: key value missing or not a string")
			}

			switchValidators, ok := vt.Cases[keyValue]

			if !ok {
				return Div("error: unknown switch case")
			}

			return validatorInput(c, field, data, path, switchValidators, values, update)

		}
	}

	return Div(Fmt("no UI available, sorry - %s", field.Name))
}

func FormField(
	c Context,
	form forms.Form,
	field forms.Field,
	data *FormData,
	path []string,
	values map[string]any,
	update func(string, any)) Element {
	/*
		- Render the UI of a given form element
	*/
	return Div(field.Name, validatorInput(c, field, data, path, field.Validators, values, update))
}

func FormAutoEditor(
	c Context,
	form forms.Form,
	values map[string]any,
	update func(string, any)) Element {
	return formAutoEditor(c, form, values, nil, []string{}, update, true)
}

func formAutoEditor(
	c Context,
	form forms.Form,
	values map[string]any,
	data *FormData,
	path []string,
	update func(string, any),
	topLevel bool) Element {

	c = c.Scope("form")

	fields := make([]Element, 0)

	if data == nil {
		data = MakeFormData(c)
	}

	for _, field := range form.Fields {
		fields = append(fields, FormField(c, form, field, data, path, values, update))
	}

	submitAction := NamedVar(c, "submitAction", "")

	onSubmit := Func[any](c, func() {

		if submitAction.Get() == "" {
			// this isn't a real update...
			return
		}

		/*
			- Map the field values of the form data structure to a real data structure
			  using the form definitions, combining them with the original data
			- Validate the new field values against the form, save any error data
			- Call the update function to update the data if it's valid
		*/

		kodex.Log.Infof("Submitted: %s", submitAction.Get())
	})

	if topLevel {
		return Form(
			Method("POST"),
			OnSubmit(onSubmit),
			fields,
			data,
			If(
				topLevel,
				Div(
					Class("bulma-field"),
					P(
						Class("bulma-control"),
						Button(
							Class("bulma-button", "bulma-is-success"),
							// we set the action to 'update'
							Value(submitAction, "update"),
							Type("submit"),
							"Submit",
						),
					),
				),
			),
		)
	} else {
		return F(fields)
	}

}
