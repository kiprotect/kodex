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

package functions

import (
	"crypto/rand"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
	"math"
	"math/big"
)

type Uniques struct {
	idField string
	epsilon float64
}

var UniquesForm = forms.Form{
	ErrorMsg: "invalid data encountered in the uniques config",
	Fields: []forms.Field{
		{
			Name: "id",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsString{},
			},
		},
		{
			Name: "epsilon",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 0.5},
				forms.IsFloat{HasMin: true, Min: 0.01, HasMax: false},
			},
		},
	},
}

func (c *Uniques) Initialize(group aggregate.Group) error {
	mapState := &StringBoolMap{M: make(map[string]bool)}
	group.Lock()
	defer group.Unlock()
	return group.Initialize(mapState)
}

func (c *Uniques) Merge(groups []aggregate.Group) (aggregate.Group, error) {
	if len(groups) == 1 {
		return groups[0], nil
	}
	newGroup := groups[0]
	newGroup.Lock()
	defer newGroup.Unlock()
	mapState, ok := newGroup.State().(*StringBoolMap)
	if !ok {
		return nil, fmt.Errorf("Expected a string map")
	}
	for i, group := range groups {
		if i == 0 {
			continue
		}
		group.Lock()
		otherMapState, ok := group.State().(*StringBoolMap)
		if !ok {
			return nil, fmt.Errorf("Expected a string map")
		}
		for k, v := range otherMapState.M {
			mapState.M[k] = v
		}
		group.Unlock()
	}
	return newGroup, nil
}

func (c *Uniques) Add(item *kodex.Item, group aggregate.Group) error {
	group.Lock()
	defer group.Unlock()
	state := group.State()
	mapState, ok := state.(*StringBoolMap)
	if !ok {
		return fmt.Errorf("Expected a string map")
	}
	idValue, ok := item.Get(c.idField)
	if !ok {
		return nil
	}
	idValueStr, ok := idValue.(string)
	if !ok {
		return fmt.Errorf("Expected a string ID")
	}
	mapState.M[idValueStr] = true
	return nil
}

func uniform() (float64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return 0, err
	}
	return float64(n.Int64()) / float64(1<<62), nil
}

func geometricNoise(epsilon float64, symmetric bool) (int64, error) {
	var k int64
	p := math.Exp(-epsilon)
	if pv, err := uniform(); err != nil {
		return 0, err
	} else {
		if symmetric {
			if pv < (1-p)/(1+p) {
				return 0, nil
			}
		} else if pv > p {
			return 0, nil
		}
	}
	if p < 1e-6 {
		return 0, nil
	}
	if pv, err := uniform(); err != nil {
		return 0, err
	} else {
		pe := 1.0 - p + p*pv
		k = int64(math.Log(1-pe) / math.Log(p))

		if symmetric {
			if pv, err := uniform(); err != nil {
				return 0, err
			} else if pv < 0.5 {
				k = -k
			}
		}

		return k, nil

	}
}

func (c *Uniques) Finalize(group aggregate.Group) (interface{}, error) {
	group.Lock()
	defer group.Unlock()
	state := group.State()
	mapState, ok := state.(*StringBoolMap)
	if !ok {
		return nil, fmt.Errorf("Expected a string map")
	}
	intState := int64(len(mapState.M))
	// we add symmetric geometric noise to the result, which makes it differentially private
	if noise, err := geometricNoise(c.epsilon, true); err != nil {
		return nil, err
	} else {
		if intState+noise < 0 {
			return 0, nil
		}
		return intState + noise, nil
	}
}

func MakeUniquesFunction(config map[string]interface{}) (aggregate.Function, error) {
	params, err := UniquesForm.Validate(config)
	if err != nil {
		return nil, err
	}
	return &Uniques{
		idField: params["id"].(string),
		epsilon: params["epsilon"].(float64),
	}, nil
}
