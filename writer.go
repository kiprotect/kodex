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
	"github.com/kiprotect/go-helpers/forms"
)

type WriterDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Internal    bool        `json:"internal"`
	Maker       WriterMaker `json:"-"`
	Form        forms.Form  `json:"form"`
}

type WriterMaker func(map[string]interface{}) (Writer, error)
type WriterDefinitions map[string]WriterDefinition

type Writer interface {
	Write(payload Payload) error
	Setup(Config) error
	Teardown() error
}

type ClosableWriter interface {
	Close() error
}

// A writer that is able to write objects for a specific model such as a stream
// or an destination. Used for internal data routing.
type ModelWriter interface {
	Writer
	SetupWithModel(Model) error
}
