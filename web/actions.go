package web

import (
	"bytes"
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/web/ui"
	//	"github.com/kiprotect/kodex/api"
)

func Actions(project kodex.Project) ElementFunction {

	return func(c Context) Element {

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				Route("/details/(?P<actionId>[^/]+)", ActionDetails(project)),
				Route("", ActionsList(project)),
			),
		)
	}
}

func ActionDetails(project kodex.Project) func(c Context, actionId string) Element {

	return func(c Context, actionId string) Element {

		action, err := project.Controller().ActionConfig(Unhex(actionId))

		// make sure this action belongs to the project...
		if !bytes.Equal(action.Project().ID(), project.ID()) {
			return nil
		}

		if err != nil {
			return nil
		}

		AddBreadcrumb(c, action.Name(), Fmt("/details/%s", Hex(action.ID())))

		router := UseRouter(c)

		name := Var(c, action.Name())

		onSubmit := Func(c, func() {
			action.SetName(name.Get())
			action.Save()
			router.RedirectUp()
			kodex.Log.Info("Changing name to %s", name)
		})

		// edit the name of the action
		editActionName := func(c Context) Element {
			return Form(
				Method("POST"),
				OnSubmit(onSubmit),
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
			)
		}

		return Div(
			H2(
				Class("bulma-subtitle"),
				router.Match(
					c,
					Route("/name/edit",
						c.ElementFunction("editName", editActionName),
					),
					Route("",
						F(
							action.Name(),
							A(
								Href(router.CurrentRoute().Path+"/name/edit"),
								"&nbsp;&nbsp;",
								I(Class("fas fa-edit")),
							),
						),
					),
				),
			),
			P(
				Fmt("Type: %s", action.ActionType()),
			),
			c.Element("actionEditor", ActionEditor(action)),
		)
	}

}

func NewAction(project kodex.Project) ElementFunction {
	return func(c Context) Element {

		name := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)

		onSubmit := Func(c, func() {

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
				router.RedirectUp()
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

func ActionsList(project kodex.Project) ElementFunction {

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
				),
			)
			ais = append(ais, actionItem)
		}

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				Route("/new", c.Element("newAction", NewAction(project))),
				Route("", F(
					ui.List(ais),
					A(
						Href(router.CurrentRoute().Path+"/new"),
						Class("bulma-button", "bulma-is-success"),
						"New Action"),
				),
				),
			),
		)
	}
}
