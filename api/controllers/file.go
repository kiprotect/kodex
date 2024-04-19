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

package controllers

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/controllers"
)

/*
The API file controller persists all data in a single API blueprint file,
which can be checked into version control and used to manage all API-related data.
*/

type FileController struct {
	*InMemoryController
}

func (f *FileController) OnUpdate(object kodex.Model) {
	// get the corresponding project, export it as a blueprint to a file. Done!
	kodex.Log.Info("API object update")
}

func MakeFileController(config map[string]interface{}, controller kodex.Controller, definitions *api.Definitions) (api.Controller, error) {

	kodexFileController, ok := controller.(*controllers.FileController)

	if !ok {
		return nil, fmt.Errorf("parent controller is not a file controller")
	}

	controller, err := MakeInMemoryController(config, kodexFileController.InMemoryController, definitions)

	if err != nil {
		return nil, err
	}

	// we coerce into an in-memory controller
	inMemoryController := controller.(*InMemoryController)

	fileController := FileController{
		InMemoryController: inMemoryController,
	}

	// we subscribe to updates from the in-memory controller
	inMemoryController.SetOnApiUpdate(fileController.OnUpdate)

	return &fileController, nil
}
