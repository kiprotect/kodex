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

package actions

import (
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"math"
)

type QuantizeAction struct {
	kodex.BaseAction
	config *QuantizeConfig
}

var QuantizeForm = forms.Form{
	ErrorMsg: "invalid data encountered in the encode/decode form",
	Fields: []forms.Field{
		{
			Name: "key",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "precision",
			Validators: []forms.Validator{
				forms.IsFloat{HasMin: true, Min: 0},
			},
		},
	},
}

type QuantizeConfig struct {
	Key       string  `json:"key"`
	Precision float64 `json:"precision"`
}

func MakeQuantizeAction(spec kodex.ActionSpecification) (kodex.Action, error) {

	quantizeConfig := &QuantizeConfig{}

	if params, err := QuantizeForm.Validate(spec.Config); err != nil {
		return nil, err
	} else if err := QuantizeForm.Coerce(params, quantizeConfig); err != nil {
		return nil, err
	} else {
		return &QuantizeAction{
			BaseAction: kodex.MakeBaseAction(spec, "quantize"),
			config:     quantizeConfig,
		}, nil
	}
}

func (a *QuantizeAction) Params() interface{} {
	return nil
}

func (a *QuantizeAction) GenerateParams(key, salt []byte) error {
	return nil
}

func (a *QuantizeAction) SetParams(params interface{}) error {
	return nil
}

func (a *QuantizeAction) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {

	v, ok := item.Get(a.config.Key)

	if !ok {
		// key is missing
		return item, nil
	}

	f, ok := v.(float64)

	if !ok {
		return nil, fmt.Errorf("expected a float64 value")
	}

	rv := math.Round(f/a.config.Precision) * a.config.Precision

	item.Set(a.config.Key, rv)

	return item, nil

}
