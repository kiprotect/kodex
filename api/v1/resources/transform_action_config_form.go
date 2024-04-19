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

package resources

import (
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

var TransformActionConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the transform action config form",
	Fields: []forms.Field{
		{
			Name: "items",
			Validators: []forms.Validator{
				// first we validate that this is a lit
				forms.IsList{
					// then we validate that each entry is an item
					Validators: []forms.Validator{
						kodex.IsItem{},
					},
				},
				// then we cast []interface{} -> []*kodex.Item
				kodex.IsItems{},
			},
		},
	},
}
