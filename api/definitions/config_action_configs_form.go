// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package definitions

import (
	"github.com/kiprotect/go-helpers/forms"
)

var AddConfigActionConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the config action config form",
	Fields: []forms.Field{
		{
			Name: "index",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsInteger{HasMin: true, Min: 0},
			},
		},
	},
}
