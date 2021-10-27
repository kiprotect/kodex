// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"github.com/kiprotect/kodex"
)

type ObjectRole interface {
	kodex.Model
	OrganizationID() []byte
	ObjectID() []byte
	OrganizationRole() string
	SetOrganizationRole(string) error
	ObjectRole() string
	SetObjectRole(string) error
	ObjectType() string
}
