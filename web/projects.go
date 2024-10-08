// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	ctrlHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	"github.com/kiprotect/kodex/web/ui"
	"io"
	"net/http"
	"strings"
	"time"
)

func InMemoryController(c Context) (api.Controller, error) {
	controller := UseController(c)
	return ctrlHelpers.InMemoryController(controller.Settings(), map[string]interface{}{}, controller.APIDefinitions())
}

func MakeProject(controller api.Controller, name string, org *api.UserOrganization) (kodex.Project, error) {

	project := controller.MakeProject(nil)

	project.SetName(name)

	if err := project.Save(); err != nil {
		return nil, fmt.Errorf("Cannot save project: %v", err)
	}

	apiOrg, err := org.ApiOrganization(controller)

	if err != nil {
		return nil, fmt.Errorf("cannot get API organization: %v", err)
	}

	// we always add admin and superuser roles
	for _, orgRole := range []string{"admin", "superuser"} {
		role := controller.MakeObjectRole(project, apiOrg)
		values := map[string]interface{}{
			"organization_role": orgRole,
			"role":              "superuser",
		}

		if err := role.Create(values); err != nil {
			return nil, fmt.Errorf("cannot create role: %v", err)
		}
		if err := role.Save(); err != nil {
			return nil, fmt.Errorf("cannot save role: %v", err)
		}
	}

	// we try to add default roles as well
	if defaultRoles, err := controller.DefaultObjectRoles(apiOrg.ID()); err != nil {
		return nil, fmt.Errorf("cannot load default roles: %v", err)
	} else {
		for _, defaultRole := range defaultRoles {
			if defaultRole.ObjectType() != "project" {
				continue
			}

			role := controller.MakeObjectRole(project, apiOrg)

			values := map[string]interface{}{
				"organization_role": defaultRole.OrganizationRole(),
				"role":              defaultRole.ObjectRole(),
			}

			if err := role.Create(values); err != nil {
				return nil, fmt.Errorf("cannot create role: %v", err)
			}
			if err := role.Save(); err != nil {
				return nil, fmt.Errorf("cannot save role: %v", err)
			}

		}
	}

	return project, nil

}

func CanCreate(user *api.ExternalUser, objectType string) bool {

	for _, userRoles := range user.Roles {
		for _, orgRole := range userRoles.Roles {
			if orgRole == "admin" || orgRole == "superuser" || orgRole == "editor" {
				return true
			}
		}
	}

	return false
}

func NewProject() ElementFunction {
	return func(c Context) Element {

		form := MakeFormData(c, "newProject", POST)
		name := form.Var("name", "")
		error := Var(c, "")
		router := UseRouter(c)
		controller := UseController(c)
		user := UseExternalUser(c)

		if !CanCreate(user, "project") {
			return Div("cannot create projects")
		}

		onSubmit := func() {

			if name.Get() == "" {
				error.Set("Please enter a name")
				return
			}

			org := UseDefaultOrganization(c)

			if org == nil {
				error.Set("Cannot get organization")
				return
			}

			controller.Begin()

			success := false

			defer func() {
				if success {
					controller.Commit()
				}
				controller.Rollback()
			}()

			project, err := MakeProject(controller, name.Get(), org)

			if err != nil {
				error.Set(Fmt("%v", err))
				return
			}

			success = true

			router.RedirectTo(Fmt("/flows/projects/%s", Hex(project.ID())))
		}

		form.OnSubmit(onSubmit)

		var errorNotice Element

		if error.Get() != "" {
			errorNotice = P(
				Class("bulma-help", "bulma-is-danger"),
				error.Get(),
			)
		}

		return form.Form(
			H1(Class("bulma-subtitle"), "New Project"),
			Div(
				Class("bulma-field"),
				errorNotice,
				Label(
					Class("bulma-label", "Name"),
					Input(
						Class("bulma-input", If(error.Get() != "", "bulma-is-danger")),
						Type("text"),
						Value(name),
						Placeholder("project name"),
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
						"Create Project",
					),
					Nbsp,
					A(
						Class("bulma-button"),
						Href("/flows/projects"),
						"Cancel",
					),
				),
			),
		)
	}
}

func Settings(project, realProject kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		router := UseRouter(c)

		return router.Match(
			c,
			Route("/roles", ProjectRolesRoutes(realProject)),
			Route("/export-blueprint",
				c.ElementFunction("downloadBlueprint",
					func(c Context) Element {

						changedBlueprint, err := kodex.ExportBlueprint(project)

						if err != nil {
							return Div("Cannot download blueprint")
						}

						bytes, err := json.MarshalIndent(changedBlueprint, "", "  ")

						if err != nil {
							return Div("Cannot download blueprint")
						}

						c.SetRespondWith(func(c Context, w http.ResponseWriter) {
							w.Header().Add("content-type", "application/json")
							w.Header().Add("content-disposition", Fmt("attachment; filename=Kodex Blueprint - %s.json;", project.Name()))
							w.Write(bytes)
						})

						return nil
					}),
			),
			Route("/delete",
				c.ElementFunction("deleteProject",
					func(c Context) Element {

						router := UseRouter(c)

						onSubmit := Func[any](c, func() {

							realProject.Delete()
							router.RedirectTo("/flows/projects")

						})

						return ui.MessageWithTitle(
							"danger",
							"Confirm Project Deletion",
							F(
								Div(
									Class("kip-col", "kip-is-lg"),
									"Do you really want to delete this project? This cannot be undone!",
								),
								Div(
									Class("kip-col", "kip-is-icon"),
									Form(
										Method("POST"),
										OnSubmit(onSubmit),
										Div(
											Class("bulma-field", "bulma-is-grouped"),
											P(
												Class("bulma-control"),
												A(
													Class("bulma-button"),
													Href(router.LastPath()),
													"Cancel",
												),
											),
											P(
												Class("bulma-control"),
												Button(
													Class("bulma-button", "bulma-is-danger"),
													"Delete",
												),
											),
										),
									),
								),
							),
						)
					}),
			),
			Route(
				"",
				c.ElementFunction("settings", SettingsTab(project, onUpdate)),
			),
		)
	}

}

func SettingsTab(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {

	return func(c Context) Element {

		error := Var(c, "")
		router := UseRouter(c)

		onSubmit := Func[any](c, func() {

			request := c.Request()

			file, _, err := request.FormFile("blueprint")

			if err != nil {
				error.Set(Fmt("Cannot retrieve file: %v", err))
				return
			}

			content, err := io.ReadAll(file)

			if err != nil {
				error.Set(Fmt("Cannot read file: %v", err))
				return
			}

			error.Set(Fmt("file length: %d", len(content)))

			var data map[string]any

			if err := json.Unmarshal(content, &data); err != nil {
				error.Set(Fmt("cannot unmarshal JSON: %v", err))
				return
			}

			data["project"] = map[string]any{
				"id":   Hex(project.ID()),
				"name": project.Name(),
			}

			blueprint := kodex.MakeBlueprint(data)
			// we do not preserve IDs here as that will lead to conflicts
			blueprint.PreserveIDs = false

			err = blueprint.CreateWithProject(project.Controller(), project)

			if err != nil {
				error.Set(Fmt("Error creating blueprint: %v", err))
				return
			}

			onUpdate(ChangeInfo{}, router.CurrentPath())
		})

		return Div(
			ui.MessageWithTitle(
				"grey",
				"Import Blueprint",
				F(
					P(
						"You can import a blueprint from a JSON file.",
					),
					Br(),
					If(
						error.Get() != "",
						P(
							Class("bulma-help", "bulma-is-danger"),
							error.Get(),
						),
					),
					IfElse(
						onUpdate != nil,
						F(
							Form(
								Method("POST"),
								Enctype("multipart/form-data"),
								OnSubmit(onSubmit),
								Div(
									Id("blueprint-file"),
									Class("bulma-file", "bulma-has-name"),
									Label(
										Class("bulma-file-label"),
										Input(
											Class("bulma-file-input"),
											Type("file"),
											Id("blueprint"),
											Name("blueprint"),
										),
										Span(
											Class("bulma-file-cta"),
											Span(
												Class("bulma-file-icon"),
												I(
													Class("fas", "fa-upload"),
												),
											),
											Span(
												Class("bulma-file-label"),
												"Info file...",
											),
										),
										Span(
											Class("bulma-file-name"),
											"please select a file",
										),
									),
								),
								Hr(),
								Button(
									Class("bulma-button", "bulma-is-success"),
									Type("submit"),
									"Import Blueprint",
								),
							),
						),
						P(
							"You need to open a change request to import first.",
						),
					),
					Script(`
						const fileInput = document.querySelector('#blueprint-file input[type=file]');
						  fileInput.onchange = () => {
						    if (fileInput.files.length > 0) {
						      const fileName = document.querySelector('#blueprint-file .bulma-file-name');
						      fileName.textContent = fileInput.files[0].name;
						    }
						  }
					`),
				),
			),
			ui.MessageWithTitle(
				"grey",
				"Export Blueprint",
				F(
					P(
						"You can export a blueprint to a JSON file.",
					),
					Br(),
					A(
						Href(Fmt("/flows/projects/%s/settings/export-blueprint", Hex(project.ID()))),
						Class("bulma-button", "bulma-is-success"),
						DataAttrib("plain", ""),
						Type("submit"),
						"Export Blueprint",
					),
				),
			),
			ui.MessageWithTitle(
				"danger",
				"Delete Project",
				F(
					P(
						"You can delete the project. This cannot be undone!",
					),
					Br(),
					A(
						Href(Fmt("/flows/projects/%s/settings/delete", Hex(project.ID()))),
						Class("bulma-button", "bulma-is-danger"),
						Type("submit"),
						"Delete Project",
					),
				),
			),
		)
	}
}

type ChangeInfo struct {
	Description string
	Data        any
}

func ProjectDetails(c Context, projectId string, tab string) Element {

	error := Var(c, "")

	controller := UseController(c)
	user := UseExternalUser(c)
	router := UseRouter(c)

	project, err := controller.Project(Unhex(projectId))

	if err != nil {
		return Div(Fmt("Cannot load project: %v", err))
	}

	objectRoles, err := controller.RolesForObject(project)

	if err != nil {
		return Div(Fmt("Cannot load object roles: %v", err))
	}

	AddBreadcrumb(c, "Projects", "/flows/projects")
	AddBreadcrumb(c, project.Name(), Fmt("/%s", Hex(project.ID())))

	// we check that the user can access the project
	if ok, err := controller.CanAccess(user, project, []string{}); !ok || err != nil {
		return nil
	}

	exportedBlueprint, err := kodex.ExportBlueprint(project)

	if err != nil {
		Log.Error("Error: %v", err)
		return nil
	}

	ctrl, err := InMemoryController(c)

	if err != nil {
		Log.Error("Error: %v", err)
		return nil
	}

	var content Element

	if tab == "" {
		tab = "actions"
	}

	changeRequestId := PersistentGlobalVar(c, "changeRequestId", "")

	var changeRequest api.ChangeRequest

	if changeRequestId.Get() != "" {

		// we retrieve the change request...
		changeRequest, err = controller.ChangeRequest(Unhex(changeRequestId.Get()))

		if err != nil {
			error.Set(Fmt("cannot load change request: %v", err))
			changeRequestId.Set("")
		} else if !bytes.Equal(changeRequest.ObjectID(), project.ID()) {
			changeRequestId.Set("")
		}
	}

	if changeRequest != nil {

		if changeRequest.Changes() != nil {
			for _, changeSet := range changeRequest.Changes() {
				if err := api.ApplyChanges(exportedBlueprint, changeSet.Changes); err != nil {
					error.Set(Fmt("cannot apply changes: %v", err))
				}
			}

		}
	}

	importedBlueprint := kodex.MakeBlueprint(exportedBlueprint)
	importedProject, err := importedBlueprint.Create(ctrl, true)

	if err != nil {

		importErr := err

		// we reset the controller
		ctrl, err = InMemoryController(c)

		if err != nil {
			error.Set(Fmt("cannot recreate in-memory controller"))
			return nil
		}

		if changeRequest != nil {

			// we reexport the original project again
			exportedBlueprint, err = kodex.ExportBlueprint(project)

			if err != nil {
				Log.Error("Error: %v", err)
				return nil
			}

			reimportedBlueprint := kodex.MakeBlueprint(exportedBlueprint)
			importedProject, err = reimportedBlueprint.Create(ctrl, true)

			// to do: fix this import problem with the in-memory controller
			if err != nil {
				return Div(Fmt("cannot import project: %v (%s)", err, changeRequest.Changes()))
			}

		} else {
			return Div(Fmt("Import error: %v", importErr))
		}

	}

	AddBreadcrumb(c, strings.Title(tab), Fmt("/%s", tab))

	onUpdate := func(change ChangeInfo, path string) {

		changedBlueprint, err := kodex.ExportBlueprint(importedProject)

		if err != nil {
			error.Set(Fmt("cannot export changes: %v", err))
			return
		}

		changes := api.DiffWithOptions(exportedBlueprint, changedBlueprint, api.DiffOptions{
			Identifiers: []string{"id", "name"},
		})

		changeSets := changeRequest.Changes()

		changeSet := api.ChangeSet{
			Description: change.Description,
			Data:        change.Data,
			Changes:     changes,
		}

		changeSets = append(changeSets, changeSet)

		if err := changeRequest.SetChanges(changeSets); err != nil {
			error.Set(Fmt("cannot set changes: %v", err))
			return
		}

		if err := api.ApplyChanges(exportedBlueprint, changeSet.Changes); err != nil {
			error.Set(Fmt("cannot apply changes: %v", err))
			return
		}

		// we create a test controller
		testCtrl, err := InMemoryController(c)

		if err != nil {
			error.Set(Fmt("cannot recreate in-memory controller"))
			return
		}

		reimportedBlueprint := kodex.MakeBlueprint(exportedBlueprint)
		_, err = reimportedBlueprint.Create(testCtrl, true)

		// we cannot recreate the modified blueprint, we abort
		if err != nil {
			error.Set(Fmt("cannot import project: %v (%s)", err, changeRequest.Changes()))
			return
		}

		if err := changeRequest.Save(); err != nil {
			error.Set(Fmt("cannot save change request: %v", err))
			return
		}

		// we redirect to the requested path...
		router := UseRouter(c)
		router.RedirectTo(path)

	}

	if changeRequest == nil {
		onUpdate = nil
	}

	userRoles := []Element{}
	foundRoles := map[string]any{}

	for _, role := range user.Roles {
		for _, orgRole := range role.Roles {
			for _, objRole := range objectRoles {
				if objRole.OrganizationRole() == orgRole {

					if _, ok := foundRoles[objRole.ObjectRole()]; ok {
						continue
					}

					foundRoles[objRole.ObjectRole()] = true

					userRoles = append(userRoles, Span(Class("bulma-tag", "bulma-is-dark"), objRole.ObjectRole()))
				}
			}
		}
	}

	roles := F(
		H3("Your object roles"),
		Div(
			Class("bulma-tags"),
			userRoles,
		),
	)

	mainContent := func(c Context) Element {

		switch tab {
		case "streams":
			content = c.Element("streams", Streams(importedProject, onUpdate))
		case "actions":
			content = c.Element("actions", Actions(importedProject, onUpdate))
		case "changes":
			content = c.Element("changes", ChangeRequests(project))
		case "settings":
			content = c.Element("settings", Settings(importedProject, project, onUpdate))
		default:
			content = Div("unknown section")
		}

		return Div(
			Div(
				Class("bulma-content"),
			),
			//			Div(
			//				Class("bulma-tags"),
			//				Span(
			//					Class("bulma-tag", "bulma-is-info", "bulma-is-light"),
			//					Fmt("last modified: %s", HumanDuration(time.Now().Sub(project.CreatedAt()))),
			//				),
			//			),
			content,
			Hr(),
			roles,
		)
	}

	projectMenu := []*SidebarItem{
		{
			Title:  project.Name(),
			Path:   Fmt("/flows/projects/%s", projectId),
			Header: true,
		},
		{
			Title: "Actions",
			Path:  Fmt("/flows/projects/%s/actions", projectId),
			Icon:  "play-circle",
		},
		{
			Title: "Streams",
			Path:  Fmt("/flows/projects/%s/streams", projectId),
			Icon:  "random",
		},
		{
			Title: "Changes",
			Path:  Fmt("/flows/projects/%s/changes", projectId),
			Icon:  "folder-open",
		},
		{
			Title: "Settings",
			Path:  Fmt("/flows/projects/%s/settings", projectId),
			Icon:  "cogs",
			Submenu: []*SidebarItem{
				{
					Title: "Roles",
					Path:  Fmt("/flows/projects/%s/settings/roles", projectId),
					Icon:  "users",
				},
			},
		},
	}

	projectsMenu := GetSidebarItemByPath(c, "/flows/projects")

	if projectsMenu == nil {
		Log.Warning("Cannot find 'projects' sidebar menu...")
	} else {
		projectsMenu.Submenu = append(projectsMenu.Submenu, projectMenu...)
	}

	canEdit, err := controller.CanAccess(user, project, []string{"editor"})

	if err != nil {
		Log.Warning("Cannot get rights: %v", err)
	}

	return F(
		If(error.Get() != "", ui.Message("danger", error.Get())),
		//ui.Modal(c, router.CurrentPath()),
		If(
			onUpdate == nil && (tab == "actions" || tab == "streams") && canEdit,
			ui.Message("warning",
				F(
					I(
						Class("fa", "fa-lock"),
					),
					" Read-only mode, please click on 'edit' to make changes.",
					Div(
						Class("bulma-is-pulled-right"),
						A(
							Style("height: 32px; margin-top: -5px;"),
							Class("bulma-button", "bulma-is-success"),
							Href(Fmt("/flows/projects/%s/changes/new", projectId)),
							"edit",
						),
					),
				),
			),
		),
		DoIf(
			changeRequest != nil && (tab == "actions" || tab == "streams"),
			func() Element {
				return F(
					ui.Message("info",
						F(
							I(
								Class("fa", "fa-check"),
							),
							" Working on change request ",
							A(
								Href(
									Fmt("/flows/projects/%s/changes/details/%s",
										projectId,
										Hex(changeRequest.ID()),
									),
								),
								changeRequest.Title(),
							),
							Div(
								Class("bulma-is-pulled-right"),
								A(
									Style("height: 32px; margin-top: -5px;"),
									Class("bulma-button", "bulma-is-success"),
									Href(Fmt("/flows/projects/%s/changes/details/%s/close", projectId, Hex(changeRequest.ID()))),
									"I'm done",
								),
							),
						),
					),
				)
			},
		),

		router.Match(
			c,
			If(tab == "streams", Route("/details/(?P<streamId>[^/]+)(?:/(?P<tab>configs|sources|settings))?", StreamDetails(importedProject, onUpdate))),
			If(tab == "actions", Route("/details/(?P<actionId>[^/]+)(?:/(?P<tab>edit|test|data))?", ActionDetails(importedProject, onUpdate))),
			Route("", mainContent),
		),
	)
}

func Projects(c Context) Element {

	user := UseExternalUser(c)
	controller := UseController(c)
	projects, err := projects(controller, user)

	if err != nil {
		// to do: redirect to error page...
		kodex.Log.Error(err)
		return nil
	}

	AddBreadcrumb(c, "Projects", "/flows/projects")

	pis := make([]any, 0, len(projects))

	for _, project := range projects {

		projectItem := A(
			Href(Fmt("/flows/projects/%s", Hex(project.ID()))),
			ui.ListItem(
				ui.ListColumn("md", project.Name()),
				ui.ListColumn("sm", HumanDuration(time.Now().Sub(project.UpdatedAt()))),
			),
		)
		pis = append(pis, projectItem)
	}

	return F(
		ui.List(
			ui.ListHeader(
				ui.ListColumn("md", "Name"),
				ui.ListColumn("sm", "Updated At"),
			),
			pis,
		),
		If(
			CanCreate(user, "project"),
			A(Href("/flows/projects/new"), Class("bulma-button", "bulma-is-success"), "New Project"),
		),
	)
}

// Helper Functions

// Projects list

func projects(controller api.Controller, user *api.ExternalUser) ([]kodex.Project, error) {

	objectRoles, err := controller.ObjectRolesForUser("project", user)

	if err != nil {
		return nil, err
	}

	ids := make([]interface{}, len(objectRoles))

	for i, role := range objectRoles {
		ids[i] = role.ObjectID()
	}

	projects, err := controller.Projects(map[string]interface{}{
		"id": kodex.In{Values: ids},
	})

	if err != nil {
		return nil, err
	}

	return projects, nil
}
