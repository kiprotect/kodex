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

package writers

import (
	"github.com/kiprotect/go-helpers/forms"
)

var CountWriterForm = forms.Form{
	ErrorMsg: "invalid data encountered in the count writer form",
	Fields: []forms.Field{
		{
			Name: "max",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 10},
				forms.IsInteger{HasMin: true, Min: 0, HasMax: true, Max: 100},
			},
		},
	},
}
