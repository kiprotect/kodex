package web

import (
	"bytes"
	. "github.com/gospel-dev/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/web/ui"
	"time"
	//	"github.com/kiprotect/kodex/api"
)

func Streams(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				Route("", StreamsList(project, onUpdate)),
			),
		)
	}
}

func StreamDetails(project kodex.Project, onUpdate func(ChangeInfo, string)) func(c Context, streamId, tab string) Element {

	return func(c Context, streamId, tab string) Element {

		if tab == "" {
			tab = "configs"
		}

		stream, err := project.Controller().Stream(Unhex(streamId))

		if err != nil {
			return Div(Fmt("can't find stream: %v (%s)", err, streamId))
		}

		// make sure this stream belongs to the project...
		if !bytes.Equal(stream.Project().ID(), project.ID()) {
			return nil
		}

		AddBreadcrumb(c, stream.Name(), Fmt("/details/%s", Hex(stream.ID())))

		router := UseRouter(c)

		name := Var(c, stream.Name())
		error := Var(c, "")

		onSubmit := Func[any](c, func() {

			if name.Get() == "" {
				error.Set("please enter a name")
				return
			}

			if err := stream.Update(map[string]any{"name": name.Get()}); err != nil {
				error.Set(Fmt("cannot set name: %v", err))
				return
			}

			if err := stream.Save(); err != nil {
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

		// edit the name of the stream
		editStreamName := func(c Context) Element {
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
		case "configs":
			content = c.Element("streamConfigs", StreamConfigs(stream, onUpdate))
		case "sources":
			content = Div("coming soon...")
		}

		mainContent := func(c Context) Element {
			return Div(
				H2(
					Class("bulma-title"),
					router.Match(
						c,
						If(onUpdate != nil,
							Route("/name/edit",
								c.ElementFunction("editName", editStreamName),
							),
						),
						Route("",
							F(
								"Stream: ",
								stream.Name(),
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
						Fmt("last modified: %s", HumanDuration(time.Now().Sub(stream.CreatedAt()))),
					),
				),
				Div(Class("bulma-content"), IfElse(stream.Description() != "", stream.Description(), "(no description given)")),
				ui.Tabs(
					ui.Tab(ui.ActiveTab(tab == "configs"), A(Href(Fmt("/projects/%s/streams/details/%s/configs", Hex(project.ID()), streamId)), "Configs")),
					ui.Tab(ui.ActiveTab(tab == "sources"), A(Href(Fmt("/projects/%s/streams/details/%s/sources", Hex(project.ID()), streamId)), "Sources")),
				),
				content,
			)
		}

		return router.Match(
			c,
			If(tab == "configs", Route("/details/(?P<configId>[^/]+)(?:/(?P<tab>actions|settings))?", StreamConfigDetails(stream, onUpdate))),
			Route("", mainContent),
		)
	}
}

func NewStream(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {
	return func(c Context) Element {

		name := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)

		onSubmit := Func[any](c, func() {

			if name.Get() == "" {
				error.Set("Please enter a name")
				return
			}

			stream := project.MakeStream(nil)

			stream.SetName(name.Get())

			if err := stream.Save(); err != nil {
				error.Set("Cannot save stream")
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
						Placeholder("stream name"),
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
						"Create Stream",
					),
				),
			),
		)
	}
}

func StreamsList(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		// we retrieve the stream configs of the project...
		streams, err := project.Controller().Streams(map[string]interface{}{
			"project.id": project.ID(),
		})

		if err != nil {
			// to do: error handling
			return nil
		}

		ais := make([]Element, 0, len(streams))

		for _, stream := range streams {
			streamItem := A(
				Href(Fmt("/projects/%s/streams/details/%s", Hex(project.ID()), Hex(stream.ID()))),
				ui.ListItem(
					ui.ListColumn("md", stream.Name()),
					ui.ListColumn("sm", HumanDuration(time.Now().Sub(stream.CreatedAt()))),
				),
			)
			ais = append(ais, streamItem)
		}

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				If(onUpdate != nil, Route("/new", c.Element("newStream", NewStream(project, onUpdate)))),
				Route("", F(

					IfElse(
						len(ais) > 0,
						ui.List(
							ui.ListHeader(
								ui.ListColumn("md", "Name"),
								ui.ListColumn("sm", "Created At"),
							),
							ais,
						),
						ui.Message(
							"info",
							"No existing streams.",
						),
					),
					If(onUpdate != nil,
						A(
							Href(router.CurrentRoute().Path+"/new"),
							Class("bulma-button", "bulma-is-success"),
							"New Stream",
						),
					),
				),
				),
			),
		)
	}
}
