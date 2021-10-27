// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources

import (
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

var TransformConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the transform config form",
	Fields: []forms.Field{
		{
			Name: "items",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{},
						kodex.IsItem{},
					},
				},
			},
		},
	},
}
