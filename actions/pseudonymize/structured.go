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

package pseudonymize

import (
	"encoding/base64"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions/pseudonymize/structured"
)

type StructuredPseudonymizer struct {
	Type             string
	prefixPreserving bool `json:"preserve-prefixes"`
	TypeParams       interface{}
	Format           string
	key              []byte
	defaultKey       []byte
	Maker            structured.TypeMaker
}

var StructuredPseudonymizerForm = forms.Form{
	ErrorMsg: "invalid data encountered in the structured pseudonymizer form",
	Fields: []forms.Field{
		{
			Name:        "key",
			Description: "The encryption key to use for this pseudonymizer (will override any default key).",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{},
			},
		},
		{
			Name:        "preserve-prefixes",
			Description: "Whether to preserve structural prefixes of a value. For example, if a date is formatted as %Y-%m-%d and prefix preservation is enabled, input dates sharing a prefix e.g. 2000-10-07 and 2000-10-13 will be mapped to pseudonyms also sharing a prefix of the same length, e.g. 2022-05-11, 2022-05-30.",
			Validators: []forms.Validator{
				forms.IsOptional{Default: false},
				forms.IsBoolean{},
			},
		},
		{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsIn{
					Choices: []interface{}{"ip", "date", "integer", "ipv4", "ipv6"},
				},
			},
		},
		{
			Name: "type-params",
			Validators: []forms.Validator{
				forms.Switch{
					Key: "type",
					Cases: map[string][]forms.Validator{
						"integer": []forms.Validator{
							forms.IsStringMap{
								Form: &forms.Form{
									Fields: []forms.Field{
										{
											Name: "min",
											Validators: []forms.Validator{
												forms.IsInteger{},
											},
										},
										{
											Name: "max",
											Validators: []forms.Validator{
												forms.IsInteger{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "format",
			Validators: []forms.Validator{
				forms.Switch{
					Key: "type",
					Cases: map[string][]forms.Validator{
						"date": []forms.Validator{
							forms.IsOptional{Default: "%Y-%m-%dT%H:%M:%SZ"},
							forms.IsString{},
						},
					},
					// to do: remove this
					Default: []forms.Validator{
						forms.IsNil{
							AllowNull: true,
						},
					},
				},
			},
		},
	},
}

func MakeStructuredPseudonymizer(config map[string]interface{}) (Pseudonymizer, error) {

	if config == nil {
		config = map[string]any{
			"type": "integer",
			"type-params": map[string]any{
				"min": 0,
				"max": 100,
			},
		}
	}

	var ok bool

	params, err := StructuredPseudonymizerForm.Validate(config)

	if err != nil {
		return nil, err
	}

	type_ := params["type"]
	format := params["format"]
	typeParams := params["type-params"]
	prefixPreserving := params["preserve-prefixes"].(bool)
	var defaultKey []byte

	var strType, strFormat string
	if strType, ok = type_.(string); !ok {
		return nil, fmt.Errorf("type: expected a string")
	}

	strFormat, _ = format.(string)

	if _, ok = structured.Types[strType]; !ok {
		return nil, fmt.Errorf("unknown type: %s", strType)
	}

	if params["key"] != nil {
		strKey := params["key"].(string)
		defaultKey = []byte(strKey)
	}

	return &StructuredPseudonymizer{
		Format:           strFormat,
		Type:             strType,
		TypeParams:       typeParams,
		Maker:            structured.Types[strType],
		defaultKey:       defaultKey,
		prefixPreserving: prefixPreserving,
	}, nil
}

func (p *StructuredPseudonymizer) Params() interface{} {
	return map[string]interface{}{
		"key": base64.StdEncoding.EncodeToString(p.key),
	}
}

func (p *StructuredPseudonymizer) SetParams(params interface{}) error {
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

func (p *StructuredPseudonymizer) GenerateParams(key, salt []byte) error {
	if key == nil {
		randomBytes, err := kodex.RandomBytes(64)
		if err != nil {
			return err
		}
		key = randomBytes
	}
	p.key = kodex.DeriveKey(key, salt, 64)
	return nil
}

func (p *StructuredPseudonymizer) Pseudonymize(value interface{}) (interface{}, error) {
	f := structured.PS
	if p.prefixPreserving {
		f = structured.PSH
	}
	return p.process(value, f)
}

func (p *StructuredPseudonymizer) Depseudonymize(value interface{}) (interface{}, error) {
	f := structured.DPS
	if p.prefixPreserving {
		f = structured.DPSH
	}
	return p.process(value, f)
}

func (p *StructuredPseudonymizer) process(value interface{}, f func(structured.CompositeType, []byte) (structured.CompositeType, error)) (interface{}, error) {

	compositeType, err := p.Maker(p.TypeParams)
	if err != nil {
		return nil, err
	}

	if err = compositeType.Unmarshal(p.Format, value); err != nil {
		return nil, err
	}

	valid := compositeType.IsValid()

	for _, v := range valid {
		if !v {
			return nil, fmt.Errorf("input value is invalid")
		}
	}

	pseudonymizedType, err := f(compositeType, p.key)

	if err != nil {
		return nil, err
	}

	pseudonymizedValue, err := pseudonymizedType.Marshal(p.Format)

	if err != nil {
		return nil, err
	}

	return pseudonymizedValue, nil
}
