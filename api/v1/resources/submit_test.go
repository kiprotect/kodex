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
	"encoding/hex"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	at "github.com/kiprotect/kodex/api/testing"
	af "github.com/kiprotect/kodex/api/testing/fixtures"
	pt "github.com/kiprotect/kodex/helpers/testing"
	pf "github.com/kiprotect/kodex/helpers/testing/fixtures"
	"testing"
	"time"
)

func TestSubmit(t *testing.T) {

	var submitFixtures = []pt.FC{
		pt.FC{pf.Settings{}, "settings"},
		pt.FC{af.Controller{}, "controller"},
		pt.FC{af.Organization{Name: "test"}, "org"},
		pt.FC{pf.Project{Name: "test"}, "project"},
		pt.FC{pf.Stream{Name: "test 1", Project: "project"}, "stream"},
		pt.FC{af.ObjectRole{
			ObjectName:       "project",
			OrganizationRole: "project:admin",
			ObjectRole:       "superuser",
			Organization:     "org"}, "projectRole"},
		pt.FC{af.User{
			EMail:        "max@mustermann.de",
			Scopes:       []string{"kiprotect:api:stream:submit"},
			Organization: "org",
			Roles:        []string{"project:admin"}}, "user"},
	}

	fixtures, err := pt.SetupFixtures(submitFixtures)
	defer pt.TeardownFixtures(submitFixtures, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	user := fixtures["user"].(*api.User)
	stream := fixtures["stream"].(kodex.Stream)
	controller := fixtures["controller"].(api.Controller)

	items := []map[string]interface{}{
		map[string]interface{}{
			"foo": "bar",
		},
		map[string]interface{}{
			"fbdsfsfdsoo": "bar",
		},
		map[string]interface{}{
			"foo":  "basdfdsdsfr",
			"4343": float64(32424),
		},
	}

	sourceData := map[string]interface{}{
		"items": items,
	}

	resp, err := at.Post(controller, user, "/v1/submit/"+hex.EncodeToString(stream.ID()), sourceData)

	if err != nil {
		t.Fatal(err)
	}

	if resp.Code != 200 {
		t.Fatalf("wrong return code: %d", resp.Code)
	}

	channel := kodex.MakeInternalChannel()

	if err := channel.Setup(controller, stream); err != nil {
		t.Fatal(err)
	}
	var payload kodex.Payload

	i := 0
	for {
		time.Sleep(time.Millisecond)
		if payload, err = channel.Read(); err != nil {
			t.Fatal(err)
		}
		if payload != nil {
			break
		}
		i++
		if i > 1000 {
			t.Fatal("no payload received")
		}
	}
	newItems := payload.Items()
	if len(newItems) != len(items) {
		t.Fatal("item count does not match")
	}
	for i, newItem := range newItems {
		item := items[i]
		ni := newItem.All()
		for k, v := range item {
			if ni[k] != v {
				t.Fatalf("key %s does not match", k)
			}
		}
	}
	// we acknowledge the payload
	if err := payload.Acknowledge(); err != nil {
		t.Fatal(err)
	}

}
