package web

import (
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/web/ui"
)

func ChangeRequests(project kodex.Project) ElementFunction {

	return func(c Context) Element {

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				Route("/details/(?P<changeRequestId>[^/]+)", ChangeRequestDetails(project)),
				Route("", ChangeRequestList(project)),
			),
		)
	}
}

func ChangeRequestDetails(project kodex.Project) func(c Context, changeRequestId string) Element {
	return func(c Context, changeRequestId string) Element {

		router := UseRouter(c)
		controller := UseController(c)

		// we retrieve the action configs of the project...
		changeRequest, err := controller.ChangeRequest(Unhex(changeRequestId))

		if err != nil {
			// to do: error handling
			return nil
		}

		onSubmit := Func(c, func() {

			changeRequestId := PersistentGlobalVar(c, "changeRequestId", "")
			changeRequestId.Set(Hex(changeRequest.ID()))
			router.RedirectUp()
		})

		return Div(
			changeRequest.Title(),
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
							"Work On Change Request",
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
			changeRequest.SetStatus(api.Draft)

			if err := changeRequest.Save(); err != nil {
				error.Set(Fmt("Cannot save change request: %v", err))
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
				),
			)
			cri = append(cri, changeRequestItem)
		}

		router := UseRouter(c)

		return F(
			router.Match(
				c,
				Route("/new", c.Element("newChangeRequest", NewChangeRequest(project))),
				Route("", F(
					ui.List(cri),
					A(
						Href(router.CurrentRoute().Path+"/new"),
						Class("bulma-button", "bulma-is-success"),
						"New Change Request"),
				),
				),
			),
		)

		return nil
	}
}
