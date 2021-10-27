// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources_test

import (
	"bytes"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Tests the deletion of a project role by a project superuser.
func TestDeleteRole(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		pt.FC{af.Organization{Name: "A"}, "orgA"},

		// we create a project role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "orgA"}, "projectRoleA"},
		// and another one that we will delete
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:master", ObjectRole: "superuser", Organization: "orgA"}, "projectRoleB"},

		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "orgA", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:project:roles"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(api.UserProfile)
	project := fixtures["project"].(kodex.Project)
	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		c.Set("userProfile", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	roles, err := controller.RolesForObject(project)

	if err != nil {
		t.Fatal(err)
	}

	var role api.ObjectRole

	for _, r := range roles {
		if r.OrganizationRole() == "project:master" {
			role = r
			break
		}
	}

	if role == nil {
		t.Fatalf("role not found")
	}

	reader := bytes.NewReader(nil)

	req, _ := http.NewRequest("DELETE", "/v1/projects/"+hex.EncodeToString(project.ID())+"/roles/"+hex.EncodeToString(role.ID()), reader)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	roles, err = controller.RolesForObject(project)

	if err != nil {
		t.Fatal(err)
	}

	if len(roles) != 1 {
		t.Fatalf("role not deleted")
	}

	role = roles[0]

	if role.OrganizationRole() != "project:admin" || role.ObjectRole() != "superuser" {
		t.Fatalf("wrong role deleted")
	}

	// we test that the user can't delete the last role that allows him/her
	// superuser access to the project

	reader = bytes.NewReader(nil)

	req, _ = http.NewRequest("DELETE", "/v1/projects/"+hex.EncodeToString(project.ID())+"/roles/"+hex.EncodeToString(role.ID()), reader)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 400 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}

// Tests that a user who is not a project superuser cannot delete project roles.
func TestRoleDeleteAuthorization(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		pt.FC{af.Organization{Name: "A"}, "orgA"},

		// we create a project role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:master", ObjectRole: "superuser", Organization: "orgA"}, "projectRole"},

		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "orgA", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:project:roles"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(api.UserProfile)
	project := fixtures["project"].(kodex.Project)
	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		c.Set("userProfile", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	roles, err := controller.RolesForObject(project)

	if err != nil {
		t.Fatal(err)
	}

	if len(roles) == 0 {
		t.Fatalf("expected at least one role")
	}

	role := roles[0]

	reader := bytes.NewReader(nil)

	req, _ := http.NewRequest("DELETE", "/v1/projects/"+hex.EncodeToString(project.ID())+"/roles/"+hex.EncodeToString(role.ID()), reader)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 404 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}