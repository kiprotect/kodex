package web

import (
	. "github.com/kiprotect/gospel"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/web/ui"
)

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

	title := GetVar[string](c, "title")

	SetVar(c, "title", "foobar")

	Log.Info("New Title: %s", title.Get())

	Log.Info(tab)

	return Div(
		Div(
			Class("bulma-content"),
			H2(Class("bulma-title"), project.Name()),
		),
		ui.Tabs(
			ui.Tab(ui.ActiveTab(tab == "configuration"), A(Href(Fmt("/projects/%s/configuration", projectId)), "Configuration")),
			ui.Tab(ui.ActiveTab(tab == "changes"), A(Href(Fmt("/projects/%s/changes", projectId)), "Changes")),
			ui.Tab(ui.ActiveTab(tab == "settings"), A(Href(Fmt("/projects/%s/settings", projectId)), "Settings")),
		),
		Div(
			"[in progress]",
		),
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
