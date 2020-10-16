// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
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

package actions

import (
	"bytes"
	"encoding/base64"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
)

var TranscodeConfigForm = forms.Form{
	ErrorMsg: "invalid data encountered in the encode/decode form",
	Fields: []forms.Field{
		{
			Name: "key",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
			},
		},
		{
			Name: "from",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
				forms.IsIn{
					Choices: []interface{}{"bytes", "string", "base64", "base64-url", "hex", "utf-8"},
				},
			},
		},
		{
			Name: "to",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
				forms.IsIn{
					Choices: []interface{}{"bytes", "string", "base64", "base64-url", "hex", "utf-8"},
				},
			},
		},
	},
}

type TranscodeAction struct {
	kodex.BaseAction
	from string
	to   string
	key  string
}

func base64Encode(source []byte, encoding *base64.Encoding) string {
	buffer := bytes.NewBuffer(make([]byte, 0))
	encoder := base64.NewEncoder(encoding, buffer)
	encoder.Write(source)
	encoder.Close()
	return string(buffer.Bytes())
}

func (t *TranscodeAction) Undo(item *kodex.Item) (*kodex.Item, error) {
	return nil, nil
}

func (t *TranscodeAction) Do(item *kodex.Item) (*kodex.Item, error) {
	return nil, nil
	/*
		value, ok := item.Get(t.key)

		if !ok {
			return nil, fmt.Errorf("key %s not found", t.key)
		}
		switch t.from {
		case "string":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("not a string")
			}
		case "bytes":
			bytesValue, ok := value.([]byte)
			if !ok {
				return nil, fmt.Errorf("not a bytestring")
			}
		case "base64":
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("not a string")
			}
			base64Value, err := base64.StdEncoding.DecodeString(string(byteValue))
			if err != nil {
				return nil, err

			}
			base64Value := base64Decode(byteValue, base64.StdEncoding)
			return strResult, nil
		case "utf-8":
			return string(byteValue), nil
		}
		return nil, fmt.Errorf("unknown/unsupported format")
		}

		byteValue, ok := value.([]byte)
		if !ok {
			strValue, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("encode/decode: expected a (byte) string")
			}
			byteValue = []byte(strValue)
		}
		switch t.format {
		case "base64":
			byteResult, err := base64.StdEncoding.DecodeString(string(byteValue))
			if err != nil {
				return nil, err

			}
			return byteResult, nil
		case "utf-8":
			return byteValue, nil
		}
		return nil, fmt.Errorf("unknown/unsupported format")
	*/

}

func (p *TranscodeAction) GenerateParams(key, salt []byte) error {
	return nil
}

func (p *TranscodeAction) SetParams(params interface{}) error {
	return nil
}

func (p *TranscodeAction) Params() interface{} {
	return nil
}

func MakeTranscodeAction(name, description string, id []byte, config map[string]interface{}) (kodex.Action, error) {

	params, err := TranscodeConfigForm.Validate(config)
	if err != nil {
		return nil, err
	}
	if baseAction, err := kodex.MakeBaseAction(name, description, "pseudonymize", id, config); err != nil {
		return nil, err
	} else {
		return &TranscodeAction{
			BaseAction: baseAction,
			from:       params["from"].(string),
			to:         params["to"].(string),
			key:        params["key"].(string),
		}, nil
	}
}
