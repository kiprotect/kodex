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
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
)

func TestCreateConfig(t *testing.T) {

	var createConfigFixtures = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},

		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},

		pt.FC{af.Organization{Name: "org"}, "org"},
		// we create two users and controllers
		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"admin", "superuser", "project:admin"}, Scopes: []string{"kiprotect:api:config:create"}}, "user"},
		// we create two users and controllers
		pt.FC{
			af.User{EMail: "mux@mastermann.de", Organization: "org", Roles: []string{"admins"}, Scopes: []string{"kiprotect:api"}}, "userB"},

		pt.FC{pf.Project{Name: "test"}, "project"},

		// we create two test streams
		pt.FC{pf.Stream{Name: "test", Project: "project"}, "stream"},

		// we create a stream role
		pt.FC{af.ObjectRole{ObjectName: "project", OrganizationRole: "project:admin", ObjectRole: "superuser", Organization: "org"}, "projectRole"},
	}

	fixtures, err := pt.SetupFixtures(createConfigFixtures)
	defer pt.TeardownFixtures(createConfigFixtures, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	controller := fixtures["controller"].(api.Controller)
	user := fixtures["user"].(*api.User)
	stream := fixtures["stream"].(kodex.Stream)

	sourceData := map[string]interface{}{
		"name":        "test",
		"description": "test",
	}

	resp, err := at.Post(controller, user, "/v1/streams/"+hex.EncodeToString(stream.ID())+"/configs", sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	var value map[string]interface{}

	if err = json.Unmarshal(resp.Body.Bytes(), &value); err != nil {
		t.Fatal("Invalid JSON")
	}

	var configMap map[string]interface{}

	if config, ok := value["data"]; !ok {
		t.Fatal("no config in response")
	} else {
		configMap, ok = config.(map[string]interface{})
		if !ok {
			t.Fatal("not a map")
		}
		if configMap["id"] == nil {
			t.Fatal("no ID in stream")
		}
	}

	configID, ok := configMap["id"].(string)

	if !ok {
		t.Fatal("ID is not a string")
	}

	configBinaryID, err := uuidToBytes(configID)

	if err != nil {
		t.Fatal(err)
	}

	config, err := controller.Config(configBinaryID)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(config.ID(), configBinaryID) {
		t.Fatalf("IDs do not match")
	}

}
