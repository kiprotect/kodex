package web

import (
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"strconv"
	"strings"
)

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
	error map[string]any,
	readOnly bool,
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

	// we can possibly get a form definition from

	for _, validator := range validators {
		switch vt := validator.(type) {
		case forms.IsOptional:
			if vv == nil || vv == "" {
				update(vt.Default)
			}
		case forms.CanBeAnything:
			// to do: support non-string values

			vs, _ := vv.(string)

			value := data.Var(name(path, field.Name), vs)

			return Div(
				Class("bulma-field"),
				Label(
					Class("bulma-label"),
					field.Name,
				),
				Input(
					If(readOnly, BooleanAttrib("disabled")()),
					Class("bulma-input", If(false, "bulma-is-danger")),
					Type("text"),
					Value(value),
				),
				If(
					field.Description != "",
					P(
						Class("bulma-help"),
						field.Description,
					),
				),
			)

		case forms.IsString:

			// the variable value should be tied to the form so that it's
			// identifiable, i.e.

			vs, _ := vv.(string)

			value := data.Var(name(path, field.Name), vs)

			return Div(
				Class("bulma-field"),
				Label(
					Class("bulma-label"),
					field.Name,
				),
				Input(
					If(readOnly, BooleanAttrib("disabled")()),
					Class("bulma-input", If(false, "bulma-is-danger")),
					Type("text"),
					Value(value),
				),
				If(
					field.Description != "",
					P(
						Class("bulma-help"),
						field.Description,
					),
				),
			)

		case forms.IsStringMap:

			mapValue, ok := values[field.Name].(map[string]any)

			if !ok {
				mapValue = map[string]any{}
			}

			if vt.Form != nil {
				return Div(
					Style("border: 1px solid #eee; padding: 10px; margin-top: 10px; margin-bottom: 10px;"),
					formAutoEditor(c, *vt.Form, mapValue, data, nil, append(path, field.Name), readOnly),
				)
			}
			return Div("map")
		case forms.IsBoolean:

			vs := ""

			if vb, _ := vv.(bool); vb {
				vs = "true"
			}

			name := name(path, field.Name)
			value := data.Var(name, vs)

			return Div(
				Class("bulma-field"),
				Label(
					Class("bulma-checkbox"),
					Nbsp,
					Input(
						If(readOnly, BooleanAttrib("disabled")()),
						If(value.Get() == "true", BooleanAttrib("checked")()),
						Type("checkbox"),
						Name(name),
						Value("true"),
					),
					field.Name,
				),
				If(
					field.Description != "",
					P(
						Class("bulma-help"),
						field.Description,
					),
				),
			)
		case forms.IsInteger:

			// the variable value should be tied to the form so that it's
			// identifiable, i.e.

			vi, _ := vv.(int)

			value := data.Var(name(path, field.Name), Fmt("%d", vi))

			return Div(
				Class("bulma-field"),
				Label(
					Class("bulma-label"),
					field.Name,
				),
				Input(
					If(readOnly, BooleanAttrib("disabled")()),
					Class("bulma-input", If(false, "bulma-is-danger")),
					Type("number"),
					If(vt.HasMin, Attrib("min")(Fmt("%d", vt.Min))),
					If(vt.HasMax, Attrib("max")(Fmt("%d", vt.Max))),
					Value(value),
				),
				If(
					field.Description != "",
					P(
						Class("bulma-help"),
						field.Description,
					),
				),
			)
		case forms.IsFloat:

			// the variable value should be tied to the form so that it's
			// identifiable, i.e.

			vi, _ := vv.(float64)

			value := data.Var(name(path, field.Name), Fmt("%f", vi))

			return Div(
				Class("bulma-field"),
				Label(
					Class("bulma-label"),
					field.Name,
				),
				Input(
					If(readOnly, BooleanAttrib("disabled")()),
					Class("bulma-input", If(false, "bulma-is-danger")),
					Type("number"),
					If(vt.HasMin, Attrib("min")(Fmt("%f", vt.Min))),
					If(vt.HasMax, Attrib("max")(Fmt("%f", vt.Max))),
					Value(value),
				),
				If(
					field.Description != "",
					P(
						Class("bulma-help"),
						field.Description,
					),
				),
			)
		case forms.IsBytes:
			return Div("bytes")
		case forms.IsNil:
			// there's nothing to edit here
			return nil
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

			return Div(
				Class("bulma-field"),
				Label(
					Class("bulma-label"),
					field.Name,
				),
				Div(
					Class("bulma-control", "bulma-is-expanded"),
					Div(
						Class("bulma-select", If(false, "bulma-is-danger")),
						Select(
							If(readOnly, BooleanAttrib("disabled")()),
							options,
							Value(value),
							OnChange("this.form.requestSubmit()"),
						),
					),
				),
				If(
					field.Description != "",
					P(
						Class("bulma-help"),
						field.Description,
					),
				),
			)

		case forms.Switch:

			keyValue, _ := values[vt.Key].(string)
			switchValidators, ok := vt.Cases[keyValue]

			if !ok && vt.Default != nil {
				switchValidators = vt.Default
				ok = true
			}

			if !ok {
				// there's no validator for this case...
				return nil
			}

			var fieldError map[string]any

			if error != nil {
				fieldError, _ = error[vt.Key].(map[string]any)
			}

			return validatorInput(c, field, data, path, switchValidators, values, fieldError, readOnly)

		}
	}

	return nil
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
	values map[string]any,
	formError map[string]any,
	readOnly bool) Element {

	element := validatorInput(c, field, data, path, field.Validators, values, formError, readOnly)

	if formError != nil {
		fieldError, ok := formError[field.Name]
		if ok {
			err, ok := fieldError.(string)
			if ok {
				return Div(
					P(
						Class("bulma-help", "bulma-is-danger"),
						err,
					),
					element,
				)
			}
		}
	}

	return element

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
		case forms.CanBeAnything:
			// to do: support non-string values

			// we initialize the variable with the existing value
			vs, _ := vv.(string)

			fullName := name(path, field.Name)

			if formValues.Has(fullName) {
				// if a new value exists in the form data, we replace vs
				vs = formValues.Get(fullName)
			}

			// we assign the value to the new data
			newValues[field.Name] = vs

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
		case forms.IsBoolean:

			value := formValues.Get(name(path, field.Name))

			if value != "" {
				newValues[field.Name] = true
			} else {
				newValues[field.Name] = false
			}

		case forms.IsInteger:

			// we get the new index
			value := formValues.Get(name(path, field.Name))

			vd, err := strconv.Atoi(value)

			if err != nil {
				kodex.Log.Errorf("%v", err)
				continue
			}

			newValues[field.Name] = vd

		case forms.IsFloat:

			// we get the new index
			value := formValues.Get(name(path, field.Name))

			vd, err := strconv.ParseFloat(value, 64)

			if err != nil {
				kodex.Log.Errorf("%v", err)
				continue
			}

			newValues[field.Name] = vd

		case forms.IsBytes:
		case forms.IsNil:
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

			keyValue, _ := newValues[vt.Key].(string)
			switchValidators, ok := vt.Cases[keyValue]

			if !ok && vt.Default != nil {
				switchValidators = vt.Default
				ok = true
			}

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

func FormAutoEditor(
	c Context,
	form forms.Form,
	values map[string]any,
	update func(map[string]any),
) Element {

	// to do: use the ID of the form as a scope
	c = c.Scope("form")
	formData := MakeFormData(c, "edit", POST)
	error := Var[map[string]any](c, nil)
	submitAction := formData.Var("submitAction", "")

	onSubmit := func() {

		if update == nil {
			return
		}

		if submitAction.Get() == "" {
			// this isn't a real update...
			return
		}

		newData := applyFormData(form, values, []string{}, formData)

		if validatedData, err := form.Validate(newData); err != nil {
			// to do: proper error handling
			kodex.Log.Errorf("Invalid form data: %v", err)
			formError, ok := err.(*forms.FormError)
			if !ok {
				return
			}
			error.Set(formError.Data())
			return
		} else {
			// we update the original value with the validated data
			update(validatedData)
		}
	}

	formData.OnSubmit(onSubmit)

	fields := formAutoEditor(c, form, copyMap(values), formData, error.Get(), []string{}, update == nil)

	if fields == nil {
		return Div(
			"It seems this validator doesn't have any fields to edit.",
		)
	}

	return formData.Form(
		fields,
		Div(
			Class("bulma-field"),
			P(
				Class("bulma-control"),
				If(
					update != nil,
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
}

func formAutoEditor(
	c Context,
	form forms.Form,
	values map[string]any,
	data *FormData,
	error map[string]any,
	path []string,
	readOnly bool) Element {

	fields := make([]Element, 0)

	for _, field := range form.Fields {

		formField := FormField(c, form, field, data, path, values, error, readOnly)

		if formField == nil {
			continue
		}

		fields = append(fields, formField)

	}

	if len(fields) == 0 {
		return nil
	}

	if error != nil {
		return Div(
			Style("background: repeating-linear-gradient(45deg, #fff, #feee 10px, #fff 10px); border: 2px solid #faa; padding: 10px; margin-top: 10px; margin-bottom: 10px;"),
			fields,
		)
	}

	return F(fields)

}
