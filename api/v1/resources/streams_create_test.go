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
	"strings"
	"testing"
)

func uuidToBytes(uuidStr string) ([]byte, error) {
	rawUUIDStr := strings.Replace(uuidStr, "-", "", -1)
	return hex.DecodeString(rawUUIDStr)
}

var createStreamFixtures = []pt.FC{

	// we create the settings
	pt.FC{pf.Settings{}, "settings"},

	// we create an SQL controller
	pt.FC{af.Controller{}, "controller"},
	pt.FC{pf.Project{Name: "test"}, "project"},
	pt.FC{af.Organization{Name: "test"}, "org"},

	pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "admin", Organization: "org"}, "projectRole"},

	pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:stream:create"}}, "user"},

	pt.FC{af.User{EMail: "mux@mustermann.de", Organization: "org", Roles: []string{"admins"}, Scopes: []string{"kiprotect:api:stream:create"}}, "userB"},
}

func TestCreateStream(t *testing.T) {

	fixtures, err := pt.SetupFixtures(createStreamFixtures)
	defer pt.TeardownFixtures(createStreamFixtures, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	controller := fixtures["controller"].(api.Controller)
	user := fixtures["user"].(*api.User)
	project := fixtures["project"].(kodex.Project)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)
	if err != nil {
		t.Fatal(err)
	}

	sourceData, _ := json.Marshal(map[string]interface{}{
		"name":        "test",
		"description": "test",
	})

	reader := bytes.NewReader(sourceData)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/projects/%s/streams", hex.EncodeToString(project.ID())), reader)
	req.Header.Set("content-type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	var value map[string]interface{}

	if err = json.Unmarshal(resp.Body.Bytes(), &value); err != nil {
		t.Fatal("Invalid JSON")
	}

	var streamMap map[string]interface{}

	if stream, ok := value["data"]; !ok {
		t.Fatal("no stream in response")
	} else {
		streamMap, ok = stream.(map[string]interface{})
		if !ok {
			t.Fatal("not a map")
		}
		if streamMap["id"] == nil {
			t.Fatal("no ID in stream")
		}
	}

	streamID, ok := streamMap["id"].(string)

	if !ok {
		t.Fatal("ID is not a string")
	}

	streamBinaryID, err := uuidToBytes(streamID)

	if err != nil {
		t.Fatal(err)
	}

	_, err = controller.Stream(streamBinaryID)

	if err != nil {
		t.Fatal(err)
	}
}

// Tests that a user who is no superuser/admin cannot create a stream
func TestInvalidCreateStream(t *testing.T) {

	fixtures, err := pt.SetupFixtures(createStreamFixtures)
	defer pt.TeardownFixtures(createStreamFixtures, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["userB"].(*api.User)
	controller := fixtures["controller"].(api.Controller)
	project := fixtures["project"].(kodex.Project)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)
	if err != nil {
		t.Fatal(err)
	}

	sourceData, _ := json.Marshal(map[string]interface{}{
		"name":        "test",
		"description": "test",
	})

	reader := bytes.NewReader(sourceData)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/projects/%s/streams", hex.EncodeToString(project.ID())), reader)
	req.Header.Set("content-type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 404 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}
