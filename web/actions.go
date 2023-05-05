package web

import (
	"bytes"
	. "github.com/kiprotect/gospel"
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
				Route("/details/(?P<actionId>[^/]+)(?:/(?P<tab>edit|test))?", ActionDetails(project, onUpdate)),
				Route("", ActionsList(project, onUpdate)),
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
			content = Div("coming soon")
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
									"&nbsp;&nbsp;",
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
					"Type: ", "&nbsp;", B(action.ActionType()),
				),
			),
			Div(Class("bulma-content"), IfElse(action.Description() != "", action.Description(), "(no description given)")),
			ui.Tabs(
				ui.Tab(ui.ActiveTab(tab == "edit"), A(Href(Fmt("/projects/%s/actions/details/%s/edit", Hex(project.ID()), actionId)), "Edit")),
				ui.Tab(ui.ActiveTab(tab == "test"), A(Href(Fmt("/projects/%s/actions/details/%s/test", Hex(project.ID()), actionId)), "Test")),
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
					ui.List(
						ui.ListHeader(
							ui.ListColumn("md", "Name"),
							ui.ListColumn("sm", "Created At"),
							ui.ListColumn("sm", "Type"),
						),
						ais,
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
