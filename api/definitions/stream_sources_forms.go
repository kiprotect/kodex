// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package definitions

import (
	"github.com/kiprotect/go-helpers/forms"
)

var AddStreamSourceForm = forms.Form{
	ErrorMsg: "invalid data encountered in the stream source adding form",
	Fields: []forms.Field{
		{
			Name: "status",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsIn{Choices: []interface{}{"active", "disabled", "testing"}},
			},
		},
	},
}
