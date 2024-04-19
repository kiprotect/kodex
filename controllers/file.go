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
	"github.com/kiprotect/kodex"
)

/*
The file based controller uses the in-memory controller and syncs all changes
done by it to disk, to a collection of YAML or JSON files. Each file corresponds
to a Project blueprint and contains all the associated data, so it's easy to
e.g. use the project files with the command line version or manage them in
version control.

- Is there an easier way to do this? We just need to capture
*/

type FileController struct {
	*InMemoryController
}

func (f *FileController) OnUpdate(object kodex.Model) {
	// get the corresponding project, export it as a blueprint to a file. Done!
}

func MakeFileController(config map[string]interface{}, settings kodex.Settings, definitions *kodex.Definitions) (kodex.Controller, error) {

	controller, err := MakeInMemoryController(config, settings, definitions)

	if err != nil {
		return nil, err
	}

	// we coerce into an in-memory controller
	inMemoryController := controller.(*InMemoryController)

	fileController := FileController{
		InMemoryController: inMemoryController,
	}

	// we subscribe to updates from the in-memory controller
	inMemoryController.SetOnUpdate(fileController.OnUpdate)

	return &fileController, nil
}
