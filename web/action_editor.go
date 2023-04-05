package web

import (
	"encoding/json"
	"github.com/kiprotect/go-helpers/forms"
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions"
)

func ActionEditor(action kodex.ActionConfig) ElementFunction {
	return func(c Context) Element {

		kodex.Log.Infof("Config data: %v", action.ConfigData())

		var content Element

		switch action.ActionType() {
		case "form":
			content = FormEditor(c, action)
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

func FormEditor(c Context, actionConfig kodex.ActionConfig) Element {

	action, err := actionConfig.Action()

	if err != nil {
		return Div("err")
	}

	formAction, ok := action.(*actions.FormAction)

	if !ok {
		return Div("errr")
	}

	form := formAction.Form()

	onUpdate := func() {

		bytes, err := json.Marshal(form)

		if err != nil {
			return
		}

		config := map[string]interface{}{}

		if err := json.Unmarshal(bytes, &config); err != nil {
			return
		}

		actionConfig.SetConfigData(config)
	}

	return Div(
		FormFields(c, actionConfig, form, onUpdate),
	)
}

func NewField(c Context, form *forms.Form, onUpdate func()) Element {

	name := Var(c, "")
	router := UseRouter(c)

	onSubmit := Func(c, func() {
		kodex.Log.Info("submitting...")
		form.Fields = append(form.Fields, forms.Field{
			Name:       name.Get(),
			Validators: []forms.Validator{},
		})
		onUpdate()
		router.RedirectTo(router.CurrentPath())
	})

	return Form(
		Class("bulma-form"),
		Method("POST"),
		OnSubmit(onSubmit),
		Fieldset(
			Div(
				Class("bulma-field", "bulma-has-addons"),
				P(
					Class("bulma-control"),
					Style("flex-grow: 1"),
					Input(
						Class("bulma-input"),
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

func FormFields(c Context, actionConfig kodex.ActionConfig, form *forms.Form, onUpdate func()) Element {

	fields := []Element{}

	for _, field := range form.Fields {
		fields = append(fields, Div(field.Name))
	}

	return Div(
		fields,
		NewField(c, form, onUpdate),
	)
}
