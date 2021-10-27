// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources

import (
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

var TransformForm = forms.Form{
	ErrorMsg: "invalid data encountered in the transform form",
	Fields: []forms.Field{
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
