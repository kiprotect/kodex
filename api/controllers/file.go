package controllers

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/controllers"
)

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
