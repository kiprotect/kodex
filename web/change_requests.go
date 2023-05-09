package web

import (
	"encoding/json"
	. "github.com/kiprotect/gospel"
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
				Route("/details/(?P<changeRequestId>[^/]+)(?:/(?P<tab>discussion|changes|description))?", ChangeRequestDetails(project)),
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
				ui.ListColumn("sm", strings.Join(pathStrs, " &gt; ")),
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

func MergeRequestNotice(c Context, project kodex.Project, changeRequest api.ChangeRequest) Element {

	router := UseRouter(c)

	onSubmit := Func[any](c, func() {

	})

	return ui.MessageWithTitle(
		"info",
		Div(
			Class("kip-col", "kip-is-lg"),
			"Do you really want to merge this change request?",
		),
		Div(
			P(
				"Merging the request will apply all changes to the current project. This cannot be undone!",
			),
			Form(
				Method("POST"),
				OnSubmit(onSubmit),
				Div(
					Class("bulma-field", "bulma-is-grouped"),
					P(
						Class("bulma-control"),
						A(
							Class("bulma-button"),
							Href(router.CurrentPath()),
							"Cancel",
						),
					),
					P(
						Class("bulma-control"),
						Button(
							Class("bulma-button", "bulma-is-success"),
							"Merge",
						),
					),
				),
			),
		),
	)
}

func ChangeRequestDetails(project kodex.Project) func(c Context, changeRequestId, tab string) Element {
	return func(c Context, changeRequestId, tab string) Element {

		if tab == "" {
			tab = "description"
		}

		router := UseRouter(c)
		controller := UseController(c)
		showMergeRequest := Var(c, false)

		// we retrieve the action configs of the project...
		changeRequest, err := controller.ChangeRequest(Unhex(changeRequestId))

		if err != nil {
			// to do: error handling
			return nil
		}

		changeRequestIdVar := PersistentGlobalVar(c, "changeRequestId", "")

		onSubmit := Func[any](c, func() {

			req := c.Request()

			action := req.FormValue("action")

			switch action {
			case "edit":
				changeRequestIdVar.Set(Hex(changeRequest.ID()))
				router.RedirectTo(router.CurrentPath())
			case "merge":
				showMergeRequest.Set(true)
			}
		})

		var content Element

		switch tab {
		case "discussion":
			content = Discussion(changeRequest)
		case "description":
			content = Div(Class("bulma-content"), IfElse(changeRequest.Description() != "", changeRequest.Description(), "(no description given)"))
		case "changes":
			content = Changes(changeRequest)
		}

		canMerge := true

		return Div(
			H2(Class("bulma-subtitle"), changeRequest.Title()),
			IfElse(
				showMergeRequest.Get(),
				MergeRequestNotice(c, project, changeRequest),
				F(
					ui.Tabs(
						ui.Tab(ui.ActiveTab(tab == "description"), A(Href(Fmt("/projects/%s/changes/details/%s/description", Hex(project.ID()), changeRequestId)), "Description")),
						ui.Tab(ui.ActiveTab(tab == "discussion"), A(Href(Fmt("/projects/%s/changes/details/%s/discussion", Hex(project.ID()), changeRequestId)), "Discussion")),
						ui.Tab(ui.ActiveTab(tab == "changes"), A(Href(Fmt("/projects/%s/changes/details/%s/changes", Hex(project.ID()), changeRequestId)), "Changes")),
					),
					content,
					Hr(),
					Form(
						Class("bulma-is-inline"),
						Method("POST"),
						OnSubmit(onSubmit),
						Div(
							Class("bulma-buttons", "bulma-has-addons"),
							If(
								changeRequestIdVar.Get() != changeRequestId,
								Button(
									Name("action"),
									Value("edit"),
									Class("bulma-button", "bulma-is-success"),
									Type("submit"),
									"Open",
								),
							),
							If(
								canMerge,
								Button(
									Name("action"),
									Value("merge"),
									Class("bulma-button", "bulma-is-primary"),
									Type("submit"),
									"Merge",
								),
							),
						),
					),
				),
			),
		)
	}
}

func NewChangeRequest(project kodex.Project) ElementFunction {
	return func(c Context) Element {

		title := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)
		user := UseApiUser(c)

		controller := UseController(c)

		onSubmit := Func[any](c, func() {

			if title.Get() == "" {
				error.Set("Please enter a title")
				return
			}

			changeRequest, err := controller.MakeChangeRequest(project, user)

			if err != nil {
				error.Set(Fmt("Cannot create change request: %v", err))
				return
			}

			changeRequest.SetTitle(title.Get())
			changeRequest.SetStatus(api.DraftCR)

			if err := changeRequest.Save(); err != nil {
				error.Set(Fmt("Cannot save change request: %v", err))
			} else {
				router.RedirectTo(Fmt("/projects/%s/changes/details/%s", Hex(project.ID()), Hex(changeRequest.ID())))
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
					Class("bulma-label", "Title"),
					Input(
						Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
						Value(title),
						Placeholder("change request title"),
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
						"Create Change Request",
					),
				),
			),
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
			return nil
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
			router.Match(
				c,
				Route("/new", c.Element("newChangeRequest", NewChangeRequest(project))),
				Route("",
					F(
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
					),
				),
			),
		)

		return nil
	}
}
