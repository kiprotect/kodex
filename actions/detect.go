// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions/pseudonymize"
	"regexp"
	"strings"
)

type DetectAction struct {
	kodex.BaseAction
	config        *DetectConfig
	pseudonymizer pseudonymize.Pseudonymizer
}

var DetectForm = forms.Form{
	ErrorMsg: "invalid data encountered in the detect form",
	Fields: []forms.Field{
		{
			Name: "key",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "_"},
				forms.IsString{},
			},
		},
		{
			Name:        "format",
			Description: "The format to use for replaced strings.",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "%s"},
				forms.IsString{},
			},
		},
		{
			Name:        "action",
			Description: "The action that should be taken for detected data.",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "pseudonymize"},
				forms.IsIn{
					Choices: []interface{}{"pseudonymize", "mask"},
				},
			},
		},
	},
}

type DetectConfig struct {
	Key    string `json:"key"`
	Format string `json:"format"`
	Action string `json:"action"`
}

var regexes = map[string]*regexp.Regexp{
	"ip":    regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`),
	"email": regexp.MustCompile(`[^@\s]+@[^@\s]+\.[a-z]+`),
	"iban":  regexp.MustCompile(`DE\d{7,20}`),
}

func MakeDetectAction(spec kodex.ActionSpecification) (kodex.Action, error) {

	detectConfig := &DetectConfig{}

	if params, err := DetectForm.Validate(spec.Config); err != nil {
		return nil, err
	} else if err := DetectForm.Coerce(detectConfig, params); err != nil {
		return nil, err
	} else {
		da := &DetectAction{
			BaseAction: kodex.MakeBaseAction(spec, "detect"),
			config:     detectConfig,
		}
		if detectConfig.Action == "pseudonymize" {
			pseudonymizer, err := pseudonymize.MakeMerenguePseudonymizer(map[string]any{})
			if err != nil {
				return nil, fmt.Errorf("cannot create pseudonymizer: %v", err)
			}
			da.pseudonymizer = pseudonymizer
		}

		return da, nil
	}
}

func (a *DetectAction) Params() interface{} {
	switch a.config.Action {
	case "pseudonymize":
		return a.pseudonymizer.Params()
	}
	return nil
}

func (a *DetectAction) GenerateParams(key, salt []byte) error {
	switch a.config.Action {
	case "pseudonymize":
		return a.pseudonymizer.GenerateParams(key, salt)
	}
	return nil
}

func (a *DetectAction) SetParams(params interface{}) error {
	switch a.config.Action {
	case "pseudonymize":
		return a.pseudonymizer.SetParams(params)
	}
	return nil
}

func (a *DetectAction) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {

	input, ok := item.Get(a.config.Key)

	if !ok {
		// key is missing
		return nil, fmt.Errorf("key missing")
	}

	inputStr, ok := input.(string)

	if !ok {
		return nil, fmt.Errorf("input is not a string")
	}

	replace := func(value string) string {
		var out string
		switch a.config.Action {
		case "pseudonymize":
			output, err := a.pseudonymizer.Pseudonymize(value)
			if err != nil {
				out = fmt.Sprintf("error: %v", err)
			}
			out = output.(string)
		case "mask":
			out = strings.Repeat("*", len(value))
		default:
			out = value
		}
		return fmt.Sprintf(a.config.Format, out)
	}

	for _, regexp := range regexes {
		inputStr = regexp.ReplaceAllStringFunc(inputStr, replace)
	}

	item.Set(a.config.Key, inputStr)

	return item, nil

}
