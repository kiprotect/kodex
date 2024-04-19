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

package resources_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
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

// Tests the creation of a project role by a project superuser.
func TestCreateRole(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		pt.FC{af.Organization{Name: "A"}, "org"},

		// we create a project role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "org"}, "projectRole"},

		pt.FC{af.User{Email: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:project:roles"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(*api.ExternalUser)
	org := fixtures["org"].(api.Organization)
	project := fixtures["project"].(kodex.Project)
	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	sourceData, _ := json.Marshal(map[string]interface{}{
		"role":              "superuser",
		"organization_role": "projectmaster",
	})

	reader := bytes.NewReader(sourceData)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/orgs/%s/projects/%s/roles", hex.EncodeToString(org.SourceID()), hex.EncodeToString(project.ID())), reader)
	req.Header.Set("content-type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	roles, err := controller.RolesForObject(project)

	if err != nil {
		t.Fatal(err)
	}

	if len(roles) != 2 {
		t.Fatal("wrong number of roles")
	}

	found := false

	for _, role := range roles {
		if role.OrganizationRole() == "projectmaster" && role.ObjectRole() == "superuser" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("role not found")
	}

	// we test invalid data handling

	for _, data := range []map[string]interface{}{
		map[string]interface{}{
			"role":              "superuses",
			"organization_role": "projectmaster",
		},
		map[string]interface{}{
			"role":              "",
			"organization_role": "projectmaster",
		},
		map[string]interface{}{
			"role":              "superuser",
			"organization_role": "",
		},
		map[string]interface{}{
			"role":              "superuser",
			"organization_role": "ölapalömablanka",
		},
	} {
		sourceData, _ = json.Marshal(data)

		reader = bytes.NewReader(sourceData)

		req, _ = http.NewRequest("POST", fmt.Sprintf("/v1/orgs/%s/projects/%s/roles", hex.EncodeToString(org.SourceID()), hex.EncodeToString(project.ID())), reader)
		req.Header.Set("content-type", "application/json")
		resp = httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != 400 {
			t.Fatalf("wrong return code: %d", resp.Code)
		}

	}

}

// Tests that a user who is not a project superuser cannot create project roles.
func TestRoleCreateAuthorization(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		pt.FC{af.Organization{Name: "Test"}, "org"},

		// we create a project role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:master", ObjectRole: "superuser", Organization: "org"}, "projectRole"},

		pt.FC{af.User{Email: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:project:roles"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	controller := fixtures["controller"].(api.Controller)
	user := fixtures["user"].(*api.ExternalUser)
	org := fixtures["org"].(api.Organization)
	project := fixtures["project"].(kodex.Project)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	sourceData, _ := json.Marshal(map[string]interface{}{
		"role":              "superuser",
		"organization_role": "projectmaster",
	})

	reader := bytes.NewReader(sourceData)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/orgs/%s/projects/%s/roles", hex.EncodeToString(org.SourceID()), hex.EncodeToString(project.ID())), reader)
	req.Header.Set("content-type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 404 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}
