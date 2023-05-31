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

package kodex

import (
	"os"
	"plugin"
	"testing"
)

type SimpleConfig interface {
	ID() []byte
}

type SimpleWriter interface {
	Teardown() error
	Write([]interface{}) error
	Setup(interface{}) error
}

type MyConfig struct {
	id []byte
}

func (c *MyConfig) ID() []byte {
	return c.id
}

func TestPlugin(t *testing.T) {

	if os.Getenv("KODEX_PLUGIN_TEST") == "" {
		t.Skip("Skipping plugin test")
	}

	p, err := plugin.Open("plugins/writers/example/example.so")
	if err != nil {
		t.Fatal(err)
	}
	v, err := p.Lookup("MakeWriter")

	if err != nil {
		t.Fatal(err)
	}

	maker, ok := v.(func(map[string]interface{}) (interface{}, error))

	if !ok {
		t.Fatalf("not a writer maker")
	} else {
		Log.Info("this is a writer maker!")
	}

	writer, err := maker(map[string]interface{}{})

	if err != nil {
		t.Fatal(err)
	}

	simpleWriter, ok := writer.(SimpleWriter)

	if !ok {
		t.Fatalf("not a simple writer")
	} else {
		Log.Info("this is a simple writer")
	}

	if err := simpleWriter.Setup(&MyConfig{
		id: []byte("hey"),
	}); err != nil {
		t.Fatal(err)
	}

}
