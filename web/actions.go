// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package web

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"time"

	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions"
	"github.com/kiprotect/kodex/web/ui"
)

func TestWithItem(c Context, actionConfig kodex.ActionConfig, item int) Element {

	controller := UseController(c)

	action, err := actionConfig.Action()

	if err != nil {
		return Div("cannot get action")
	}

	var form *forms.Form

	// we check if this is a form action
	if formAction, ok := action.(*actions.FormAction); ok {
		form = formAction.Form()
	}

	data, ok := actionConfig.Data().(map[string]any)

	if !ok {
		return Div("cannot get data (not a map)")
	}

	// to do: improve parsing...

	rawItems, ok := data["items"].([]any)

	if !ok {
		return Div("cannot get data")
	}

	if len(rawItems) <= item {
		return Div("no data")
	}

	dataItem, ok := rawItems[item].(map[string]any)

	if !ok {
		return Div("cannot get data item")
	}

	dataItemData, ok := dataItem["data"].(map[string]any)

	if !ok {
		return Div("cannot get data item data")
	}

	items := []*kodex.Item{kodex.MakeItem(dataItemData)}

	parameterSet, err := kodex.MakeParameterSet([]kodex.Action{action}, controller.ParameterStore())

	writer := kodex.MakeInMemoryChannelWriter()

	processor, err := kodex.MakeProcessor(parameterSet, writer, nil)

	if err != nil {
		return Div("cannot create processor")
	}

	if newItems, err := processor.Process(items, nil); err != nil {
		return Div("Cannot process")
	} else {
		channels := make(map[string]interface{})
		channels["items"] = newItems
		for k, v := range writer.Items {
			channels[k] = v
		}

		// channels["errors"] = writer.Errors
		// channels["messages"] = writer.Messages
		// channels["warnings"] = writer.Warnings

		errors := []Element{}

		for _, error := range writer.Errors {
			errors = append(errors, Li(Fmt("%v", error.Error)))
		}

		if len(newItems) == 0 {
			return Div(
				Class("bulma-message", "bulma-is-danger"),
				Div(
					Class("bulma-message-body"),
					Ul(
						errors,
					),
				),
			)
		}

		var errorNotice Element

		formData := MakeFormData(c, "dataQuery", GET)
		query := formData.Var("query", "")

		return F(
			formData.Form(
				Fieldset(
					errorNotice,
					Div(
						Class("bulma-field"),
						P(
							Class("bulma-control"),
							Input(Class("bulma-control", "bulma-input"), Placeholder("filter fields by key or value"), Value(query)),
						),
					),
				),
			),
			Hr(),
			CodeDiff(c, newItems[0], items[0], form),
		)
	}

}

func ActionTest(actionConfig kodex.ActionConfig, onUpdate func(ChangeInfo, string)) ElementFunction {
	return func(c Context) Element {

		data, ok := actionConfig.Data().(map[string]any)

		if !ok {
			return Div("cannot get data (not a map)")
		}

		// to do: improve parsing...

		rawItems, ok := data["items"].([]any)

		if !ok {
			return Div("cannot get data")
		}

		router := UseRouter(c)
		dataItem := 0

		rv := router.Query().Get("dataItem")

		if len(rv) > 0 {
			var err error

			if dataItem, err = strconv.Atoi(string(rv[0])); err != nil {
				kodex.Log.Error(err)
			}
		}

		values := []Element{}

		for i, item := range rawItems {
			itemMap, ok := item.(map[string]any)

			if !ok {
				continue
			}

			values = append(values, Option(If(i == dataItem, BooleanAttrib("selected")()), Value(Fmt("%d", i)), itemMap["name"]))
		}

		content := TestWithItem(c, actionConfig, dataItem)

		return Div(
			Form(
				Id("itemForm"),
				Method("GET"),
				Div(
					Class("bulma-select", "bulma-is-fullwidth"),
					Select(
						values,
						Attrib("autocomplete")("off"),
						Id("itemSelect"),
						OnChange("itemForm.requestSubmit()"),
						Name("dataItem"),
					),
				),
			),
			Hr(),
			content,
		)
	}
}

func ActionData(action kodex.ActionConfig, onUpdate func(ChangeInfo, string)) ElementFunction {
	return func(c Context) Element {

		error := Var(c, "")
		router := UseRouter(c)

		// we fetch the existing data
		existingData := action.Data()

		dataMap, ok := existingData.(map[string]any)

		if !ok {
			dataMap = map[string]any{}
		}

		newData := map[string]any{}

		// we copy the existing map
		for k, v := range dataMap {
			newData[k] = v
		}

		dataItems, ok := newData["items"].([]any)

		if !ok {
			dataItems = []any{}
		}

		newDataItems := make([]any, len(dataItems))

		copy(newDataItems, dataItems)

		onSubmit := Func[any](c, func() {

			request := c.Request()

			file, header, err := request.FormFile("data")

			if err != nil {
				error.Set(Fmt("Cannot retrieve file: %v", err))
				return
			}

			content, err := io.ReadAll(file)

			if err != nil {
				error.Set(Fmt("Cannot read file: %v", err))
				return
			}

			var data map[string]any

			if err := json.Unmarshal(content, &data); err != nil {
				error.Set(Fmt("cannot unmarshal JSON: %v", err))
				return
			}

			newDataItems = append(newDataItems, map[string]any{
				"name": header.Filename,
				"data": data,
			})

			newData["items"] = newDataItems

			// we update the data map
			if err := action.SetData(newData); err != nil {
				error.Set(Fmt("Cannot set data: %v", err))
				return
			}

			if err := action.Save(); err != nil {
				error.Set(Fmt("Cannot save action: %v", err))
				return
			}

			onUpdate(ChangeInfo{}, router.CurrentPath())
		})

		items := []Element{}

		for i, dataItem := range dataItems {
			itemMap, ok := dataItem.(map[string]any)

			if !ok {
				continue
			}

			deleteDataItem := Func[any](c, func() {

				newData["items"] = append(newDataItems[:i], newDataItems[i+1:]...)

				// we update the data map
				if err := action.SetData(newData); err != nil {
					error.Set(Fmt("Cannot set data: %v", err))
					return
				}

				if err := action.Save(); err != nil {
					error.Set(Fmt("Cannot save action: %v", err))
					return
				}

				onUpdate(ChangeInfo{}, router.CurrentPath())

			})

			item := ui.ListItem(
				ui.ListColumn("md", itemMap["name"]),
				ui.ListColumn("icon",
					If(
						onUpdate != nil,
						Form(
							Method("POST"),
							OnSubmit(deleteDataItem),
							A(
								Href("#"),
								OnClick("this.parentElement.requestSubmit()"),
								Type("submit"),
								I(
									Class("fas", "fa-trash"),
								),
							),
						),
					),
				),
			)

			items = append(items, item)

		}

		return F(
			ui.List(
				ui.ListHeader(
					ui.ListColumn("md", "Name"),
					ui.ListColumn("icon", "Menu"),
				),
				items,
			),

			ui.MessageWithTitle(
				"grey",
				"Import Data",
				F(
					P(
						"You can import data from a JSON file.",
					),
					Br(),
					If(
						error.Get() != "",
						P(
							Class("bulma-help", "bulma-is-danger"),
							error.Get(),
						),
					),
					IfElse(
						onUpdate != nil,
						F(
							Form(
								Method("POST"),
								Enctype("multipart/form-data"),
								OnSubmit(onSubmit),
								Div(
									Id("data-file"),
									Class("bulma-file", "bulma-has-name"),
									Label(
										Class("bulma-file-label"),
										Input(
											Class("bulma-file-input"),
											Type("file"),
											Id("data"),
											Name("data"),
										),
										Span(
											Class("bulma-file-cta"),
											Span(
												Class("bulma-file-icon"),
												I(
													Class("fas", "fa-upload"),
												),
											),
											Span(
												Class("bulma-file-label"),
												"Info file...",
											),
										),
										Span(
											Class("bulma-file-name"),
											"please select a file",
										),
									),
								),
								Hr(),
								Button(
									Class("bulma-button", "bulma-is-success"),
									Type("submit"),
									"Import Data",
								),
							),
						),
						P(
							"You need to open a change request to import first.",
						),
					),
					Script(`
						const fileInput = document.querySelector('#data-file input[type=file]');
						  fileInput.onchange = () => {
						    if (fileInput.files.length > 0) {
						      const fileName = document.querySelector('#data-file .bulma-file-name');
						      fileName.textContent = fileInput.files[0].name;
						    }
						  }
					`),
				),
			),
		)
	}
}

func ActionDetails(project kodex.Project, onUpdate func(ChangeInfo, string)) func(c Context, actionId, tab string) Element {

	return func(c Context, actionId, tab string) Element {

		if tab == "" {
			tab = "edit"
		}

		action, err := project.Controller().ActionConfig(Unhex(actionId))

		if err != nil {
			return nil
		}

		// make sure this action belongs to the project...
		if !bytes.Equal(action.Project().ID(), project.ID()) {
			return nil
		}

		AddBreadcrumb(c, action.Name(), Fmt("/details/%s", Hex(action.ID())))

		router := UseRouter(c)

		name := Var(c, action.Name())
		error := Var(c, "")

		onSubmit := Func[any](c, func() {

			if name.Get() == "" {
				error.Set("please enter a name")
				return
			}

			if err := action.Update(map[string]any{"name": name.Get()}); err != nil {
				error.Set(Fmt("cannot set name: %v", err))
				return
			}

			if err := action.Save(); err != nil {
				error.Set(Fmt("cannot save: %v", err))
				return
			}

			onUpdate(ChangeInfo{}, router.LastPath())
		})

		var errorNotice Element

		if error.Get() != "" {
			errorNotice = P(
				Class("bulma-help", "bulma-is-danger"),
				error.Get(),
			)
		}

		// edit the name of the action
		editActionName := func(c Context) Element {
			return Form(
				Method("POST"),
				OnSubmit(onSubmit),
				Fieldset(
					errorNotice,
					Div(
						Class("bulma-field", "bulma-has-addons"),
						P(
							Class("bulma-control"),
							Input(Class("bulma-control", "bulma-input"), Value(name)),
						),
						P(
							Class("bulma-control"),
							Button(
								Class("bulma-button", "bulma-is-success"),
								Type("submit"),
								"Change",
							),
						),
					),
				),
			)
		}

		var content Element

		switch tab {
		case "edit":
			content = c.Element("actionEditor",
				ActionEditor(action, onUpdate),
			)
		case "test":
			content = c.Element("actionTest", ActionTest(action, onUpdate))
		case "data":
			content = c.Element("actionData", ActionData(action, onUpdate))
		}

		return Div(
			H2(
				Class("bulma-title"),
				router.Match(
					c,
					If(onUpdate != nil,
						Route("/name/edit",
							c.ElementFunction("editName", editActionName),
						),
					),
					Route("",
						F(
							action.Name(),
							If(onUpdate != nil,
								A(
									Style("float: right"),
									Href(router.CurrentRoute().Path+"/name/edit"),
									Nbsp,
									Nbsp,
									I(Class("fas fa-edit")),
								),
							),
						),
					),
				),
			),
			Div(
				Class("bulma-tags"),
				Span(
					Class("bulma-tag", "bulma-is-info", "bulma-is-light"),
					Fmt("last modified: %s", HumanDuration(time.Now().Sub(action.CreatedAt()))),
				),
				Span(
					Class("bulma-tag", "bulma-is-info", "bulma-is-light"),
					"Type: ", Nbsp, B(action.ActionType()),
				),
			),
			Div(Class("bulma-content"), IfElse(action.Description() != "", action.Description(), "(no description given)")),
			ui.Tabs(
				ui.Tab(ui.ActiveTab(tab == "edit"), A(Href(Fmt("/flows/projects/%s/actions/details/%s/edit", Hex(project.ID()), actionId)), "Edit")),
				ui.Tab(ui.ActiveTab(tab == "test"), A(Href(Fmt("/flows/projects/%s/actions/details/%s/test", Hex(project.ID()), actionId)), "Test")),
				ui.Tab(ui.ActiveTab(tab == "data"), A(Href(Fmt("/flows/projects/%s/actions/details/%s/data", Hex(project.ID()), actionId)), "Data")),
			),
			content,
		)
	}

}

func NewAction(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {
	return func(c Context) Element {

		name := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)

		onSubmit := Func[any](c, func() {

			if name.Get() == "" {
				error.Set("Please enter a name")
				return
			}

			action := project.MakeActionConfig(nil)

			action.SetName(name.Get())
			action.SetActionType("form")
			action.SetConfigData(map[string]any{
				"fields": []any{},
			})

			if err := action.Save(); err != nil {
				error.Set("Cannot save action")
			} else {
				onUpdate(ChangeInfo{}, router.CurrentPath())
			}
		})

		var errorNotice Element

		if error.Get() != "" {
			errorNotice = P(
				Class("bulma-help", "bulma-is-danger"),
				error.Get(),
			)
		}

		return Form(
			Method("POST"),
			OnSubmit(onSubmit),
			Div(
				Class("bulma-field"),
				errorNotice,
				Label(
					Class("bulma-label", "Name"),
					Input(
						Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
						Type("text"),
						Value(name),
						Placeholder("action name"),
					),
				),
			),
			Div(
				Class("bulma-field"),
				P(
					Class("bulma-control"),
					Button(
						Class("bulma-button", "bulma-is-success"),
						Type("submit"),
						"Create Action",
					),
					Nbsp,
					A(
						Class("bulma-button"),
						Href(Fmt("/flows/projects/%s/actions", Hex(project.ID()))),
						"Cancel",
					),
				),
			),
		)
	}
}

func Actions(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		// we retrieve the action configs of the project...
		actions, err := project.Controller().ActionConfigs(map[string]interface{}{
			"project.id": project.ID(),
		})

		if err != nil {
			// to do: error handling
			return nil
		}

		ais := make([]Element, 0, len(actions))

		for _, action := range actions {
			actionItem := A(
				Href(Fmt("/flows/projects/%s/actions/details/%s", Hex(project.ID()), Hex(action.ID()))),
				ui.ListItem(
					ui.ListColumn("md", action.Name()),
					ui.ListColumn("sm", HumanDuration(time.Now().Sub(action.CreatedAt()))),
					ui.ListColumn("sm", action.ActionType()),
				),
			)
			ais = append(ais, actionItem)
		}

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				If(onUpdate != nil, Route("/new", c.Element("newAction", NewAction(project, onUpdate)))),
				Route("", F(
					IfElse(
						len(ais) > 0,
						ui.List(
							ui.ListHeader(
								ui.ListColumn("md", "Name"),
								ui.ListColumn("sm", "Created At"),
								ui.ListColumn("sm", "Type"),
							),
							ais,
						),
						ui.Message(
							"info",
							"No existing actions.",
						),
					),
					If(onUpdate != nil,
						A(
							Href(router.CurrentRoute().Path+"/new"),
							Class("bulma-button", "bulma-is-success"),
							"New Action",
						),
					),
				),
				),
			),
		)
	}
}
