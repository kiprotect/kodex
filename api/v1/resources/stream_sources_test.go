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
	"testing"
)

// Tests adding an source to a stream
func TestAddStreamSource(t *testing.T) {

	var fixturesConfig = []pt.FC{
		pt.FC{pf.Settings{}, "settings"},
		pt.FC{af.Controller{}, "controller"},
		pt.FC{pf.Project{Name: "test"}, "project"},
		pt.FC{pf.Stream{Name: "test 1", Project: "project"}, "stream"},
		pt.FC{pf.Source{
			Name:       "source",
			SourceType: "bytes",
			Project:    "project",
			Config: map[string]interface{}{
				"input":  []byte{},
				"format": "json",
			},
		}, "source"},
		pt.FC{af.Organization{Name: "test"}, "org"},
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "org"}, "projectRole"},
		pt.FC{af.User{Email: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:stream:write", "kiprotect:api:source:write"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		user := fixtures["user"]
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)
	if err != nil {
		t.Fatal(err)
	}

	stream := fixtures["stream"].(kodex.Stream)
	source := fixtures["source"].(kodex.Source)

	fmt.Printf("Adding source %s to stream %s\n", hex.EncodeToString(source.ID()), hex.EncodeToString(stream.ID()))

	for _, sourceStatus := range []string{"active", "disabled", "testing"} {

		sourceData, _ := json.Marshal(map[string]interface{}{
			"status": sourceStatus,
		})

		reader := bytes.NewReader(sourceData)

		req, _ := http.NewRequest("POST", "/v1/streams/"+hex.EncodeToString(stream.ID())+"/sources/"+hex.EncodeToString(source.ID()), reader)
		req.Header.Set("content-type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != 200 {
			t.Fatalf("wrong return code: %d", resp.Code)
		}

		if err := stream.Refresh(); err != nil {
			t.Fatal(err)
		}

		sources, err := stream.Sources()

		if err != nil {
			t.Fatal(err)
		}

		if len(sources) != 1 {
			t.Fatalf("expected 1 source, got %d", len(sources))
		}

		var sourceMap kodex.SourceMap
		for _, v := range sources {
			sourceMap = v
			break
		}

		if !bytes.Equal(sourceMap.Source().ID(), source.ID()) || sourceMap.Status() != kodex.SourceStatus(sourceStatus) {
			t.Fatalf("IDs or statuses do not match")
		}

	}
}

// Tests the removal of a stream source
func TestRemoveStreamSource(t *testing.T) {

	var fixturesConfig = []pt.FC{
		pt.FC{pf.Settings{}, "settings"},
		pt.FC{af.Controller{}, "controller"},
		pt.FC{pf.Project{Name: "test"}, "project"},
		pt.FC{pf.Stream{Name: "test 1", Project: "project"}, "stream"},
		pt.FC{pf.Source{
			Name:       "source",
			SourceType: "bytes",
			Project:    "project",
			Config: map[string]interface{}{
				"input":  []byte{},
				"format": "json",
			},
		}, "source"},
		pt.FC{pf.SourceAdder{Source: "source", Stream: "stream", Status: "active"}, "sourceMap"},
		pt.FC{af.Organization{Name: "test"}, "org"},
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "org"}, "projectRole"},
		pt.FC{af.User{Email: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:stream:write", "kiprotect:api:source:write"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		user := fixtures["user"]
		c.Set("user", user)
	}

	router, err := at.Router(controller, withUser)
	if err != nil {
		t.Fatal(err)
	}

	stream := fixtures["stream"].(kodex.Stream)
	source := fixtures["source"].(kodex.Source)

	fmt.Printf("Adding source %s to stream %s\n", hex.EncodeToString(source.ID()), hex.EncodeToString(stream.ID()))

	req, _ := http.NewRequest("DELETE", "/v1/streams/"+hex.EncodeToString(stream.ID())+"/sources/"+hex.EncodeToString(source.ID()), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	if err := stream.Refresh(); err != nil {
		t.Fatal(err)
	}

	sources, err := stream.Sources()

	if err != nil {
		t.Fatal(err)
	}

	if len(sources) != 0 {
		t.Fatalf("expected 0 sources")
	}

}
