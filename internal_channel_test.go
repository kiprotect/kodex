// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
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

package kiprotect_test

import (
	"github.com/kiprotect/kiprotect"
	pt "github.com/kiprotect/kiprotect/helpers/testing"
	pf "github.com/kiprotect/kiprotect/helpers/testing/fixtures"
	"testing"
	"time"
)

func TestAPIReader(t *testing.T) {

	var fixtureConfig = []pt.FC{
		pt.FC{&pf.Settings{}, "settings"},
		pt.FC{&pf.Controller{}, "controller"},
		pt.FC{&pf.Project{Name: "test"}, "project"},
		pt.FC{&pf.Stream{Name: "test", Project: "project"}, "stream"},
	}

	fixtures, err := pt.SetupFixtures(fixtureConfig)
	defer pt.TeardownFixtures(fixtureConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	controller := fixtures["controller"].(kiprotect.Controller)
	stream := fixtures["stream"].(kiprotect.Stream)

	channel := kiprotect.MakeInternalChannel()

	if err := channel.Setup(controller, stream); err != nil {
		t.Fatal(err)
	}
	if err := channel.Purge(); err != nil {
		t.Fatal(err)
	}

	items := []*kiprotect.Item{
		kiprotect.MakeItem(map[string]interface{}{
			"foo": "bar",
		}),
	}

	if err := channel.Write(kiprotect.MakeBasicPayload(items, map[string]interface{}{}, false)); err != nil {
		t.Fatal(err)
	}

	var payload kiprotect.Payload

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

	// we acknowledge the payload
	if err := payload.Acknowledge(); err != nil {
		t.Fatal(err)
	}

}
