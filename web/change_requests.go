package web

import (
	"bytes"
	"encoding/json"
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/web/ui"
	"strings"
	"time"
)

func ChangeRequests(project kodex.Project) ElementFunction {

	return func(c Context) Element {

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				Route("/details/(?P<changeRequestId>[^/]+)(?:/(?P<tab>discussion|changes|description|merge|close))?", ChangeRequestDetails(project)),
				Route("", ChangeRequestList(project)),
			),
		)
	}
}

func Discussion(changeRequest api.ChangeRequest) Element {
	return Div("coming soon...")
}

func Changes(changeRequest api.ChangeRequest) Element {

	cri := make([]Element, 0)

	for a, changeSet := range changeRequest.Changes() {

		for b, change := range changeSet.Changes {

			pathStrs := make([]string, len(change.Path))

			for i, pathElement := range change.Path {
				pathStrs[i] = pathElement.String()
			}

			value, _ := json.MarshalIndent(change.Value, "", "  ")

			changeItem := ui.ListItem(
				ui.ListColumn("xs", Fmt("%d.%d", a+1, b+1)),
				ui.ListColumn("sm", Strong(api.OpName(change.Op))),
				ui.ListColumn("sm", strings.Join(pathStrs, " > ")),
				ui.ListColumn("md", Pre(string(value))),
			)

			cri = append(cri, changeItem)

		}
	}

	return ui.List(
		ui.ListHeader(
			ui.ListColumn("xs", "#No"),
			ui.ListColumn("sm", "Operator"),
			ui.ListColumn("sm", "Path"),
			ui.ListColumn("md", "Change Value (JSON)"),
		),
		cri,
	)

}

func CloseChangeRequest(c Context, project kodex.Project, changeRequest api.ChangeRequest) Element {

	router := UseRouter(c)
	error := Var(c, "")

	onSubmit := Func[any](c, func() {

		/*
			// we export the current blueprint
			exportedBlueprint, err := kodex.ExportBlueprint(project)

			if err != nil {
				error.Set(Fmt("Cannot export blueprint: %v", err))
				return
			}

			for _, changeSet := range changeRequest.Changes() {
				if err := api.ApplyChanges(exportedBlueprint, changeSet.Changes); err != nil {
					error.Set(Fmt("cannot apply changes: %v", err))
					return
				}
			}

			importedBlueprint := kodex.MakeBlueprint(exportedBlueprint)

			// we re-import the blueprint
			_, err = importedBlueprint.Create(project.Controller(), false)

			if err != nil {
				error.Set(Fmt("Cannot import project: %v", err))
				return
			}

			if err := changeRequest.Delete(); err != nil {
				error.Set(Fmt("Cannot delete change request: %v", err))
				return
			}
		*/

		// we remove the current change request
		changeRequestId := PersistentGlobalVar(c, "changeRequestId", "")
		changeRequestId.Set("")

		router.RedirectTo(Fmt("/projects/%s", Hex(project.ID())))
	})

	cancelButton := A(
		Class("bulma-button"),
		Href(Fmt("/projects/%s/changes/details/%s", Hex(project.ID()), Hex(changeRequest.ID()))),
		"Cancel",
	)

	return ui.FormModal(
		c,
		[]any{
			Method("POST"),
			OnSubmit(onSubmit),
		},
		F(
			IfElse(
				error.Get() != "",
				P(
					"There was an error merging this change request: ",
					Strong(error.Get()),
				),
				P(
					"If you're done editing this change request and want it to be reviewed, please choose ", Strong("Ready for review"),
					". If you just want to close it for now but keep editing it later, choose ", Strong("Close for now"), ".",
				),
			),
		),
		F(
			cancelButton,
			Span(
				Style("flex-grow: 1"),
			),
			F(
				Button(
					Class("bulma-button", "bulma-is-primary"),
					"Close for now",
				),
				Button(
					Class("bulma-button", "bulma-is-success"),
					"Ready for review",
				),
			),
		),
		"Close Change Request",
		Fmt("/projects/%s/changes", Hex(project.ID())),
	)
}

func ChangeRequestDetails(project kodex.Project) func(c Context, changeRequestId, tab string) Element {
	return func(c Context, changeRequestId, tab string) Element {

		if tab == "" {
			tab = "description"
		}

		controller := UseController(c)

		// we retrieve the action configs of the project...
		changeRequest, err := controller.ChangeRequest(Unhex(changeRequestId))

		if err != nil {
			// to do: error handling
			return nil
		}

		var content Element

		switch tab {
		case "discussion":
			content = Discussion(changeRequest)
		case "description":
			content = Div(Class("bulma-content"), IfElse(changeRequest.Description() != "", changeRequest.Description(), "(no description given)"))
		case "changes":
			content = Changes(changeRequest)
		}

		return Div(
			H2(Class("bulma-subtitle"), changeRequest.Title()),
			If(
				tab == "close",
				CloseChangeRequest(c, project, changeRequest),
			),
			F(
				ui.Tabs(
					ui.Tab(ui.ActiveTab(tab == "description"), A(Href(Fmt("/projects/%s/changes/details/%s/description", Hex(project.ID()), changeRequestId)), "Description")),
					ui.Tab(ui.ActiveTab(tab == "discussion"), A(Href(Fmt("/projects/%s/changes/details/%s/discussion", Hex(project.ID()), changeRequestId)), "Discussion")),
					ui.Tab(ui.ActiveTab(tab == "changes"), A(Href(Fmt("/projects/%s/changes/details/%s/changes", Hex(project.ID()), changeRequestId)), "Changes")),
				),
				content,
			),
		)
	}
}

func NewChangeRequest(project kodex.Project, confirmed bool) ElementFunction {
	return func(c Context) Element {

		title := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)
		user := UseApiUser(c)

		controller := UseController(c)

		if !confirmed {

			existingChangeRequests, err := controller.ChangeRequests(project)

			if err != nil {
				// to do: proper error view
				return Div("Cannot load change request")
			}

			changeRequestId := Var(c, "")

			onSubmit := Func[any](c, func() {
				Log.Info("Opening existing change request '%s'...", changeRequestId.Get())
				changeRequestIdVar := PersistentGlobalVar(c, "changeRequestId", "")
				// to do: verify ID
				changeRequestIdVar.Set(changeRequestId.Get())
				router.RedirectTo(router.CurrentPath())
			})

			for _, changeRequest := range existingChangeRequests {
				if changeRequest.Status() == api.DraftCR && bytes.Equal(changeRequest.Creator().SourceID(), user.SourceID()) {
					return ui.Modal(
						c,
						"Existing change request found",
						Span(
							"There already is a draft change request ",
							A(
								Href(
									Fmt("/projects/%s/changes/details/%s", Hex(project.ID()), Hex(changeRequest.ID())),
								),
								Strong(changeRequest.Title()),
							),
							", do you want to work on this one instead?",
						),
						F(
							A(
								Class("bulma-button"),
								Href(Fmt("/projects/%s", Hex(project.ID()))),
								"Cancel",
							),
							Span(Style("flex-grow: 1")),
							Span(
								Form(
									Class("bulma-is-inline"),
									Method("POST"),
									OnSubmit(onSubmit),
									Input(
										Type("hidden"),
										Value(changeRequestId, Hex(changeRequest.ID())),
									),
									Button(
										Name("action"),
										Value("edit"),
										Class("bulma-button", "bulma-is-warning"),
										Type("submit"),
										"Use Existing",
									),
								),

								A(
									Class("bulma-button", "bulma-is-success"),
									Href(Fmt("/projects/%s/changes/new/confirm", Hex(project.ID()))),
									"Open New One",
								),
							),
						),
						router.CurrentPath(),
					)
				}
			}
		}

		onSubmit := Func[any](c, func() {

			if title.Get() == "" {
				error.Set("Please enter a title")
				return
			}

			changeRequest, err := controller.MakeChangeRequest(nil, project, user)

			if err != nil {
				error.Set(Fmt("Cannot create change request: %v", err))
				return
			}

			changeRequest.SetTitle(title.Get())
			changeRequest.SetStatus(api.DraftCR)

			if err := changeRequest.Save(); err != nil {
				error.Set(Fmt("Cannot save change request: %v", err))
			} else {

				// we open the change request
				changeRequestIdVar := PersistentGlobalVar(c, "changeRequestId", "")
				changeRequestIdVar.Set(Hex(changeRequest.ID()))

				// we redirect to the main project view
				router.RedirectTo(Fmt("/projects/%s", Hex(project.ID())))
			}
		})

		var errorNotice Element

		if error.Get() != "" {
			errorNotice = P(
				Class("bulma-help", "bulma-is-danger"),
				error.Get(),
			)
		}

		return ui.FormModal(
			c,
			[]any{
				Method("POST"),
				OnSubmit(onSubmit),
			},
			F(
				Div(
					Class("bulma-field"),
					errorNotice,
					Label(
						Class("bulma-label", "Title"),
						Input(
							Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
							Value(title),
							Placeholder("change request title"),
						),
					),
				),
			),
			F(
				Div(
					Class("bulma-field"),
					P(
						Class("bulma-control"),
						Button(
							Class("bulma-button", "bulma-is-success"),
							Type("submit"),
							"Create Change Request",
						),
					),
				),
			),
			"New Change Request",
			Fmt("/projects/%s/changes", Hex(project.ID())),
		)
	}
}

func ChangeRequestList(project kodex.Project) ElementFunction {

	return func(c Context) Element {

		controller := UseController(c)

		// we retrieve the action configs of the project...
		changeRequests, err := controller.ChangeRequests(project)

		if err != nil {
			// to do: error handling
			return Div(Fmt("Cannot show change requests: %v", err))
		}

		cri := make([]Element, 0, len(changeRequests))

		for _, changeRequest := range changeRequests {
			changeRequestItem := A(
				Href(Fmt("/projects/%s/changes/details/%s", Hex(project.ID()), Hex(changeRequest.ID()))),
				ui.ListItem(
					ui.ListColumn("md", changeRequest.Title()),
					ui.ListColumn("sm", changeRequest.Creator().DisplayName()),
					ui.ListColumn("sm", HumanDuration(time.Now().Sub(changeRequest.CreatedAt()))),
					ui.ListColumn("icon", string(changeRequest.Status())),
				),
			)
			cri = append(cri, changeRequestItem)
		}

		router := UseRouter(c)

		return F(
			IfElse(
				len(cri) > 0,
				F(
					ui.List(
						ui.ListHeader(
							ui.ListColumn("md", "Name"),
							ui.ListColumn("sm", "Created By"),
							ui.ListColumn("sm", "Created At"),
							ui.ListColumn("icon", "Status"),
						),
						cri),
				),
				ui.Message(
					"info",
					"No open change requests.",
				),
			),
			A(
				Href(router.CurrentRoute().Path+"/new"),
				Class("bulma-button", "bulma-is-success"),
				"New Change Request",
			),
			router.Match(
				c,
				Route("/new/confirm", c.Element("newChangeRequest", NewChangeRequest(project, true))),
				Route("/new", c.Element("newChangeRequest", NewChangeRequest(project, false))),
			),
		)

		return nil
	}
}
