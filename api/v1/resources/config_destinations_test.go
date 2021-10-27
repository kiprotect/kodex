// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

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

// Tests the retrieval of the streams list.
func TestAddConfigDestination(t *testing.T) {

	var fixturesConfig = []pt.FC{
		pt.FC{pf.Settings{}, "settings"},
		pt.FC{af.Controller{}, "controller"},
		pt.FC{pf.Project{Name: "test"}, "project"},
		pt.FC{pf.Stream{Name: "test 1", Project: "project"}, "stream"},
		pt.FC{af.Organization{Name: "Test"}, "org"},
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "org"}, "projectRole"},
		pt.FC{pf.Config{Stream: "stream", Name: "test 1", Version: "1", Source: "api", Status: "active"}, "config"},
		pt.FC{pf.Destination{Name: "stream", Project: "project", DestinationType: "in-memory", Config: map[string]interface{}{}}, "destination"},
		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:destination:write", "kiprotect:api:config:write"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	controller := fixtures["controller"].(api.Controller)

	withUser := func(c *gin.Context) {
		user := fixtures["user"]
		c.Set("userProfile", user)
	}

	router, err := at.Router(controller, withUser)
	if err != nil {
		t.Fatal(err)
	}

	config := fixtures["config"].(kodex.Config)
	destination := fixtures["destination"].(kodex.Destination)

	for _, destinationStatus := range []string{"active", "disabled", "testing"} {

		sourceData, _ := json.Marshal(map[string]interface{}{
			"status": destinationStatus,
			"name":   "test",
		})

		reader := bytes.NewReader(sourceData)

		req, _ := http.NewRequest("POST", "/v1/configs/"+hex.EncodeToString(config.ID())+"/destinations/"+hex.EncodeToString(destination.ID()), reader)
		req.Header.Set("content-type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != 200 {
			t.Fatalf("wrong return code: %d", resp.Code)
		}

		if err := config.Refresh(); err != nil {
			t.Fatal(err)
		}

		destinations, err := config.Destinations()

		if err != nil {
			t.Fatal(err)
		}

		if len(destinations) != 1 {
			t.Fatalf("expected 1 destination, got %d", len(destinations))
		}

		var destinationName string
		var destinationMaps []kodex.DestinationMap
		for destinationName, destinationMaps = range destinations {
			break
		}

		if len(destinationMaps) != 1 {
			t.Fatalf("Expected 1 destination map, got %d", len(destinationMaps))
		}

		destinationMap := destinationMaps[0]

		if !bytes.Equal(destinationMap.Destination().ID(), destination.ID()) ||
			destinationMap.Status() != kodex.DestinationStatus(destinationStatus) ||
			destinationName != "test" {
			t.Fatalf("IDs or statuses do not match")
		}

	}
}

// Tests the retrieval of the streams list.
func TestRemoveConfigDestination(t *testing.T) {

	var fixturesConfig = []pt.FC{
		pt.FC{pf.Settings{}, "settings"},
		pt.FC{af.Controller{}, "controller"},
		pt.FC{pf.Project{Name: "test"}, "project"},
		pt.FC{pf.Stream{Name: "test 1", Project: "project"}, "stream"},
		pt.FC{pf.Config{Stream: "stream", Name: "test 1", Version: "1", Source: "api", Status: "active"}, "config"},
		pt.FC{pf.Destination{Name: "stream", Project: "project", DestinationType: "in-memory", Config: map[string]interface{}{}}, "destination"},
		pt.FC{pf.DestinationAdder{Destination: "destination", Name: "stream", Config: "config", Status: "active"}, "destinationMap"},
		pt.FC{af.Organization{Name: "test"}, "org"},
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "org"}, "projectRole"},
		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:destination:write", "kiprotect:api:config:write"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	withUser := func(c *gin.Context) {
		user := fixtures["user"]
		c.Set("userProfile", user)
	}

	controller := fixtures["controller"].(api.Controller)

	router, err := at.Router(controller, withUser)
	if err != nil {
		t.Fatal(err)
	}

	config := fixtures["config"].(kodex.Config)
	destination := fixtures["destination"].(kodex.Destination)

	req, _ := http.NewRequest("DELETE", "/v1/configs/"+hex.EncodeToString(config.ID())+"/destinations/"+hex.EncodeToString(destination.ID()), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	if err := config.Refresh(); err != nil {
		t.Fatal(err)
	}

	destinations, err := config.Destinations()

	if err != nil {
		t.Fatal(err)
	}

	if len(destinations) != 0 {
		t.Fatalf("expected 0 destinations")
	}

}
