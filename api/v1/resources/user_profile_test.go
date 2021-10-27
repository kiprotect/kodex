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
	"encoding/json"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
)

type UserProfile struct {
}

func TestUserProfile(t *testing.T) {

	var fixturesConfig = []pt.FC{

		// we create the settings
		pt.FC{pf.Settings{}, "settings"},
		// we create an SQL controller
		pt.FC{af.Controller{}, "controller"},
		pt.FC{af.Organization{Name: "test"}, "org"},
		// we create a user
		pt.FC{af.User{EMail: "max@mustermann.de", Organization: "org", Roles: []string{"admin"}, Scopes: []string{"kiprotect:api"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(fixturesConfig)

	user := fixtures["user"].(api.UserProfile)
	controller := fixtures["controller"].(api.Controller)

	resp, err := at.Get(controller, user, "/v1/user", map[string]interface{}{})

	if err != nil {
		t.Fatal(err)
	}

	var values map[string]interface{}

	if err = json.Unmarshal(resp.Body.Bytes(), &values); err != nil {
		t.Fatal("invalid JSON")
	}

	_, ok := maps.ToStringMap(values["data"])

	if !ok {
		t.Fatal("no data found")
	}

}
