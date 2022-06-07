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

func TestDeleteStream(t *testing.T) {

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

		pt.FC{af.Organization{Name: "test"}, "org"},

		// we create two stream roles
		pt.FC{
			af.ObjectRole{
				ObjectName:       "projectA",
				ObjectRole:       "superuser",
				OrganizationRole: "project:a:superuser",
				Organization:     "org",
			},
			"projectRoleA",
		},
		pt.FC{
			pf.Config{
				Stream: "streamA",
				Name:   "test",
				Status: "active",
			},
			"config",
		},

		pt.FC{
			pf.ActionConfig{
				Name:    "test 1",
				Type:    "pseudonymize",
				Project: "projectA",
				Config: map[string]interface{}{
					"key":    "foo",
					"method": "merengue",
					"config": map[string]interface{}{},
				},
			},
			"action",
		},

		pt.FC{
			pf.ActionMap{
				Config: "config",
				Action: "action",
				Index:  0,
			},
			"actionMap",
		},

		pt.FC{
			af.ObjectRole{
				ObjectName:       "projectB",
				OrganizationRole: "project:b:admin",
				ObjectRole:       "admin",
				Organization:     "org",
			},
			"projectRoleB",
		},

		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:a:superuser", "project:b:admin"}, Scopes: []string{"kiprotect:api:stream:write"}}, "user"},
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

	streamA := fixtures["streamA"].(kodex.Stream)
	streamB := fixtures["streamB"].(kodex.Stream)

	reader := bytes.NewReader(nil)
	req, _ := http.NewRequest("DELETE", "/v1/streams/"+
		hex.EncodeToString(streamA.ID()), reader)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	// we try to delete a stream for which we don't have the permission
	reader = bytes.NewReader(nil)
	req, _ = http.NewRequest("DELETE", "/v1/streams/"+
		hex.EncodeToString(streamB.ID()), reader)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// this user doesn't have a matching role to delete the stream
	if resp.Code != 404 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	// we make sure the stream is deleted
	_, err = controller.Stream(streamA.ID())

	if err == nil {
		t.Fatalf("stream A is not deleted")
	}

	// we make sure the other stream isn't
	_, err = controller.Stream(streamB.ID())

	if err != nil {
		t.Fatalf("stream B is deleted")
	}

}
