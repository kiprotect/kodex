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

package api

import (
	"github.com/kiprotect/kodex"
)

type User interface {
	kodex.Model
	Source() string
	SourceID() []byte
	Email() string
	DisplayName() string
	Superuser() bool
	SetDisplayName(string) error
	SetEmail(string) error
	SetSuperuser(bool) error
	SetSource(string) error
	SetSourceID([]byte) error
	Data() interface{}
	SetData(interface{}) error
}
