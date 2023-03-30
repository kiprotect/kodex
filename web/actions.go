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
				Route("/(?P<actionId>[^/]+)", ActionDetails(project)),
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

		router := UseRouter(c)

		name := Var(c, action.Name())

		onSubmit := Func(c, func() {
			action.SetName(name.Get())
			action.Save()
			router.RedirectUp()
			kodex.Log.Info("Changing name to %s", name)
		})

		editActionName := func(c Context) Element {
			return F(
				Form(
					Method("POST"),
					OnSubmit(onSubmit),
					Input(Class("bulma-control", "bulma-input"), Value(name)),
					Button(
						Class("bulma-button", "bulma-is-success"),
						Type("submit"),
						"Change",
					),
				),
			)
		}

		return Div(
			H1(
				Class("bulma-title"),
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
								"&nbsp;",
								I(Class("fas fa-edit")),
							),
						),
					),
				),
			),
			P(
				Fmt("Type: %s", action.ActionType()),
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
				Href(Fmt("/projects/%s/actions/%s", Hex(project.ID()), Hex(action.ID()))),
				ui.ListItem(
					ui.ListColumn("md", action.Name()),
				),
			)
			ais = append(ais, actionItem)
		}

		return F(
			ui.List(ais),
			Button(Class("bulma-button", "bulma-is-success"), "New Action"),
		)
	}
}
