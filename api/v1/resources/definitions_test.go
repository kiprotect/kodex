// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources_test

import (
	"encoding/json"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
)

// Tests the retrieval of definitions
func TestDefinitions(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},
		pt.FC{af.Organization{Name: "org"}, "org"},
		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{}, Scopes: []string{"kiprotect:api:definitions"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(api.UserProfile)
	controller := fixtures["controller"].(api.Controller)

	resp, err := at.Get(controller, user, "/v1/definitions", nil)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	var values map[string]interface{}

	if err = json.Unmarshal(resp.Body.Bytes(), &values); err != nil {
		t.Error("Invalid JSON")
	}

	_, ok := values["data"]

	if !ok {
		t.Fatalf("definitions missing")
	}

}
