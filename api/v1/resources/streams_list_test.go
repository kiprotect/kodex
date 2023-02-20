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

// Tests the retrieval of the streams list.
func TestStreams(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test A"}, "projectA"},
		pt.FC{pf.Project{Name: "test B"}, "projectB"},

		// we create two test streams
		pt.FC{pf.Stream{Name: "test 1", Project: "projectA"}, "streamA"},
		pt.FC{pf.Stream{Name: "test 2", Project: "projectB"}, "streamB"},

		pt.FC{af.Organization{Name: "A"}, "orgA"},
		pt.FC{af.Organization{Name: "B"}, "orgB"},

		// we create two project roles
		pt.FC{af.ObjectRole{ObjectName: "projectA", OrganizationRole: "project:a:admin", ObjectRole: "admin", Organization: "orgA"}, "projectRoleA"},
		pt.FC{af.ObjectRole{ObjectName: "projectB", OrganizationRole: "project:b:admin", ObjectRole: "admin", Organization: "orgB"}, "projectRoleB"},

		pt.FC{af.User{Email: "max@mustermann.de", Organization: "orgA", Roles: []string{"project:a:admin"}, Scopes: []string{"kiprotect:api:stream:read", "kiprotect:api:project:read"}}, "userA"},
		pt.FC{af.User{Email: "max@mustermann.de", Organization: "orgB", Roles: []string{"project:b:admin"}, Scopes: []string{"kiprotect:api:stream:read", "kiprotect:api:project:read"}}, "userB"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	controller := fixtures["controller"].(api.Controller)
	project := fixtures["projectA"].(kodex.Project)

	withUser := func(c *gin.Context) {
		user := fixtures["userA"]
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/projects/%s/streams", hex.EncodeToString(project.ID())), nil)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	var values map[string]interface{}

	if err = json.Unmarshal(resp.Body.Bytes(), &values); err != nil {
		t.Error("Invalid JSON")
	}

	streams, ok := values["data"]

	if !ok {
		t.Fatalf("streams missing")
	}

	streamsList, ok := streams.([]interface{})

	if !ok {
		t.Fatalf("streams not a list")
	}

	if len(streamsList) != 1 {
		t.Fatalf("expected 1 item")
	}

}
