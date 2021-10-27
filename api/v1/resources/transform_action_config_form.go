// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

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
