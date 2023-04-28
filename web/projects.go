package web

import (
	"bytes"
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	ctrlHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	"github.com/kiprotect/kodex/web/ui"
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

		onSubmit := Func(c, func() {

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

	msg := PersistentVar(c, []api.Change{})

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
			// to do: error handling
			return nil
		}

		if !bytes.Equal(changeRequest.ObjectID(), project.ID()) {
			error.Set(Fmt("change request not valid for this project"))
			return nil
		}

		return changeRequest

	})

	changeRequest := changeRequestVar.Get()

	if changeRequest != nil {

		if changeRequest.Changes() != nil {
			if err := api.ApplyChanges(exportedBlueprint, changeRequest.Changes()); err != nil {
				error.Set(Fmt("cannot apply changes: %v", err))
			}
		}
	}

	importedBlueprint := kodex.MakeBlueprint(exportedBlueprint)
	importedProject, err := importedBlueprint.Create(ctrl, true)

	if err != nil {
		Log.Error("Import error: %v", err)
		return nil
	}

	AddBreadcrumb(c, strings.Title(tab), Fmt("/%s", tab))

	onUpdate := func(change api.Change, path string) {

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

		msg.Set(changes)

		existingChanges := changeRequest.Changes()

		if existingChanges != nil {
			changes = append(existingChanges, changes...)
		}

		if err := changeRequest.SetChanges(changes); err != nil {
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

	switch tab {
	case "actions":
		content = c.Element("actions", Actions(importedProject, onUpdate))
	case "changes":
		content = c.Element("changes", ChangeRequests(importedProject))
	default:
		content = Div("...")
	}

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
		Fmt("Change: %v", msg.Get()),
		Fmt("Change request: %s", changeRequestId.Get()),
		If(error.Get() != "", ui.Message("danger", error.Get())),
		ui.Tabs(
			ui.Tab(ui.ActiveTab(tab == "actions"), A(Href(Fmt("/projects/%s/actions", projectId)), "Actions")),
			ui.Tab(ui.ActiveTab(tab == "changes"), A(Href(Fmt("/projects/%s/changes", projectId)), "Changes")),
			ui.Tab(ui.ActiveTab(tab == "settings"), A(Href(Fmt("/projects/%s/settings", projectId)), "Settings")),
		),
		content,
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
		kodex.Log.Infof("Name: %s", project.Name())
		projectItem := A(
			Href(Fmt("/projects/%s", Hex(project.ID()))),
			ui.ListItem(
				ui.ListColumn("md", project.Name()),
			),
		)
		pis = append(pis, projectItem)
	}

	return F(
		ui.List(
			ui.ListHeader(
				ui.ListColumn("md", "Name"),
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
