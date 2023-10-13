package web

import (
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/web/ui"
	"time"
)

func NewConfigAction(config kodex.Config, onUpdate func(ChangeInfo, string)) ElementFunction {
	return func(c Context) Element {

		actionId := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)

		actions, err := config.Stream().Project().Controller().ActionConfigs(map[string]any{})

		if err != nil {
			return Div("cannot load actions")
		}

		onSubmit := Func[any](c, func() {

			var action kodex.ActionConfig

			id := Unhex(actionId.Get())

			for _, possibleAction := range actions {
				if string(possibleAction.ID()) == string(id) {
					action = possibleAction
					break
				}
			}

			if action == nil {
				error.Set("invalid action")
				return
			}

			configActions, err := config.ActionConfigs()

			if err != nil {
				error.Set(Fmt("Cannot load config actions: %v", err))
				return
			}

			for _, configAction := range configActions {
				if string(configAction.ID()) == string(action.ID()) {
					error.Set(Fmt("action already in config"))
					return
				}
			}

			if err := config.AddActionConfig(action, len(configActions)); err != nil {
				error.Set(Fmt("Cannot add config action: %v", err))
				return
			}

			if err := config.Save(); err != nil {
				error.Set(Fmt("Cannot save config: %v", err))
				return
			}

			onUpdate(ChangeInfo{}, router.CurrentPath())

		})

		var errorNotice Element

		if error.Get() != "" {
			errorNotice = P(
				Class("bulma-help", "bulma-is-danger"),
				error.Get(),
			)
		}

		options := make([]Element, 0, len(actions))

		for _, action := range actions {
			options = append(options,
				Option(
					Value(Hex(action.ID())),
					action.Name(),
				),
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
					Div(
						Class("bulma-control", "bulma-is-expanded"),
						Div(
							Class("bulma-select", "bulma-is-fullwidth", If(error.Get() != "", "bulma-is-danger")),
							Select(
								options,
								Value(actionId),
							),
						),
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
						"Add Action",
					),
				),
			),
		)
	}
}

func ConfigActionsList(config kodex.Config, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		// we retrieve the config configs of the config...
		actions, err := config.ActionConfigs()

		if err != nil {
			// to do: error handling
			return nil
		}

		ais := make([]Element, 0, len(actions))

		for _, action := range actions {
			actionItem := A(
				Href(Fmt("/projects/%s/actions/details/%s", Hex(config.Stream().Project().ID()), Hex(action.ID()))),
				ui.ListItem(
					ui.ListColumn("md", action.Name()),
					ui.ListColumn("sm", HumanDuration(time.Now().Sub(action.CreatedAt()))),
				),
			)
			ais = append(ais, actionItem)
		}

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				If(onUpdate != nil, Route("/new", c.Element("newConfigAction", NewConfigAction(config, onUpdate)))),
				Route("", F(
					ui.List(
						ui.ListHeader(
							ui.ListColumn("md", "Name"),
							ui.ListColumn("sm", "Created At"),
						),
						ais,
					),
					If(onUpdate != nil,
						A(
							Href(router.CurrentRoute().Path+"/new"),
							Class("bulma-button", "bulma-is-success"),
							"Add Action",
						),
					),
				),
				),
			),
		)
	}
}
