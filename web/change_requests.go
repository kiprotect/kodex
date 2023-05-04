package web

import (
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/web/ui"
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

func ChangeRequestDetails(project kodex.Project) func(c Context, changeRequestId, tab string) Element {
	return func(c Context, changeRequestId, tab string) Element {

		if tab == "" {
			tab = "description"
		}

		router := UseRouter(c)
		controller := UseController(c)

		// we retrieve the action configs of the project...
		changeRequest, err := controller.ChangeRequest(Unhex(changeRequestId))

		if err != nil {
			// to do: error handling
			return nil
		}

		changeRequestIdVar := PersistentGlobalVar(c, "changeRequestId", "")

		onSubmit := Func(c, func() {
			changeRequestIdVar.Set(Hex(changeRequest.ID()))
			router.RedirectTo(router.CurrentPath())
		})

		var content Element

		switch tab {
		case "discussion":
		case "description":
			content = Div(Class("bulma-content"), IfElse(changeRequest.Description() != "", changeRequest.Description(), "(no description given)"))
		case "changes":
		}

		return Div(
			H2(Class("bulma-subtitle"), changeRequest.Title()),
			ui.Tabs(
				ui.Tab(ui.ActiveTab(tab == "description"), A(Href(Fmt("/projects/%s/changes/details/%s/description", Hex(project.ID()), changeRequestId)), "Description")),
				ui.Tab(ui.ActiveTab(tab == "discussion"), A(Href(Fmt("/projects/%s/changes/details/%s/discussion", Hex(project.ID()), changeRequestId)), "Discussion")),
				ui.Tab(ui.ActiveTab(tab == "changes"), A(Href(Fmt("/projects/%s/changes/details/%s/changes", Hex(project.ID()), changeRequestId)), "Changes")),
			),
			content,
			If(
				changeRequestIdVar.Get() != changeRequestId,
				F(
					Hr(),
					Form(
						Method("POST"),
						OnSubmit(onSubmit),
						Div(
							Class("bulma-field"),
							P(
								Class("bulma-control"),
								Button(
									Class("bulma-button", "bulma-is-success"),
									Type("submit"),
									"Work on this change request",
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

		onSubmit := Func(c, func() {

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
					ui.ListColumn("icon", "*"),
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
