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

package api

import (
	"github.com/kiprotect/go-helpers/forms"
)

var BPObjectRoleForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "objectType",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "organizationID",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "objectID",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "organizationRole",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "objectRole",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
	},
}

var BPOrganizationForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "id",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
	},
}

var AccessTokenForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "token",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "scopes",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsString{},
					},
				},
			},
		},
	},
}

var RoleForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "roles",
			Validators: []forms.Validator{
				forms.IsStringList{},
			},
		},
		{
			Name: "organizationID",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
	},
}

var UserForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "email",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{},
			},
		},
		{
			Name: "displayName",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{},
			},
		},
		{
			Name: "superuser",
			Validators: []forms.Validator{
				forms.IsOptional{Default: false},
				forms.IsBoolean{},
			},
		},
		{
			Name: "accessToken",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &AccessTokenForm,
				},
			},
		},
		{
			Name: "roles",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &RoleForm,
						},
					},
				},
			},
		},
	},
}

var BlueprintConfigForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "users",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &UserForm,
						},
					},
				},
			},
		},
		{
			Name: "roles",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &BPObjectRoleForm,
						},
					},
				},
			},
		},
		{
			Name: "organizations",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &BPOrganizationForm,
						},
					},
				},
			},
		},
	},
}

type Blueprint struct {
	config map[string]interface{}
}

func initRoles(controller Controller, config map[string]interface{}) error {
	return nil
}

func initUsers(controller Controller, config map[string]interface{}) error {
	return nil
}

func MakeBlueprint(config map[string]interface{}) *Blueprint {
	return &Blueprint{
		config: config,
	}
}

func (b *Blueprint) Create(controller Controller) error {
	if err := initUsers(controller, b.config); err != nil {
		return err
	}
	if err := initRoles(controller, b.config); err != nil {
		return err
	}
	return nil
}
