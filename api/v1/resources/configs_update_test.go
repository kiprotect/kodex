// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources_test

import (
	"encoding/hex"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
)

// Tests updating of streams by a stream superuser.
func TestUpdateConfig(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		// we create a test stream
		pt.FC{pf.Stream{Name: "test 1", Project: "project"}, "stream"},
		pt.FC{pf.Config{Name: "test 1", Stream: "stream"}, "config"},

		pt.FC{af.Organization{Name: "A"}, "org"},

		// we create a stream role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", Organization: "org", ObjectRole: "superuser"}, "projectRole"},
		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:config:write"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)
	defer pt.TeardownFixtures(fixturesConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(api.UserProfile)
	config := fixtures["config"].(kodex.Config)
	controller := fixtures["controller"].(api.Controller)

	sourceData := map[string]interface{}{
		"name":        "another test",
		"description": "another test",
	}

	resp, err := at.Put(controller, user, "/v1/configs/"+hex.EncodeToString(config.ID()), sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	if err := config.Refresh(); err != nil {
		t.Fatal(err)
	}

	if config.Name() != "another test" || config.Description() != "another test" {
		t.Fatal("update did not work")
	}

	// we test only updating the name

	sourceData = map[string]interface{}{
		"name": "third test",
	}

	resp, err = at.Put(controller, user, "/v1/configs/"+hex.EncodeToString(config.ID()), sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	if err := config.Refresh(); err != nil {
		t.Fatal(err)
	}

	if config.Name() != "third test" || config.Description() != "another test" {
		t.Fatal("update did not work")
	}

	// we test only updating the description

	sourceData = map[string]interface{}{
		"description": "third test",
	}

	resp, err = at.Put(controller, user, "/v1/configs/"+hex.EncodeToString(config.ID()), sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	if err := config.Refresh(); err != nil {
		t.Fatal(err)
	}

	if config.Name() != "third test" || config.Description() != "third test" {
		t.Fatal("update did not work")
	}

	// we test using invalid data

	sourceData = map[string]interface{}{
		"name": "1",
	}

	resp, err = at.Put(controller, user, "/v1/configs/"+hex.EncodeToString(config.ID()), sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 400 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}
