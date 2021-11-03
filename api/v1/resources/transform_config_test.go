// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
)

var transformConfigFixtures = []pt.FC{

	// we create the settings
	pt.FC{pf.Settings{}, "settings"},

	// we create an SQL controller
	pt.FC{af.Controller{}, "controller"},

	pt.FC{af.Organization{Name: "test"}, "org"},

	pt.FC{pf.Project{Name: "test"}, "project"},
	// we create a test stream
	pt.FC{pf.Stream{Name: "test", Project: "project"}, "stream"},

	// we create a stream role
	pt.FC{
		af.ObjectRole{
			ObjectName:       "project",
			OrganizationRole: "project:admin",
			ObjectRole:       "admin",
			Organization:     "org",
		},
		"projectRole",
	},

	pt.FC{pf.Config{Stream: "stream", Name: "test", Version: "1", Source: "api", Status: "active"}, "config"},

	pt.FC{pf.ActionConfig{Name: "test", Type: "pseudonymize", Project: "project", Config: map[string]interface{}{
		"key":    "foo",
		"method": "merengue",
		"config": map[string]interface{}{},
	}}, "action0"},

	pt.FC{pf.ActionMap{
		Config: "config",
		Action: "action0",
		Index:  0,
	}, "actionMap"},

	// we create a user
	pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"project:admin"}, Scopes: []string{"kiprotect:api:config:transform"}}, "user"},
}

func TestTransformConfig(t *testing.T) {

	fixtures, err := pt.SetupFixtures(transformConfigFixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(api.User)
	config := fixtures["config"].(kodex.Config)
	controller := fixtures["controller"].(api.Controller)

	sourceItems := []map[string]interface{}{
		map[string]interface{}{
			"foo": "bar",
		},
	}

	sourceData := map[string]interface{}{
		"items": sourceItems,
	}

	resp, err := at.Post(controller, user, "/v1/configs/"+hex.EncodeToString(config.ID())+"/transform", sourceData)

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

func TestTransformConfigError(t *testing.T) {

	fixtures, err := pt.SetupFixtures(transformConfigFixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(api.User)
	config := fixtures["config"].(kodex.Config)
	controller := fixtures["controller"].(api.Controller)

	sourceItems := []map[string]interface{}{
		map[string]interface{}{
			"fooz": "bar",
		},
	}

	sourceData := map[string]interface{}{
		"items": sourceItems,
	}

	resp, err := at.Post(controller, user, "/v1/configs/"+hex.EncodeToString(config.ID())+"/transform", sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

}
