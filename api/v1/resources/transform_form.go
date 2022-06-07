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

package resources

import (
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

var TransformForm = forms.Form{
	ErrorMsg: "invalid data encountered in the transform form",
	Fields: []forms.Field{
		{
			Name: "undo",
			Validators: []forms.Validator{
				forms.IsOptional{Default: false},
				forms.IsBoolean{},
			},
		},
		{
			Name: "key",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsBytes{
					Encoding: "base64",
				},
			},
		},
		{
			Name: "salt",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsBytes{
					Encoding: "base64",
				},
			},
		},
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
		{
			Name: "actions",
			Validators: []forms.Validator{
				forms.IsOptional{Default: []map[string]interface{}{}},
				// first we validate that this is a list
				forms.IsList{
					Validators: []forms.Validator{
						// then we validate that each entry is an action spec
						kodex.IsActionSpecification{},
					},
				},
				// then we cast []interface{} -> []kodex.ActionSpecification
				kodex.IsActionSpecifications{},
			},
		},
	},
}
