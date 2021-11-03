// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - Germany
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

package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"github.com/kiprotect/kodex/api/decorators"
	"github.com/kiprotect/kodex/api/v1/resources"
	"github.com/kiprotect/kodex/api/v1/resources/admin"
	"github.com/kiprotect/kodex/api/v1/resources/pcap"
)

func Initialize(group *gin.RouterGroup,
	controller api.Controller, meter kodex.Meter) error {
	/*
	   Here we define the routes for the V1 of the API
	*/

	// Version 1
	endpoints := group.Group("/v1")
	settings := controller.Settings()

	// All endpoints require a valid user
	endpoints.Use(decorators.ValidUser(settings, []string{"kiprotect:api*"}, false))

	if meter != nil {
		decorators.BaseRateLimitSetup(settings, endpoints, meter)
		endpoints.GET("/usage", admin.UsageEndpoint(meter, "organizationMeterId"))
	}

	// Hello, world endpoint
	endpoints.GET("/hello", resources.SayHello)

	// User endpoint

	userEndpoints := endpoints.Group("")
	userEndpoints.Use(decorators.ValidUser(settings, []string{"kiprotect:api:user"}, false))
	userEndpoints.GET("/user", resources.User)

	transformEndpoints := endpoints.Group("")
	transformEndpoints.Use(decorators.ValidUser(settings, []string{"kiprotect:api:transform"}, false))
	// Data Transformatin Endpoints
	transformEndpoints.POST("/transform", resources.TransformEndpoint(meter))

	// Data Protection Endpoints
	transformEndpoints.POST("/protect/pcap", pcap.Protect)
	transformEndpoints.POST("/unprotect/pcap", pcap.Unprotect)

	// Definitions for readers, writers and actions
	definitionsEndpoints := endpoints.Group("")
	definitionsEndpoints.Use(decorators.ValidUser(settings, []string{"kiprotect:api:definitions"}, false))
	definitionsEndpoints.GET("/definitions", resources.Definitions)

	// Item Submission Endpoint
	submitEndpoints := endpoints.Group("")
	submitEndpoints.Use(decorators.ValidObject(settings,
		"stream", []string{"admin", "superuser", "writer"}, []string{"kiprotect:api:stream:submit"}))
	submitEndpoints.POST("/submit/:streamID", resources.Submit)

	transformConfigEndpoints := endpoints.Group("")
	transformConfigEndpoints.Use(decorators.ValidObject(settings,
		"config", []string{"admin", "superuser", "writer"}, []string{"kiprotect:api:config:transform"}))
	transformConfigEndpoints.POST("/configs/:configID/transform", resources.TransformConfigEndpoint(meter))

	transformActionEndpoints := endpoints.Group("")
	transformActionEndpoints.Use(decorators.ValidObject(settings,
		"action", []string{"admin", "superuser", "writer"}, []string{"kiprotect:api:action:transform"}))
	transformActionEndpoints.POST("/actions/:actionID/transform", resources.TransformActionConfigEndpoint(meter))

	for _, objectAdaptor := range controller.APIDefinitions().ObjectAdaptors {

		objectType := objectAdaptor.Type()
		dependsOn := objectAdaptor.DependsOn()

		// Object Management
		objectEndpoints := endpoints.Group("")

		if err := objectAdaptor.Initialize(controller, objectEndpoints); err != nil {
			return err
		}

		objectEndpoints.Use(decorators.ObjectType(objectType))

		objectBase := "/"

		objectListEndpoints := objectEndpoints.Group("")
		// All endpoints require a valid user
		objectListEndpoints.Use(decorators.ValidUser(settings, []string{fmt.Sprintf("kiprotect:api:%s:read", objectType)}, false))

		if dependsOn != "" {
			objectBase = fmt.Sprintf("/%ss/:%sID/",
				dependsOn, dependsOn)
			// we require the object this resources depend on
			objectListEndpoints.Use(decorators.ValidObject(settings,
				dependsOn, []string{"admin", "superuser", "viewer"}, []string{}))

		}

		// every user can view objects for which he/she has one role
		objectListEndpoints.GET(fmt.Sprintf("%s%ss",
			objectBase, objectType), resources.Objects)

		if _, ok := objectAdaptor.(api.ListAllObjectAdaptor); ok {
			objectListAllEndpoints := objectEndpoints.Group("")
			// All endpoints require a valid user
			objectListEndpoints.Use(decorators.ValidUser(settings, []string{fmt.Sprintf("kiprotect:api:%s:read", objectType)}, false))
			objectListAllEndpoints.GET(fmt.Sprintf("%ss", objectType), resources.AllObjects)
		}

		objectDetailsEndpoints := objectEndpoints.Group("")
		objectDetailsEndpoints.Use(decorators.ValidObject(settings,
			objectType, []string{"superuser", "admin", "viewer"}, []string{fmt.Sprintf("kiprotect:api:%s:read", objectType)}))

		// we expose object statistics via the usage endpoint
		if modelMeter, ok := meter.(kodex.ModelMeter); ok {
			usageEndpoint := admin.UsageEndpoint(meter, "objMeterId")
			objectStatsEndpoints := objectEndpoints.Group("")
			objectStatsEndpoints.Use(decorators.ValidObject(settings,
				objectType, []string{"superuser", "admin", "viewer"}, []string{fmt.Sprintf("kiprotect:api:%s:stats", objectType), fmt.Sprintf("kiprotect:api:%s:read", objectType)}))
			objectStatsEndpoints.GET(fmt.Sprintf("/%ss/:%sID/stats", objectType,
				objectType), func(c *gin.Context) {

				obj := resources.GetObj(c, "objectType")

				if obj == nil {
					return
				}

				meterId := modelMeter.ModelID(obj)
				c.Set("objMeterId", meterId)

				usageEndpoint(c)
			})
		}

		// object details endpoint
		objectDetailsEndpoints.GET(fmt.Sprintf("/%ss/:%sID",
			objectType, objectType), resources.ObjectDetails)

		// only object superusers can perform advanced object operations
		objectSuperusers := objectEndpoints.Group("")
		objectSuperusers.Use(decorators.ValidObject(settings,
			objectType, []string{"superuser"}, []string{fmt.Sprintf("kiprotect:api:%s:write", objectType)}))
		// delete a object
		objectSuperusers.DELETE(fmt.Sprintf("/%ss/:%sID", objectType,
			objectType), resources.DeleteObject)
		// update a object
		objectSuperusers.PATCH(fmt.Sprintf("/%ss/:%sID", objectType,
			objectType), resources.UpdateObject)

		if _, ok := objectAdaptor.(api.CreateObjectAdaptor); ok {

			// only admins or superusers can create a new object
			newObjectEndpoints := objectEndpoints.Group("")

			if dependsOn == "" {

				newObjectEndpoints.Use(decorators.ValidUser(settings, []string{fmt.Sprintf("kiprotect:api:%s:create", objectType)}, false))

				newObjectEndpoints.Use(decorators.ValidOrganization([]string{"admin",
					"superuser"}))
				// create a object
				newObjectEndpoints.POST(fmt.Sprintf("/orgs/:organizationID%s%ss",
					objectBase, objectType), resources.CreateObject)
			} else {

				newObjectEndpoints.Use(decorators.ValidObject(settings,
					dependsOn, []string{"admin", "superuser"}, []string{fmt.Sprintf("kiprotect:api:%s:create", objectType)}))

				newObjectEndpoints.Use(decorators.ValidObject(settings,
					dependsOn, []string{"admin", "superuser"}, []string{}))

				// create a object
				newObjectEndpoints.POST(fmt.Sprintf("%s%ss",
					objectBase, objectType), resources.CreateObject)
			}
		}

		if dependsOn == "" {
			objectRoleEndpoints := objectEndpoints.Group("")
			objectRoleEndpoints.Use(decorators.ValidObject(settings,
				objectType, []string{"superuser", "admin"}, []string{fmt.Sprintf("kiprotect:api:%s:roles", objectType)}))
			// get object roles
			objectRoleEndpoints.GET(fmt.Sprintf("/%ss/:%sID/roles",
				objectType, objectType), resources.ObjectRoles)
			// delete a object role
			objectRoleEndpoints.DELETE(fmt.Sprintf("/%ss/:%sID/roles/:roleID",
				objectType, objectType),
				resources.DeleteObjectRole)
			newObjectSuperusers := objectRoleEndpoints.Group("")
			// the user must have at least one role in the organization
			newObjectSuperusers.Use(decorators.ValidOrganization([]string{"*"}))
			// create a new object role
			newObjectSuperusers.POST(fmt.Sprintf("/orgs/:organizationID/%ss/:%sID/roles",
				objectType, objectType),
				resources.CreateObjectRole)
		}
	}

	for _, associateAdaptor := range controller.APIDefinitions().AssociateAdaptors {

		left := associateAdaptor.LeftType()
		right := associateAdaptor.RightType()

		// Object Management
		associateEndpoints := endpoints.Group("")

		associateEndpoints.Use(decorators.AssociateType(left, right))

		associateUrl := fmt.Sprintf("/%ss/:%sID/%ss/:%sID", left, left, right, right)

		getUrl := fmt.Sprintf("/%ss/:%sID/%ss", left, left, right)

		getEndpoints := associateEndpoints.Group("")

		associateEndpoints.Use(decorators.ValidObject(settings,
			left, []string{"admin", "superuser"}, []string{fmt.Sprintf("kiprotect:api:%s:write", left)}))

		associateEndpoints.Use(decorators.ValidObject(settings,
			right, []string{"admin", "superuser"}, []string{fmt.Sprintf("kiprotect:api:%s:write", right)}))

		getEndpoints.Use(decorators.ValidObject(settings,
			left, []string{"admin", "superuser", "viewer"}, []string{fmt.Sprintf("kiprotect:api:%s:read", left)}))

		getEndpoints.GET(getUrl, resources.AssociatedObjects)

		associateEndpoints.DELETE(associateUrl,
			resources.DissociateObjects)

		associateEndpoints.POST(associateUrl,
			resources.AssociateObjects)

	}

	return nil

}
