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

package fixtures

import (
	"fmt"
	"github.com/kiprotect/kodex"
)

type Project struct {
	Name string
}

func (c Project) Setup(fixtures map[string]interface{}) (interface{}, error) {
	controller, err := GetController(fixtures)
	if err != nil {
		return nil, err
	}

	project := controller.MakeProject()

	values := map[string]interface{}{
		"name": c.Name,
	}

	if err := project.Create(values); err != nil {
		return nil, err
	}

	return project, project.Save()
}

func (c Project) Teardown(fixture interface{}) error {
	return nil
}

type Stream struct {
	Name    string
	Project string
}

func (c Stream) Setup(fixtures map[string]interface{}) (interface{}, error) {

	project, ok := fixtures[c.Project].(kodex.Project)

	if !ok {
		return nil, fmt.Errorf("project missing")
	}

	stream := project.MakeStream()

	values := map[string]interface{}{
		"name":        c.Name,
		"status":      string(kodex.ActiveStream),
		"description": "",
	}

	if err := stream.Create(values); err != nil {
		return nil, err
	}

	return stream, stream.Save()

}

func (c Stream) Teardown(fixture interface{}) error {
	return nil
}

type Config struct {
	Stream  string
	Name    string
	Version string
	Source  string
	Status  kodex.ConfigStatus
}

func (c Config) Setup(fixtures map[string]interface{}) (interface{}, error) {

	stream, ok := fixtures[c.Stream].(kodex.Stream)
	if !ok {
		return nil, fmt.Errorf("not a stream")
	}

	config := stream.MakeConfig()

	values := map[string]interface{}{
		"status":  string(c.Status),
		"name":    c.Name,
		"version": c.Version,
		"source":  c.Source,
	}

	if err := config.Create(values); err != nil {
		return nil, err
	} else {
		return config, config.Save()
	}
}

func (c Config) Teardown(fixture interface{}) error {
	return nil
}

type Source struct {
	Name       string
	Project    string
	SourceType string
	Config     map[string]interface{}
}

func (i Source) Setup(fixtures map[string]interface{}) (interface{}, error) {

	project, ok := fixtures[i.Project].(kodex.Project)

	if !ok {
		return nil, fmt.Errorf("project missing")
	}

	source := project.MakeSource()

	values := map[string]interface{}{
		"type":        i.SourceType,
		"name":        i.Name,
		"description": "",
		"config":      i.Config,
	}

	if err := source.Create(values); err != nil {
		return nil, err
	}

	return source, source.Save()
}

func (i Source) Teardown(fixture interface{}) error {
	return nil
}

type Destination struct {
	Name            string
	Project         string
	DestinationType string
	Config          map[string]interface{}
}

func (o Destination) Setup(fixtures map[string]interface{}) (interface{}, error) {

	project, ok := fixtures[o.Project].(kodex.Project)

	if !ok {
		return nil, fmt.Errorf("project missing")
	}

	destination := project.MakeDestination()

	values := map[string]interface{}{
		"type":        o.DestinationType,
		"name":        o.Name,
		"description": "",
		"config":      o.Config,
	}

	if err := destination.Create(values); err != nil {
		return nil, err
	}

	return destination, destination.Save()
}

func (o Destination) Teardown(fixture interface{}) error {
	return nil
}

type DestinationAdder struct {
	Destination string
	Config      string
	Status      string
	Name        string
}

func (o DestinationAdder) Setup(fixtures map[string]interface{}) (interface{}, error) {

	psDestination, ok := fixtures[o.Destination].(kodex.Destination)

	if !ok {
		return nil, fmt.Errorf("destination not found")
	}

	psConfig, ok := fixtures[o.Config].(kodex.Config)

	if !ok {
		return nil, fmt.Errorf("config not found")
	}

	return nil, psConfig.AddDestination(psDestination, o.Name, kodex.DestinationStatus(o.Status))

}

func (o DestinationAdder) Teardown(fixture interface{}) error {
	return nil
}

type SourceAdder struct {
	Source string
	Stream string
	Status string
}

func (i SourceAdder) Setup(fixtures map[string]interface{}) (interface{}, error) {

	psSource, ok := fixtures[i.Source].(kodex.Source)

	if !ok {
		return nil, fmt.Errorf("source not found")
	}

	psStream, ok := fixtures[i.Stream].(kodex.Stream)

	if !ok {
		return nil, fmt.Errorf("stream not found")
	}

	return nil, psStream.AddSource(psSource, kodex.SourceStatus(i.Status))

}

func (i SourceAdder) Teardown(fixture interface{}) error {
	return nil
}

type ActionMap struct {
	Config string
	Action string
	Index  int
}

func (a ActionMap) Setup(fixtures map[string]interface{}) (interface{}, error) {

	config, ok := fixtures[a.Config].(kodex.Config)

	if !ok {
		return nil, fmt.Errorf("action map config missing")
	}

	actionConfig, ok := fixtures[a.Action].(kodex.ActionConfig)

	if !ok {
		return nil, fmt.Errorf("action config missing")
	}

	if err := config.AddActionConfig(actionConfig, a.Index); err != nil {
		return nil, err
	}

	return nil, nil
}

func (a ActionMap) Teardown(fixture interface{}) error {
	return nil
}

type ActionConfig struct {
	Name    string
	Project string
	Type    string
	Config  map[string]interface{}
}

func (a ActionConfig) Setup(fixtures map[string]interface{}) (interface{}, error) {

	project, ok := fixtures[a.Project].(kodex.Project)

	if !ok {
		return nil, fmt.Errorf("project missing")
	}

	actionConfig := project.MakeActionConfig()

	values := map[string]interface{}{
		"name":   a.Name,
		"type":   a.Type,
		"config": a.Config,
	}

	if err := actionConfig.Create(values); err != nil {
		return nil, err
	}

	return actionConfig, actionConfig.Save()
}

func (a ActionConfig) Teardown(fixture interface{}) error {
	return nil
}

func GetController(fixtures map[string]interface{}) (kodex.Controller, error) {

	controllerObj := fixtures["controller"]

	if controllerObj == nil {
		return nil, fmt.Errorf("A controller is required")
	}

	controller, ok := controllerObj.(kodex.Controller)

	if !ok {
		return nil, fmt.Errorf("Controller should be an controller")
	}

	return controller, nil

}
