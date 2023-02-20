// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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
	"time"
)

// Tests that a user can retrieve project roles.
func TestListRoles(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		pt.FC{af.Organization{Name: "A"}, "orgA"},

		// we create a project role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "orgA"}, "projectRoleA"},

		// and another one
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:margarita", ObjectRole: "superuser", Organization: "orgA"}, "projectRoleB"},

		pt.FC{af.User{Email: "max@mustermann.de", Organization: "orgA", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:project:roles"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(*api.ExternalUser)
	project := fixtures["project"].(kodex.Project)
	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	originalRoles, err := controller.RolesForObject(project)

	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(nil)

	req, _ := http.NewRequest("GET", "/v1/projects/"+hex.EncodeToString(project.ID())+"/roles", reader)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	var values map[string][]map[string]interface{}

	if err = json.Unmarshal(resp.Body.Bytes(), &values); err != nil {
		t.Error("Invalid JSON")
	}

	roles, ok := values["data"]

	if !ok {
		t.Fatalf("roles missing")
	}

	for _, role := range roles {
		if !ok {
			t.Fatalf("expected a string map")
		}
		binaryID, err := hex.DecodeString(role["object_id"].(string))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(binaryID, project.ID()) {
			t.Fatalf("IDs do not match")
		}
		found := false
		for _, originalRole := range originalRoles {
			if originalRole.OrganizationRole() == role["organization_role"].(string) &&
				originalRole.ObjectRole() == role["object_role"].(string) &&
				originalRole.CreatedAt().Format(time.RFC3339Nano) == role["created_at"].(string) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("didn't find a matching role")
		}
	}

}

// Tests that a user who is not a project superuser cannot list project roles.
func TestRoleListAuthorization(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		pt.FC{af.Organization{Name: "A"}, "orgA"},

		// we create a project role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:master", ObjectRole: "superuser", Organization: "orgA"}, "projectRole"},

		pt.FC{af.User{Email: "max@mustermann.de", Organization: "orgA", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:project:roles"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(*api.ExternalUser)
	project := fixtures["project"].(kodex.Project)
	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(nil)

	req, _ := http.NewRequest("GET", "/v1/projects/"+hex.EncodeToString(project.ID())+"/roles", reader)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 404 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}
