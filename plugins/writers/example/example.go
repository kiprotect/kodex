// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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

package main

import (
    "fmt"
)

type TestWriter struct {}

type Config interface {
    ID() []byte
}

type Item interface {
    Get(string) (interface{}, error)
}

func MakeWriter(config map[string]interface{}) (interface{}, error) {
    return &TestWriter{}, nil
}

// test
func (t *TestWriter) Setup(data interface{}) error {
    if config, ok := data.(Config); !ok {
        return fmt.Errorf("not a config object")
    } else {
        fmt.Printf("Setting up writer with config %s", string(config.ID()))
    }
    return nil
}

func (t *TestWriter) Teardown() error {
    return nil
}

func (t *TestWriter) Validate() error {
    return nil
}

func (t *TestWriter) Write(items []interface{}) error {
    return nil
}

func main() {
    // this does nothing, just needs to be defined for "go get" to work
}
