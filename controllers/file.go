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
