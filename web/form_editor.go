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
) Element {

	var vv any
	var ok bool

	vv, ok = values[field.Name]

	if !ok {
		vv = ""
	}

	update := func(value any) {
		vv = value
		values[field.Name] = value
	}

	for _, validator := range validators {
		switch vt := validator.(type) {
		case forms.IsOptional:

			if vv == nil || vv == "" {
				update(vt.Default)
			}

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
				mapValue = map[string]any{}
			}

			if vt.Form != nil {
				return Div(
					Style("border: 1px solid #eee; padding: 10px; margin-top: 10px; margin-bottom: 10px;"),
					formAutoEditor(c, *vt.Form, mapValue, data, append(path, field.Name), nil),
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

			value := data.Var(name(path, field.Name), "")

			for i, choice := range vt.Choices {

				d := Fmt("%d", i)

				selected := false

				if value.Get() == d {
					// we mark this option as selected
					selected = true
				} else if value.Get() == "" && vv == choice {
					// this is the current value, we mark it as selected
					selected = true
				}

				if selected {
					// we update the current value so it can go through the
					// other validators correctly...
					update(choice)
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

			return validatorInput(c, field, data, path, switchValidators, values)

		}
	}

	return Div(Fmt("no UI available, sorry - %s", field.Name))
}

func copyMap(values map[string]any) map[string]any {
	newValues := make(map[string]any)

	for k, v := range values {
		newValues[k] = v
	}

	return newValues
}

func FormField(
	c Context,
	form forms.Form,
	field forms.Field,
	data *FormData,
	path []string,
	values map[string]any) Element {
	return Div(field.Name, validatorInput(c, field, data, path, field.Validators, values))
}

func FormAutoEditor(
	c Context,
	form forms.Form,
	values map[string]any,
	update func(map[string]any),
) Element {
	return formAutoEditor(c, form, values, nil, []string{}, update)
}

func applyValidators(field forms.Field, validators []forms.Validator, values map[string]any, newValues map[string]any, path []string, data *FormData) {

	formValues := data.Data()
	vv, _ := values[field.Name]

	update := func(value any) {
		vv = value
		values[field.Name] = value
	}

	for _, validator := range validators {
		switch vt := validator.(type) {
		case forms.IsOptional:

			if vv == nil || vv == "" {
				update(vt.Default)
			}

		case forms.IsString:

			// we initialize the variable with the existing value
			vs, _ := vv.(string)

			fullName := name(path, field.Name)

			if formValues.Has(fullName) {
				// if a new value exists in the form data, we replace vs
				vs = formValues.Get(fullName)
			}

			// we assign the value to the new data
			newValues[field.Name] = vs
		case forms.IsStringMap:

			vm, ok := vv.(map[string]any)

			if !ok {
				vm = map[string]any{}
			}

			if vt.Form != nil {
				// we recurse into the subform
				newValues[field.Name] = applyFormData(*vt.Form, copyMap(vm), append(path, field.Name), data)
			} else {
				// we directly assign the new value
				newValues[field.Name] = vm
			}
		case forms.IsInteger:
		case forms.IsFloat:
		case forms.IsBytes:
		case forms.IsIn:

			// we get the new index
			value := formValues.Get(name(path, field.Name))

			if value != "" {

				for i, choice := range vt.Choices {

					d := Fmt("%d", i)

					if value == d {
						// we update the value
						newValues[field.Name] = choice
						break
					}
				}

			} else if vv != nil {
				// we directly assign the existing value
				newValues[field.Name] = vv
			}

		case forms.Switch:

			keyValue, ok := newValues[vt.Key].(string)

			if !ok {
				continue
			}

			switchValidators, ok := vt.Cases[keyValue]

			if !ok {
				continue
			}

			applyValidators(field, switchValidators, values, newValues, path, data)

		}
	}

}

func applyFormData(form forms.Form, values map[string]any, path []string, data *FormData) map[string]any {
	newValues := map[string]any{}

	for _, field := range form.Fields {
		applyValidators(field, field.Validators, copyMap(values), newValues, path, data)
	}

	return newValues
}

func formAutoEditor(
	c Context,
	form forms.Form,
	values map[string]any,
	data *FormData,
	path []string,
	update func(map[string]any)) Element {

	// to do: use the ID of the form as a scope
	c = c.Scope("form")

	fields := make([]Element, 0)

	if data == nil {
		data = MakeFormData(c)
	}

	copiedValues := copyMap(values)

	for _, field := range form.Fields {
		fields = append(fields, FormField(c, form, field, data, path, copiedValues))
	}

	submitAction := NamedVar(c, "submitAction", "")

	onSubmit := Func[any](c, func() {

		if submitAction.Get() == "" {
			// this isn't a real update...
			return
		}

		newData := applyFormData(form, values, []string{}, data)

		if validatedData, err := form.Validate(newData); err != nil {
			// to do: proper error handling
			kodex.Log.Errorf("Invalid data: %v", err)
			return
		} else {
			// we update the original value with the validated data
			update(validatedData)
		}
	})

	if update != nil {
		return Form(
			Method("POST"),
			OnSubmit(onSubmit),
			fields,
			data,
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
		)
	} else {
		return F(fields)
	}

}
