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

package api

import (
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
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
				forms.IsBytes{Encoding: "hex"},
			},
		},
		{
			Name: "organizationSource",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "inMemory"},
				forms.IsIn{Choices: []interface{}{"inMemory"}},
			},
		},
		{
			Name: "objectID",
			Validators: []forms.Validator{
				forms.IsBytes{Encoding: "hex"},
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
			Name: "source",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "inMemory"},
				forms.IsIn{Choices: []interface{}{"inMemory"}},
			},
		},
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "default",
			Validators: []forms.Validator{
				forms.IsOptional{Default: false},
				forms.IsBoolean{},
			},
		},
		{
			Name: "id",
			Validators: []forms.Validator{
				forms.IsBytes{Encoding: "hex"},
			},
		},
	},
}

var AccessTokenForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "token",
			Validators: []forms.Validator{
				forms.IsBytes{Encoding: "hex"},
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

var RolesForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "roles",
			Validators: []forms.Validator{
				forms.IsStringList{},
			},
		},
		{
			Name: "organization",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &BPOrganizationForm,
				},
			},
		},
	},
}

var BlueprintUserForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "source",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "inMemory"},
				forms.IsIn{Choices: []interface{}{"inMemory"}},
			},
		},
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
							Form: &RolesForm,
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
				forms.IsOptional{},
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &BlueprintUserForm,
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
	},
}

type Blueprint struct {
	config map[string]interface{}
}

type ObjectRoleSpec struct {
	ObjectRole         string `json:"objectRole"`
	OrganizationRole   string `json:"organizationRole"`
	OrganizationID     []byte `json:"organizationID"`
	OrganizationSource string `json:"organizationSource"`
	ObjectID           []byte `json:"objectID"`
	ObjectType         string `json:"objectType"`
}

type BlueprintSpec struct {
	Users []*ExternalUser   `json:"users"`
	Roles []*ObjectRoleSpec `json:"roles"`
}

func InitRoles(controller Controller, roles []*ObjectRoleSpec) error {
	for _, role := range roles {
		var obj kodex.Model
		var err error
		switch role.ObjectType {
		case "project":
			obj, err = controller.Project(role.ObjectID)
			if err != nil {
				kodex.Log.Error("project not found")
				return err
			}
		default:
			return fmt.Errorf("invalid object type: %s", role.ObjectType)
		}
		if org, err := controller.Organization(role.OrganizationSource, role.OrganizationID); err != nil {
			return err
		} else {
			objRole := controller.MakeObjectRole(obj, org)
			if err := objRole.SetObjectRole(role.ObjectRole); err != nil {
				return err
			} else if err := objRole.SetOrganizationRole(role.OrganizationRole); err != nil {
				return err
			} else if err := objRole.Save(); err != nil {
				return err
			}
		}

	}
	return nil
}

func InitUsers(controller Controller, users []*ExternalUser) error {

	if users == nil {
		return nil
	}

	userProvider, err := controller.UserProvider()

	if err != nil {
		return err
	}

	createUserProvider, ok := userProvider.(CreateUserProvider)

	if !ok {
		return fmt.Errorf("cannot create users")
	}

	for _, user := range users {
		for _, roles := range user.Roles {
			// we ensure all organizations are created in the controller
			if _, err := roles.Organization.ApiOrganization(controller); err != nil {
				return err
			}
		}
		if err := createUserProvider.Create(user); err != nil {
			return err
		}
	}
	return nil
}

func MakeBlueprint(config map[string]interface{}) *Blueprint {
	return &Blueprint{
		config: config,
	}
}

func (b *Blueprint) Create(controller Controller) error {

	if params, err := BlueprintConfigForm.Validate(b.config); err != nil {
		return err
	} else {
		spec := &BlueprintSpec{}
		if err := BlueprintConfigForm.Coerce(spec, params); err != nil {
			return err
		}
		// we initialize users
		if err := InitUsers(controller, spec.Users); err != nil {
			return err
		}
		// we initialize roles
		if err := InitRoles(controller, spec.Roles); err != nil {
			return err
		}
		return nil
	}
}
