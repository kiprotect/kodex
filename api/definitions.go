// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
)

type Definitions struct {
	kodex.Definitions
	APIControllerDefinitions APIControllerDefinitions
	Routes                   []Routes
	ObjectAdaptors           map[string]ObjectAdaptor
	AssociateAdaptors        map[string]AssociateAdaptor
}

func (d Definitions) Marshal() map[string]interface{} {
	objectAdaptors := make(map[string]map[string]interface{})
	associateAdaptors := make(map[string]map[string]interface{})
	for name, oa := range d.ObjectAdaptors {
		objectAdaptors[name] = map[string]interface{}{
			"type":       oa.Type(),
			"depends-on": oa.DependsOn(),
		}
	}
	return map[string]interface{}{
		"objects":      objectAdaptors,
		"associations": associateAdaptors,
	}
}

func (d Definitions) MarshalJSON() ([]byte, error) {
	ed := d.Marshal()
	dd := d.Definitions.Marshal()
	for k, v := range dd {
		ed[k] = v
	}
	return json.Marshal(ed)
}

func MergeDefinitions(a, b Definitions) Definitions {
	cc := kodex.MergeDefinitions(a.Definitions, b.Definitions)
	c := Definitions{
		Definitions:              cc,
		Routes:                   []Routes{},
		APIControllerDefinitions: map[string]APIControllerMaker{},
		ObjectAdaptors:           map[string]ObjectAdaptor{},
		AssociateAdaptors:        map[string]AssociateAdaptor{},
	}
	for _, obj := range []Definitions{a, b} {
		for _, route := range obj.Routes {
			c.Routes = append(c.Routes, route)
		}
		for k, v := range obj.APIControllerDefinitions {
			c.APIControllerDefinitions[k] = v
		}
		for k, v := range obj.ObjectAdaptors {
			c.ObjectAdaptors[k] = v
		}
		for k, v := range obj.AssociateAdaptors {
			c.AssociateAdaptors[k] = v
		}
	}
	return c
}

type ObjectAdaptor interface {
	// Returns the object for the given ID and its role object, if possible
	Get(Controller, *gin.Context, []byte) (kodex.Model, kodex.Model, error)
	Initialize(Controller, *gin.RouterGroup) error
	Objects(*gin.Context) []kodex.Model
	Type() string
	DependsOn() string
}

type ListAllObjectAdaptor interface {
	AllObjects(*gin.Context) []kodex.Model
}

type CreateObjectAdaptor interface {
	MakeObject(*gin.Context) kodex.Model
}

type UpdateObjectAdaptor interface {
	UpdateObject(kodex.Model, map[string]interface{}) (kodex.Model, error)
	SaveUpdated(updatedObject, object kodex.Model) error
}

type AssociateAdaptor interface {
	Associate(c *gin.Context, left, right kodex.Model) bool
	Dissociate(c *gin.Context, left, right kodex.Model) bool
	Get(c *gin.Context, left kodex.Model) interface{}
	LeftType() string
	RightType() string
}

type Routes func(*gin.RouterGroup, Controller, kodex.Meter) error
