package web

import (
	"bytes"
	. "github.com/gospel-dev/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/web/ui"
	"time"
	//	"github.com/kiprotect/kodex/api"
)

func StreamConfigs(stream kodex.Stream, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				Route("/details/(?P<configId>[^/]+)(?:/(?P<tab>edit|test))?", StreamConfigDetails(stream, onUpdate)),
				Route("", StreamConfigsList(stream, onUpdate)),
			),
		)
	}
}

func StreamConfigDetails(stream kodex.Stream, onUpdate func(ChangeInfo, string)) func(c Context, configId, tab string) Element {

	return func(c Context, configId, tab string) Element {

		if tab == "" {
			tab = "edit"
		}

		config, err := stream.Config(configId)

		if err != nil {
			return nil
		}

		// make sure this config belongs to the stream...
		if !bytes.Equal(config.Stream().ID(), stream.ID()) {
			return nil
		}

		AddBreadcrumb(c, config.Name(), Fmt("/details/%s", Hex(config.ID())))

		router := UseRouter(c)

		name := Var(c, config.Name())
		error := Var(c, "")

		onSubmit := Func[any](c, func() {

			if name.Get() == "" {
				error.Set("please enter a name")
				return
			}

			if err := config.Update(map[string]any{"name": name.Get()}); err != nil {
				error.Set(Fmt("cannot set name: %v", err))
				return
			}

			if err := config.Save(); err != nil {
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

		// edit the name of the config
		editStreamConfigName := func(c Context) Element {
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
		}

		return Div(
			H2(
				Class("bulma-subtitle"),
				router.Match(
					c,
					If(onUpdate != nil,
						Route("/name/edit",
							c.ElementFunction("editName", editStreamConfigName),
						),
					),
					Route("",
						F(
							config.Name(),
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
					Fmt("last modified: %s", HumanDuration(time.Now().Sub(config.CreatedAt()))),
				),
			),
			Div(Class("bulma-content"), IfElse(config.Description() != "", config.Description(), "(no description given)")),
			ui.Tabs(
				ui.Tab(ui.ActiveTab(tab == "edit"), A(Href(Fmt("/streams/details/%s/configs/details/%s/edit", Hex(stream.ID()), configId)), "Edit")),
				ui.Tab(ui.ActiveTab(tab == "test"), A(Href(Fmt("/streams/details/%s/configs/details/%s/test", Hex(stream.ID()), configId)), "Test")),
			),
			content,
		)
	}

}

func NewStreamConfig(stream kodex.Stream, onUpdate func(ChangeInfo, string)) ElementFunction {
	return func(c Context) Element {

		name := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)

		onSubmit := Func[any](c, func() {

			if name.Get() == "" {
				error.Set("Please enter a name")
				return
			}

			config := stream.MakeConfig(nil)

			config.SetName(name.Get())

			if err := config.Save(); err != nil {
				error.Set("Cannot save config")
			} else {
				onUpdate(ChangeInfo{}, router.CurrentPath())
			}

			configs, _ := stream.Configs()

			kodex.Log.Infof("Configs: %s", configs)

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
						Placeholder("config name"),
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
						"Create StreamConfig",
					),
				),
			),
		)
	}
}

func StreamConfigsList(stream kodex.Stream, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		// we retrieve the config configs of the stream...
		configs, err := stream.Configs()

		if err != nil {
			// to do: error handling
			return nil
		}

		ais := make([]Element, 0, len(configs))

		for _, config := range configs {
			configItem := A(
				Href(Fmt("/projects/%s/streams/details/%s/configs/details/%s", Hex(stream.Project().ID()), Hex(stream.ID()), Hex(config.ID()))),
				ui.ListItem(
					ui.ListColumn("md", config.Name()),
					ui.ListColumn("sm", HumanDuration(time.Now().Sub(config.CreatedAt()))),
				),
			)
			ais = append(ais, configItem)
		}

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				If(onUpdate != nil, Route("/new", c.Element("newStreamConfig", NewStreamConfig(stream, onUpdate)))),
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
							"New StreamConfig",
						),
					),
				),
				),
			),
		)
	}
}
