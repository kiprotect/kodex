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

package functions

import (
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
)

var CountForm = forms.Form{
	ErrorMsg: "invalid data encountered in the count config",
	Fields: []forms.Field{
		{
			Name: "epsilon",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 0.5},
				forms.IsFloat{HasMin: true, Min: 0.01, HasMax: false},
			},
		},
		{
			Name: "treshold",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 0},
				forms.IsInteger{HasMin: true, Min: 0, HasMax: false},
			},
		},
	},
}

type Count struct {
	epsilon  float64
	treshold int64
}

func (c *Count) Initialize(group aggregate.Group) error {
	intState := &Int64{
		I: 0,
	}
	group.Lock()
	defer group.Unlock()
	return group.Initialize(intState)
}

func (c *Count) Add(item *kodex.Item, group aggregate.Group) error {
	group.Lock()
	defer group.Unlock()
	state := group.State()
	intState, ok := state.(*Int64)
	if !ok {
		return errors.MakeInternalError("Expected an integer state", "COUNT", nil, nil)
	}
	intState.I += 1
	return nil
}

func (c *Count) Merge(groups []aggregate.Group) (aggregate.Group, error) {
	if len(groups) == 1 {
		return groups[0], nil
	}
	newGroup := groups[0]
	newGroup.Lock()
	defer newGroup.Unlock()
	intState, ok := newGroup.State().(*Int64)
	if !ok {
		return nil, errors.MakeInternalError("Expected an integer", "COUNT", nil, nil)
	}
	for i, group := range groups {
		if i == 0 {
			continue
		}
		group.Lock()
		otherIntState, ok := group.State().(*Int64)
		if !ok {
			return nil, errors.MakeInternalError("Expected an integer", "COUNT", nil, nil)
		}
		intState.I += otherIntState.I
		group.Unlock()
	}
	return newGroup, nil
}

func (c *Count) Finalize(group aggregate.Group) (interface{}, error) {
	group.Lock()
	defer group.Unlock()
	state := group.State()
	intState, ok := state.(*Int64)
	if !ok {
		return nil, errors.MakeInternalError("Expected an integer state", "COUNT", nil, nil)
	}
	i := intState.I
	var noise int64
	var err error
	if noise, err = geometricNoise(c.epsilon, true); err != nil {
		return nil, err
	}
	// we do not report values below or at the treshold
	if i+noise <= c.treshold {
		return nil, nil
	}
	return i + noise, nil
}

func MakeCountFunction(config map[string]interface{}) (aggregate.Function, error) {
	params, err := CountForm.Validate(config)
	if err != nil {
		return nil, err
	}
	return &Count{
		epsilon:  params["epsilon"].(float64),
		treshold: params["treshold"].(int64),
	}, nil
}
