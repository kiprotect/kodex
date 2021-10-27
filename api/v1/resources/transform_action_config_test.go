// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources_test

import (
	"encoding/hex"
	"encoding/json"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
)

var transformActionConfigFixtures = []pt.FC{

	// we create the settings
	pt.FC{pf.Settings{}, "settings"},

	// we create an SQL controller
	pt.FC{af.Controller{}, "controller"},

	pt.FC{af.Organization{Name: "test"}, "org"},

	pt.FC{pf.Project{Name: "test"}, "project"},

	pt.FC{pf.ActionConfig{Name: "test", Type: "pseudonymize", Project: "project", Config: map[string]interface{}{
		"key":    "foo",
		"method": "merengue",
		"config": map[string]interface{}{},
	}}, "action"},

	// we create an action role
	pt.FC{
		af.ObjectRole{
			ObjectName:       "project",
			OrganizationRole: "project:admin",
			ObjectRole:       "admin",
			Organization:     "org",
		},
		"projectRole",
	},

	// we create a user
	pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:action:transform"}}, "user"},
}

func TestTransformActionConfig(t *testing.T) {

	fixtures, err := pt.SetupFixtures(transformActionConfigFixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(api.UserProfile)
	action := fixtures["action"].(kodex.ActionConfig)
	controller := fixtures["controller"].(api.Controller)

	sourceItems := []map[string]interface{}{
		map[string]interface{}{
			"foo": "bar",
		},
	}

	sourceData := map[string]interface{}{
		"items": sourceItems,
	}

	resp, err := at.Post(controller, user, "/v1/actions/"+hex.EncodeToString(action.ID())+"/transform", sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	var values map[string]interface{}
	if err = json.Unmarshal(resp.Body.Bytes(), &values); err != nil {
		t.Fatal("invalid JSON")
	}

	data, ok := maps.ToStringMap(values["data"])

	if !ok {
		t.Fatal("data missing")
	}

	items, ok := maps.ToStringMapList(data["items"])

	if !ok {
		t.Fatal("items missing")
	}

	if len(items) != 1 {
		t.Fatal("expected 1 item")
	}

}

func TestTransformActionConfigError(t *testing.T) {

	fixtures, err := pt.SetupFixtures(transformActionConfigFixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(api.UserProfile)
	action := fixtures["action"].(kodex.ActionConfig)
	controller := fixtures["controller"].(api.Controller)

	sourceItems := []map[string]interface{}{
		map[string]interface{}{
			"fooz": "bar",
		},
	}

	sourceData := map[string]interface{}{
		"items": sourceItems,
	}

	resp, err := at.Post(controller, user, "/v1/actions/"+hex.EncodeToString(action.ID())+"/transform", sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}
