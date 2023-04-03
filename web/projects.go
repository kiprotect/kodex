package web

import (
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	ctrlHelpers "github.com/kiprotect/kodex/api/helpers/controller"
	"github.com/kiprotect/kodex/web/ui"
)

func InMemoryController(c Context) (api.Controller, error) {
	controller := UseController(c)
	return ctrlHelpers.InMemoryController(controller.Settings(), map[string]interface{}{}, controller.APIDefinitions())
}

// Project details

func ProjectDetails(c Context, projectId string, tab string) Element {

	controller := UseController(c)
	user := UseExternalUser(c)

	// we fetch the project
	project, err := controller.Project(Unhex(projectId))

	if err != nil {
		// to do: return error
		Log.Error("%v", err)
		return nil
	}

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

	title := GlobalVar[string](c, "title", "")

	title.Set(Fmt("Projects > %s", project.Name()))

	Log.Info("New Title: %s", title.Get())

	var content Element

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
		ui.List(pis),
		Button(Class("bulma-button", "bulma-is-success"), "New Project"),
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
