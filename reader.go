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

package kiprotect

import (
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
)

type ReaderDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Internal    bool        `json:"internal"`
	Maker       ReaderMaker `json:"-"`
	Form        forms.Form  `json:"form"`
}

type ReaderMaker func(map[string]interface{}) (Reader, error)
type ReaderDefinitions map[string]ReaderDefinition

type Reader interface {
	Read() (Payload, error)
	Setup(Stream) error
	Purge() error
	Teardown() error
}

var EOS = fmt.Errorf("end of stream")

// A SchemaReader is able to generate its own schema
type SchemaReader interface {
	Schema() (DataSchema, error)
}

// A peeking reader is able to "peek" into the data stream, i.e. to read a
// payload but immediately put it back to
type PeekingReader interface {
	// Read a payload but immediately reject it (if possible)
	Peek() (Payload, error)
	Reader
}

// A reader that is able to write objects for a specific model such as a stream.
type ModelReader interface {
	Reader
	SetupWithModel(Model) error
}
