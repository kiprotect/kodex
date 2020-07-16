// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
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

package pseudonymize

import (
	"encoding/base64"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kiprotect"
	"github.com/kiprotect/kiprotect/actions/pseudonymize/merengue"
)

var MerengueConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the Merengue pseudonymizer form",
	Fields: []forms.Field{
		{
			Name: "key",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{},
			},
		},
		{
			Name: "encode",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "base64"},
				forms.IsIn{Choices: []interface{}{"base64"}},
			},
		},
	},
}

type MerenguePseudonymizer struct {
	key    []byte
	encode string
}

func toByteString(value interface{}) ([]byte, error) {
	var input []byte
	b, ok := value.([]byte)
	if ok {
		input = b
	} else {
		s, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("Merengue: Expected a string or byte array")
		}
		input = []byte(s)
	}
	return input, nil
}

func (p *MerenguePseudonymizer) Pseudonymize(value interface{}) (interface{}, error) {
	var input []byte
	var err error
	if input, err = toByteString(value); err != nil {
		return nil, err
	}
	result := merengue.Pseudonymize(input, uint(len(input)*8), p.key, merengue.Sha256)
	switch p.encode {
	case "base64":
		return base64.StdEncoding.EncodeToString(result), nil
	}
	return result, nil
}

func (p *MerenguePseudonymizer) Depseudonymize(value interface{}) (interface{}, error) {
	var input []byte
	var err error
	if input, err = toByteString(value); err != nil {
		return nil, err
	}
	switch p.encode {
	case "base64":
		if input, err = base64.StdEncoding.DecodeString(string(input)); err != nil {
			return nil, err
		}
	}
	return string(merengue.Depseudonymize(input, uint(len(input)*8), p.key, merengue.Sha256)), nil
}

func (p *MerenguePseudonymizer) GenerateParams(key, salt []byte) error {
	if key == nil {
		randomBytes, err := kiprotect.RandomBytes(64)
		if err != nil {
			return err
		}
		key = randomBytes
	}
	return p.GenerateParamsFromSeed(key, salt)
}

func (p *MerenguePseudonymizer) GenerateParamsFromSeed(key, salt []byte) error {
	p.key = kiprotect.DeriveKey(key, salt, 64)
	return nil
}

func (p *MerenguePseudonymizer) Params() interface{} {
	return map[string]interface{}{
		"key": base64.StdEncoding.EncodeToString(p.key),
	}
}

func (p *MerenguePseudonymizer) SetParams(params interface{}) error {
	paramsMap, ok := maps.ToStringMap(params)
	if !ok {
		return fmt.Errorf("Expected a map as parameters")
	}
	key, ok := paramsMap["key"]
	if !ok {
		return fmt.Errorf("Key missing from parameters map")
	}
	strKey, ok := key.(string)
	if !ok {
		return fmt.Errorf("Key should be a string or byte sequence")
	}
	byteKey, err := base64.StdEncoding.DecodeString(strKey)
	if err != nil {
		return err
	}
	p.key = byteKey
	return nil
}

func MakeMerenguePseudonymizer(config map[string]interface{}) (Pseudonymizer, error) {
	params, err := MerengueConfigForm.Validate(config)
	if err != nil {
		return nil, err
	}
	p := &MerenguePseudonymizer{}
	if params["encode"] != nil {
		p.encode = params["encode"].(string)
	}
	return p, nil
}
