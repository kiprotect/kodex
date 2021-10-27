// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package resources

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
)

func objectType(c *gin.Context, key string) string {
	objectType, ok := c.Get(key)

	if !ok {
		return ""
	}

	objectTypeStr, _ := objectType.(string)

	// will return an empty string if the conversion failed
	return objectTypeStr
}

func getAdaptor(c *gin.Context) api.ObjectAdaptor {
	return getAdaptorForType(objectType(c, "objectType"), c)
}

func getAdaptorForType(objectType string, c *gin.Context) api.ObjectAdaptor {

	controller, ok := c.Get("controller")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("controller missing"))
		return nil
	}

	apiController, ok := controller.(api.Controller)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("not an API controller"))
		return nil
	}

	definitions := apiController.APIDefinitions()

	adaptor, ok := definitions.ObjectAdaptors[objectType]

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid object type"))
		return nil
	}

	return adaptor

}

func getAssociateAdaptor(c *gin.Context) api.AssociateAdaptor {

	leftType := objectType(c, "leftType")
	rightType := objectType(c, "rightType")

	associateName := fmt.Sprintf("%s-%s", leftType, rightType)

	controller, ok := c.Get("controller")

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("controller missing"))
		return nil
	}

	apiController, ok := controller.(api.Controller)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("not an API controller"))
		return nil
	}

	definitions := apiController.APIDefinitions()

	adaptor, ok := definitions.AssociateAdaptors[associateName]

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid associate object type"))
		return nil
	}

	return adaptor

}

func getAllObjs(c *gin.Context) []kodex.Model {

	adaptor := getAdaptor(c)

	if adaptor == nil {
		return nil
	}

	if allObjectsAdaptor, ok := adaptor.(api.ListAllObjectAdaptor); !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid adaptor type"))
		return nil
	} else {
		return allObjectsAdaptor.AllObjects(c)
	}

}

func GetObjs(c *gin.Context) []kodex.Model {

	adaptor := getAdaptor(c)

	if adaptor == nil {
		return nil
	}

	return adaptor.Objects(c)

}

func makeObject(c *gin.Context) kodex.Model {

	adaptor := getAdaptor(c)

	if adaptor == nil {
		return nil
	}

	createAdaptor, ok := adaptor.(api.CreateObjectAdaptor)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("cannot create objects of this type"))
		return nil
	}

	obj := createAdaptor.MakeObject(c)

	if obj == nil {
		api.HandleError(c, 500, fmt.Errorf("cannot create object"))
		return nil
	}

	return obj
}

// Return the object that was initialized via the ValidObject decorator
func GetObj(c *gin.Context, key string) kodex.Model {

	objectType := objectType(c, key)

	objectObj, ok := c.Get(objectType)

	if !ok {
		api.HandleError(c, 404, fmt.Errorf("object not found"))
		return nil
	}

	object, ok := objectObj.(kodex.Model)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid object"))
		return nil
	}

	return object

}

func getRoleObj(c *gin.Context) (kodex.Model, bool) {

	roleObjObj, ok := c.Get("roleObject")

	if !ok {
		return nil, true
	}

	roleObj, ok := roleObjObj.(kodex.Model)

	if !ok {
		api.HandleError(c, 500, fmt.Errorf("invalid object"))
		return nil, false
	}

	return roleObj, true

}
