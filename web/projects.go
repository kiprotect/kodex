package web

import (
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	ctrlHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	"github.com/kiprotect/kodex/web/ui"
	"strings"
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

// Project details

func ProjectDetails(c Context, projectId string, tab string) Element {

	controller := UseController(c)
	user := UseExternalUser(c)

	project, err := controller.Project(Unhex(projectId))

	if err != nil {
		// to do: return error
		Log.Error("%v", err)
		return nil
	}

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

	importedBlueprint := kodex.MakeBlueprint(exportedBlueprint)

	importedProject, err := importedBlueprint.Create(ctrl, true)

	if err != nil {
		Log.Error("Import error: %v", err)
		return nil
	}

	var content Element

	if tab == "" {
		tab = "actions"
	}

	AddBreadcrumb(c, strings.Title(tab), Fmt("/%s", tab))

	switch tab {
	case "actions":
		content = c.Element("actions", Actions(importedProject))
	case "changes":
		content = c.Element("changes", ChangeRequests(importedProject))
	default:
		content = Div("...")
	}

	router := UseRouter(c)
	req := router.Request()

	// we persist the project changes (if there were any)
	if req.Method == "POST" && c.Interactive() {
		Log.Error("Updating blueprint...")

		exportedBlueprint, err = kodex.ExportBlueprint(importedProject)

		if err != nil {
			Log.Error("cannot export blueprint: %v", err)
			return nil
		}

		Log.Info("%v", exportedBlueprint)

		importedBlueprint = kodex.MakeBlueprint(exportedBlueprint)

		// we store the blueprint again
		if _, err := importedBlueprint.Create(controller, false); err != nil {
			Log.Error("Error saving blueprint: %v", err)

			return nil
		}

	}

	Log.Info("%v", exportedBlueprint)

	return Div(
		Div(
			Class("bulma-content"),
			H2(Class("bulma-title"), project.Name()),
		),
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
