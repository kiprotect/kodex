package web

import (
	"encoding/json"
	"bytes"
	"io"
	. "github.com/gospel-dev/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/web/ui"
	"time"
	//	"github.com/kiprotect/kodex/api"
)

func Actions(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				Route("/details/(?P<actionId>[^/]+)(?:/(?P<tab>edit|test|data))?", ActionDetails(project, onUpdate)),
				Route("", ActionsList(project, onUpdate)),
			),
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

		dataItems, ok := dataMap["items"].([]any)

		if !ok {
			dataItems = []any{}
		}

		onSubmit := Func[any](c, func() {

			request := c.Request()

			file, header, err := request.FormFile("data")

			kodex.Log.Info(header.Filename)

			if err != nil {
				error.Set(Fmt("Cannot retrieve file: %v", err))
				return
			}

			content, err := io.ReadAll(file)

			if err != nil {
				error.Set(Fmt("Cannot read file: %v", err))
				return
			}

			error.Set(Fmt("file length: %d", len(content)))

			var data map[string]any

			if err := json.Unmarshal(content, &data); err != nil {
				error.Set(Fmt("cannot unmarshal JSON: %v", err))
				return
			}

			dataItems = []any{}

			dataItems = append(dataItems, map[string]any{
				"name": header.Filename,
				"data": data,
			})

			dataMap["items"] = dataItems

			// we update the data map
			if err := action.SetData(dataMap); err != nil {
				error.Set(Fmt("Cannot set data: %v", err))
				return
			}


			if err := action.Save(); err != nil {
				error.Set(Fmt("Cannot save action: %v", err))
				return
			}

			kodex.Log.Info("Success")

			onUpdate(ChangeInfo{}, router.CurrentPath())
		})

		items := []Element{}

		ed, _ := json.Marshal(dataMap)

		for _, dataItem := range dataItems {
			itemMap, ok := dataItem.(map[string]any)

			if !ok {
				continue
			}


			item := ui.ListItem(
				ui.ListColumn("md", itemMap["name"]),
				ui.ListColumn("sm", ""),
			)

			items = append(items, item)

		}

		return F(
			string(ed),
			ui.List(
				ui.ListHeader(
					ui.ListColumn("md", "Name"),
					ui.ListColumn("sm", "Type"),
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
						console.log("hey");
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
						Class("bHavelwelle,  14471 Brandenburg - PotsdamHavelwelle,  14471 Brandenburg - Potsdamulma-field", "bulma-has-addons"),
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
		case "data":
			content = c.Element("actionData", ActionData(action, onUpdate))
		}

		return Div(
			H2(
				Class("bulma-subtitle"),
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
				ui.Tab(ui.ActiveTab(tab == "edit"), A(Href(Fmt("/projects/%s/actions/details/%s/edit", Hex(project.ID()), actionId)), "Edit")),
				ui.Tab(ui.ActiveTab(tab == "test"), A(Href(Fmt("/projects/%s/actions/details/%s/test", Hex(project.ID()), actionId)), "Test")),
				ui.Tab(ui.ActiveTab(tab == "data"), A(Href(Fmt("/projects/%s/actions/details/%s/data", Hex(project.ID()), actionId)), "Data")),
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
				),
			),
		)
	}
}

func ActionsList(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {

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
				Href(Fmt("/projects/%s/actions/details/%s", Hex(project.ID()), Hex(action.ID()))),
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
