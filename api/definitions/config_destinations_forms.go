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

package definitions

import (
	"github.com/kiprotect/go-helpers/forms"
	"regexp"
)

var AddConfigDestinationForm = forms.Form{
	ErrorMsg: "invalid data encountered in the config destination form",
	Fields: []forms.Field{
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsIn{Choices: []interface{}{"active", "disabled", "testing"}},
			},
		},
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.MatchesRegex{Regexp: regexp.MustCompile("^[a-z0-9-]{3,40}$")},
			},
		},
	},
}
