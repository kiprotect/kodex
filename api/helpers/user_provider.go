// IRIS Endpoint-Server (EPS)
// Copyright (C) 2021-2021 The IRIS Endpoint-Server Authors (see AUTHORS.md)
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

package helpers

import (
	"fmt"
	"github.com/kiprotect/kodex/api"
)

func UserProvider(controller api.Controller) (api.UserProvider, error) {
	userProviderType, ok := controller.Settings().String("user-provider.type")
	if !ok {
		return nil, fmt.Errorf("user provider type missing")
	}
	definition, ok := controller.APIDefinitions().UserProviders[userProviderType]
	if !ok {
		return nil, fmt.Errorf("invalid user provider type '%s'", userProviderType)
	}
	return definition.Maker(controller.Settings())
}
