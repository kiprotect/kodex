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

package fixtures

import (
	"fmt"
	"github.com/kiprotect/kiprotect"
)

type Blueprint struct {
	Config  map[string]interface{}
	Project string
}

func (c Blueprint) Setup(fixtures map[string]interface{}) (interface{}, error) {

	project, ok := fixtures[c.Project].(kiprotect.Project)

	if !ok {
		return nil, fmt.Errorf("project is missing")
	}

	blueprint := kiprotect.MakeBlueprint(c.Config)

	if err := blueprint.Create(project); err != nil {
		return nil, err
	}

	return blueprint, nil
}

func (c Blueprint) Teardown(fixture interface{}) error {
	return nil
}
