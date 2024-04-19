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

package testing

import (
	"github.com/kiprotect/kodex"
)

type Fixture interface {
	Setup(map[string]interface{}) (interface{}, error)
	Teardown(interface{}) error
}

type FC struct {
	F    Fixture
	Name string
}

func TeardownFixtures(fixtureConfigs []FC, fixtures map[string]interface{}) error {
	var teardownErr error
	for _, fixtureConfig := range fixtureConfigs {
		if err := fixtureConfig.F.Teardown(fixtures[fixtureConfig.Name]); err != nil {
			kodex.Log.Errorf("error tearing down fixture %s", fixtureConfig.Name)
			teardownErr = err
		}
	}
	return teardownErr
}

func SetupFixtures(fixtureConfigs []FC) (map[string]interface{}, error) {

	fixtures := make(map[string]interface{})

	for _, fixtureConfig := range fixtureConfigs {
		var result interface{}
		var err error
		if result, err = fixtureConfig.F.Setup(fixtures); err != nil {
			kodex.Log.Errorf("error creating fixture %s", fixtureConfig.Name)
			return nil, err
		}
		if fixtureConfig.Name == "" {
			// we skip fixtures with empty names (they only provide side-effects)
			continue
		}
		oldValue := fixtures[fixtureConfig.Name]
		if oldValue != nil {
			if oldList, ok := oldValue.([]interface{}); ok {
				fixtures[fixtureConfig.Name] = append(oldList, result)
			} else {
				fixtures[fixtureConfig.Name] = []interface{}{oldValue, result}
			}
		} else {
			fixtures[fixtureConfig.Name] = result
		}
	}

	return fixtures, nil

}
