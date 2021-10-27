// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"github.com/kiprotect/kodex"
)

type Organization interface {
	kodex.Model
	Source() string
	SourceID() []byte
	Name() string
	SetName(string) error
	Description() string
	SetDescription(string) error
	SetSource(string) error
	SetSourceID([]byte) error
	Data() interface{}
	SetData(interface{}) error
}
