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
)

// Tests updating of streams by a stream superuser.
func TestUpdateStream(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		// we create a test stream
		pt.FC{pf.Stream{Name: "test 1", Project: "project"}, "stream"},

		pt.FC{af.Organization{Name: "Org"}, "org"},
		// we create a stream role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "org"}, "projectRole"},

		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:stream:write"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(*api.ExternalUser)
	stream := fixtures["stream"].(kodex.Stream)
	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	sourceData, _ := json.Marshal(map[string]interface{}{
		"name":        "another test",
		"description": "another test",
	})

	reader := bytes.NewReader(sourceData)

	req, _ := http.NewRequest("PATCH", "/v1/streams/"+hex.EncodeToString(stream.ID()), reader)
	req.Header.Set("content-type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	if err := stream.Refresh(); err != nil {
		t.Fatal(err)
	}

	if stream.Name() != "another test" || stream.Description() != "another test" {
		t.Fatal("update did not work")
	}

	// we test only updating the name

	sourceData, _ = json.Marshal(map[string]interface{}{
		"name": "third test",
	})

	reader = bytes.NewReader(sourceData)

	req, _ = http.NewRequest("PATCH", "/v1/streams/"+hex.EncodeToString(stream.ID()), reader)
	req.Header.Set("content-type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	if err := stream.Refresh(); err != nil {
		t.Fatal(err)
	}

	if stream.Name() != "third test" || stream.Description() != "another test" {
		t.Fatal("update did not work")
	}

	// we test only updating the description

	sourceData, _ = json.Marshal(map[string]interface{}{
		"description": "third test",
	})

	reader = bytes.NewReader(sourceData)

	req, _ = http.NewRequest("PATCH", "/v1/streams/"+hex.EncodeToString(stream.ID()), reader)
	req.Header.Set("content-type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	if err := stream.Refresh(); err != nil {
		t.Fatal(err)
	}

	if stream.Name() != "third test" || stream.Description() != "third test" {
		t.Fatal("update did not work")
	}

	// we test using invalid data

	sourceData, _ = json.Marshal(map[string]interface{}{
		"name": "1",
	})

	reader = bytes.NewReader(sourceData)

	req, _ = http.NewRequest("PATCH", "/v1/streams/"+hex.EncodeToString(stream.ID()), reader)
	req.Header.Set("content-type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 400 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}

// Tests that only stream superusers can update a stream.
func TestUpdateStreamAuthorization(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},
		// we create a test stream
		pt.FC{pf.Stream{Name: "test 1", Project: "project"}, "stream"},

		pt.FC{af.Organization{Name: "Org"}, "org"},

		// we create a project role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:master", ObjectRole: "superuser", Organization: "org"}, "projectRole"},

		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:stream:write"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(*api.ExternalUser)
	stream := fixtures["stream"].(kodex.Stream)
	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)

	if err != nil {
		t.Fatal(err)
	}

	sourceData, _ := json.Marshal(map[string]interface{}{
		"name":        "another test",
		"description": "another test",
	})

	reader := bytes.NewReader(sourceData)

	req, _ := http.NewRequest("PATCH", "/v1/streams/"+hex.EncodeToString(stream.ID()), reader)
	req.Header.Set("content-type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 404 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}
