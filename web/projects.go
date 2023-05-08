package web

import (
	"bytes"
	"encoding/json"
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	ctrlHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	"github.com/kiprotect/kodex/web/ui"
	"io"
	"strings"
	"time"
)

func InMemoryController(c Context) (api.Controller, error) {
	controller := UseController(c)
	return ctrlHelpers.InMemoryController(controller.Settings(), map[string]interface{}{}, controller.APIDefinitions())
}

func NewProject() ElementFunction {
	return func(c Context) Element {

		name := Var(c, "")
		error := Var(c, "")
		router := UseRouter(c)
		controller := UseController(c)

		onSubmit := Func[any](c, func() {

			if name.Get() == "" {
				error.Set("Please enter a name")
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

			project := controller.MakeProject(nil)

			project.SetName(name.Get())

			if err := project.Save(); err != nil {
				error.Set("Cannot save project")
				return
			}

			org := UseDefaultOrganization(c)

			if org == nil {
				error.Set("Cannot get organization")
				return
			}

			apiOrg, err := org.ApiOrganization(controller)

			if err != nil {
				error.Set("Cannot get organization")
				return
			}

			// we always add admin and superuser roles
			for _, orgRole := range []string{"admin", "superuser"} {
				role := controller.MakeObjectRole(project, apiOrg)
				values := map[string]interface{}{
					"organization_role": orgRole,
					"role":              "superuser",
				}

				if err := role.Create(values); err != nil {
					error.Set(Fmt("Cannot create role: %v", err))
					return
				}
				if err := role.Save(); err != nil {
					error.Set(Fmt("Cannot save role: %v", err))
					return
				}
			}

			// we try to add default roles as well
			if defaultRoles, err := controller.DefaultObjectRoles(apiOrg.ID()); err != nil {
				kodex.Log.Errorf("Cannot load default roles: %v", err)
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
						return
					}
					if err := role.Save(); err != nil {
						return
					}

				}
			}

			success = true

			router.RedirectTo(Fmt("/projects/%s", Hex(project.ID())))
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
				),
			),
		)
	}
}

func Settings(project kodex.Project, onUpdate func(ChangeInfo, string)) ElementFunction {

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

			err = blueprint.CreateWithProject(project.Controller(), project)

			if err != nil {
				error.Set(Fmt("Error creating blueprint: %v", err))
				return
			}

			actionConfigs, _ := project.Controller().ActionConfigs(map[string]any{})

			error.Set(Fmt("Success: %d", len(actionConfigs)))

			onUpdate(ChangeInfo{}, router.CurrentPath())
		})

		return F(
			ui.MessageWithTitle(
				"info",
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
											"screenshot",
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
						console.log("hey");
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

	// we load the project
	projectVar := CachedVar(c, func() kodex.Project {

		Log.Error("Loading project....")

		project, err := controller.Project(Unhex(projectId))

		if err != nil {
			error.Set(Fmt("Cannot load project: %v", err))
			// to do: return error
			Log.Error("%v", err)
			return nil
		}

		return project

	})

	// we retrieve the project...
	project := projectVar.Get()

	if project == nil {
		// to do: error handling...
		return nil
	}

	objectRolesVar := CachedVar(c, func() []api.ObjectRole {
		roles, err := controller.RolesForObject(project)

		if err != nil {
			error.Set(Fmt("Cannot load object roles: %v", err))
			Log.Error("%v", err)
			return nil
		}
		return roles

	})

	objectRoles := objectRolesVar.Get()

	if objectRoles == nil {
		// to do: error handling...
		return nil
	}

	AddBreadcrumb(c, "Projects", "/projects")
	AddBreadcrumb(c, project.Name(), Fmt("/%s", Hex(project.ID())))

	// we check that the user can access the project
	if ok, err := controller.CanAccess(user, project, []string{"read", "write", "admin"}); !ok || err != nil {
		Log.Error("cannot access")
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

	changeRequestVar := CachedVar(c, func() api.ChangeRequest {

		if changeRequestId.Get() == "" {
			return nil
		}

		// we retrieve the action configs of the project...
		changeRequest, err := controller.ChangeRequest(Unhex(changeRequestId.Get()))

		if err != nil {
			error.Set(Fmt("cannot load change request: %v", err))
			changeRequestId.Set("")
			return nil
		}

		if !bytes.Equal(changeRequest.ObjectID(), project.ID()) {
			changeRequestId.Set("")
			return nil
		}

		return changeRequest

	})

	changeRequest := changeRequestVar.Get()

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

		if changeRequest != nil {

			exportedBlueprint, err = kodex.ExportBlueprint(project)

			if err != nil {
				Log.Error("Error: %v", err)
				return nil
			}

			importedBlueprint = kodex.MakeBlueprint(exportedBlueprint)
			importedProject, err = importedBlueprint.Create(ctrl, true)

			if err != nil {
				return Div(Fmt("uh oh: %v (%s)", err, changeRequest.Changes()))
			}

		} else {
			Log.Error("Import error: %v", err)
			return nil
		}

	}

	AddBreadcrumb(c, strings.Title(tab), Fmt("/%s", tab))

	onUpdate := func(change ChangeInfo, path string) {

		changeRequest := changeRequestVar.Get()

		// we persist the project changes (if there were any)
		Log.Error("Updating blueprint...")

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
			error.Set(Fmt("cannot apply changes: %v", err))
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

	if changeRequestVar.Get() == nil {
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

	switch tab {
	case "actions":
		content = c.Element("actions", Actions(importedProject, onUpdate))
	case "changes":
		content = c.Element("changes", ChangeRequests(importedProject))
	case "settings":
		content = c.Element("settings", Settings(importedProject, onUpdate))
	default:
		content = Div("...")
	}

	onDoneEditing := Func[any](c, func() {
		changeRequestId.Set("")
		router := UseRouter(c)
		router.RedirectTo(router.CurrentPath())
	})

	return Div(
		Div(
			Class("bulma-content"),
			H2(Class("bulma-title"), project.Name()),
		),
		Div(
			Class("bulma-tags"),
			Span(
				Class("bulma-tag", "bulma-is-info", "bulma-is-light"),
				Fmt("last modified: %s", HumanDuration(time.Now().Sub(project.CreatedAt()))),
			),
		),

		If(
			onUpdate == nil,
			ui.Message("warning",
				F(
					I(
						Class("fa", "fa-lock"),
					),
					" Read-only mode, please open a change request to edit project.",
				),
			),
		),
		DoIf(
			changeRequest != nil,
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
									Fmt("/projects/%s/changes/details/%s",
										projectId,
										Hex(changeRequest.ID()),
									),
								),
								changeRequest.Title(),
							),
							Fmt(", %d changes so far.", len(changeRequest.Changes())),
							Form(
								Method("POST"),
								OnSubmit(onDoneEditing),
								Div(
									Class("bulma-field"),
									P(
										Class("bulma-control"),
										Button(
											Class("bulma-button", "bulma-is-info"),
											Type("submit"),
											"Finish work",
										),
									),
								),
							),
						),
					),
				)
			},
		),
		If(error.Get() != "", ui.Message("danger", error.Get())),
		ui.Tabs(
			ui.Tab(ui.ActiveTab(tab == "actions"), A(Href(Fmt("/projects/%s/actions", projectId)), "Actions")),
			ui.Tab(ui.ActiveTab(tab == "changes"), A(Href(Fmt("/projects/%s/changes", projectId)), "Change Requests")),
			ui.Tab(ui.ActiveTab(tab == "settings"), A(Href(Fmt("/projects/%s/settings", projectId)), "Settings")),
		),
		content,
		Hr(),
		roles,
	)
}

func Projects(c Context) Element {

	externalUser := UseExternalUser(c)
	controller := UseController(c)

	projects, err := projects(controller, externalUser)

	if err != nil {
		// to do: redirect to error page...
		kodex.Log.Error(err)
		return nil
	}

	AddBreadcrumb(c, "Projects", "/projects")

	pis := make([]any, 0, len(projects))

	for _, project := range projects {

		projectItem := A(
			Href(Fmt("/projects/%s", Hex(project.ID()))),
			ui.ListItem(
				ui.ListColumn("md", project.Name()),
				ui.ListColumn("sm", HumanDuration(time.Now().Sub(project.CreatedAt()))),
			),
		)
		pis = append(pis, projectItem)
	}

	return F(
		ui.List(
			ui.ListHeader(
				ui.ListColumn("md", "Name"),
				ui.ListColumn("sm", "Created At"),
			),
			pis,
		),
		A(Href("/projects/new"), Class("bulma-button", "bulma-is-success"), "New Project"),
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
