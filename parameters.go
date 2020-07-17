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

package kiprotect

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/go-helpers/maps"
	"sort"
)

type ParameterStoreDefinition struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Internal    bool                `json:"internal"`
	Maker       ParameterStoreMaker `json:"-"`
	Form        forms.Form          `json:"form"`
}

type ParameterStoreDefinitions map[string]ParameterStoreDefinition

type ParameterStoreMaker func(map[string]interface{}, Definitions) (ParameterStore, error)

// An interface that manages action parameters
type ParameterStore interface {

	// Returns parameters for a given action config (or nil if none exist). If
	// there is an error (e.g. parameters exist for the same action config
	// as identified by their ID but they are incompatible) and error will
	// be returned.
	ParametersById(id []byte) (*Parameters, error)
	Parameters(action Action, parameterGroup *ParameterGroup) (*Parameters, error)
	ParameterSet(hash []byte) (*ParameterSet, error)
	SaveParameterSet(*ParameterSet) (bool, error)
	SaveParameters(*Parameters) (bool, error)
}

func MakeParameterStore(settings Settings, definitions Definitions) (ParameterStore, error) {
	config, err := settings.Get("parameter-store")

	if err != nil {
		config = map[string]interface{}{
			"type": "in-memory",
		}
	}

	configMap, ok := maps.ToStringMap(config)

	if !ok {
		return nil, fmt.Errorf("not a valid config for parameter store")
	}

	storeType, ok := configMap["type"].(string)

	if !ok {
		return nil, fmt.Errorf("type is missing")
	}

	definition, ok := definitions.ParameterStoreDefinitions[storeType]

	if !ok {
		return nil, fmt.Errorf("not a valid parameter store type: %s", storeType)
	}

	return definition.Maker(configMap, definitions)
}

type ParameterSet struct {
	parameterStore ParameterStore
	parameters     []*Parameters
	hash           []byte
}

func MakeParameterSet(actions []Action, parameterStore ParameterStore) (*ParameterSet, error) {
	parameters := make([]*Parameters, len(actions))

	for i, action := range actions {
		parameters[i] = MakeParameters(action, parameterStore, nil, nil)
	}

	ps := &ParameterSet{
		parameterStore: parameterStore,
		parameters:     parameters,
		hash:           []byte("default"),
	}

	return ps, ps.UpdateHash()
}

func (a ParameterSet) MarshalJSON() ([]byte, error) {

	parametersList := make([]string, len(a.parameters))

	for i, parameters := range a.parameters {
		parametersList[i] = hex.EncodeToString(parameters.ID())
	}

	return json.Marshal(map[string]interface{}{
		"parameters": parametersList,
		"hash":       hex.EncodeToString(a.hash),
	})
}

func (p *ParameterSet) SetParameterStore(parameterStore ParameterStore) {
	p.parameterStore = parameterStore
}

func (p *ParameterSet) ParameterStore() ParameterStore {
	return p.parameterStore
}

func (p *ParameterSet) Actions() []Action {
	actions := make([]Action, len(p.parameters))
	for i, parameters := range p.parameters {
		actions[i] = parameters.Action()
	}
	return actions
}

func (p *ParameterSet) ParametersFor(action Action, parameterGroup *ParameterGroup) (*Parameters, bool, error) {
	if p.parameterStore == nil {
		return nil, false, fmt.Errorf("no parameter store given")
	}
	found := false
	var i = 0
	var parameters *Parameters
	for i, parameters = range p.parameters {
		if bytes.Equal(parameters.Action().ID(), action.ID()) {
			if valid, err := parameters.Valid(action, parameterGroup); err == nil && valid {
				return parameters, false, nil
			} else if err != nil {
				return nil, false, err
			}
			found = true
			break
		}
	}
	if !found {
		// the action config does not exist, this should not happen
		return nil, false, fmt.Errorf("action not found")
	}
	// we try to load valid parameters from the store
	validParameters, err := p.parameterStore.Parameters(action, parameterGroup)
	if err != nil {
		return nil, false, err
	}
	// if we have found some valid parameters in the store we replace them in
	// the parameter set
	if validParameters != nil {
		p.parameters[i] = validParameters
		return validParameters, true, p.UpdateHash()
	}
	// we haven't found any valid parameters for this group
	return nil, false, nil
}

func (p *ParameterSet) Hash() []byte {
	return p.hash
}

func (p *ParameterSet) UpdateParameters(action Action, params interface{}, parameterGroup *ParameterGroup) error {

	found := false
	var i int
	var parameters *Parameters
	for i, parameters = range p.parameters {
		if bytes.Equal(parameters.Action().ID(), action.ID()) {
			found = true
			break
		}
	}
	// we haven't found the action config (this should never happen)
	if !found {
		return fmt.Errorf("action config not found")
	}

	// we try to generate new parameters for the given parameter group
	newParameters := MakeParameters(action, p.ParameterStore(), params, parameterGroup)
	if err := newParameters.Save(); err != nil {
		return err
	}

	// we replace the parameters for the given action
	p.parameters[i] = newParameters

	return p.UpdateHash()
}

// The hash uniquely identifies a given parameters set based on the IDs of
// the constitutent paremeters.
func (p *ParameterSet) UpdateHash() error {
	ids := make([]string, 0)
	for _, parameters := range p.parameters {
		id := parameters.ID()
		if id == nil {
			return fmt.Errorf("parameters ID is null")
		}
		ids = append(ids, string(id))
	}
	// we sort the parameter IDs
	sort.Strings(ids)
	if hash, err := StructuredHash(ids); err != nil {
		return err
	} else {
		p.hash = hash
	}
	return nil
}

// Returns the parameters (in order) for the set.
func (p *ParameterSet) Parameters() []*Parameters {
	return p.parameters
}

func (p *ParameterSet) Empty() bool {
	for _, parameters := range p.parameters {
		if parameters.Parameters() != nil {
			return false
		}
	}
	return true
}

// Saves the parameter set
func (p *ParameterSet) Save() error {
	if p.Empty() {
		return nil
	}
	if p.parameterStore == nil {
		return fmt.Errorf("no parameter store given")
	}
	_, err := p.parameterStore.SaveParameterSet(p)
	return err
}

type ParameterGroup struct {
	data map[string]interface{} `json:"data"`
	hash []byte                 `json:"hash"`
}

func (p *ParameterGroup) Hash() []byte {
	return p.hash
}

func (p *ParameterGroup) Data() map[string]interface{} {
	return p.data
}

// Represents parameters for a given ActionConfig
type Parameters struct {
	parameterStore ParameterStore
	parameters     interface{}
	action         Action
	id             []byte
	parameterGroup *ParameterGroup
}

var ParameterForm = forms.Form{
	ErrorMsg: "invalid data encountered in the parameter form",
	Fields: []forms.Field{
		{
			Name: "action",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsStringMap{
					Form: &forms.Form{
						Fields: []forms.Field{
							{
								Name: "type",
								Validators: []forms.Validator{
									forms.IsRequired{},
									forms.IsString{},
								},
							},
							{
								Name: "id",
								Validators: []forms.Validator{
									forms.IsRequired{},
									forms.IsBytes{Encoding: "hex"},
								},
							},
							{
								Name: "name",
								Validators: []forms.Validator{
									forms.IsRequired{},
									forms.IsString{},
								},
							},
							{
								Name: "description",
								Validators: []forms.Validator{
									forms.IsRequired{},
									forms.IsString{},
								},
							},
							{
								Name: "config",
								Validators: []forms.Validator{
									forms.IsRequired{},
									forms.IsStringMap{},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "parameter_group",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsStringMap{
					Form: &forms.Form{
						Fields: []forms.Field{
							{
								Name: "data",
								Validators: []forms.Validator{
									forms.IsRequired{},
								},
							},
							{
								Name: "hash",
								Validators: []forms.Validator{
									forms.IsRequired{},
									forms.IsBytes{
										Encoding: "hex",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "id",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsBytes{Encoding: "hex"},
			},
		},
		{
			Name: "parameters",
			Validators: []forms.Validator{
				forms.IsOptional{Default: nil},
			},
		},
	},
}

var ParameterSetForm = forms.Form{
	ErrorMsg: "invalid data encountered in the parameter set form",
	Fields: []forms.Field{
		{
			Name: "parameters",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsBytes{
							Encoding: "hex",
						},
					},
				},
			},
		},
		{
			Name: "hash",
			Validators: []forms.Validator{
				forms.IsRequired{},
				forms.IsBytes{
					Encoding: "hex",
				},
			},
		},
	},
}

func RestoreParameterSet(data map[string]interface{}, parameterStore ParameterStore) (*ParameterSet, error) {
	config, err := ParameterSetForm.Validate(data)
	if err != nil {
		return nil, err
	}
	paramsIds := config["parameters"].([]interface{})
	parametersList := make([]*Parameters, len(paramsIds))
	for i, id := range paramsIds {
		parameters, err := parameterStore.ParametersById(id.([]byte))
		if err != nil {
			return nil, err
		}
		if parameters == nil {
			return nil, fmt.Errorf("parameters not found")
		}
		parametersList[i] = parameters
	}
	return &ParameterSet{
		parameterStore: parameterStore,
		parameters:     parametersList,
		hash:           config["hash"].([]byte),
	}, nil
}

func RestoreParameters(data map[string]interface{}, parameterStore ParameterStore, definitions Definitions) (*Parameters, error) {
	config, err := ParameterForm.Validate(data)
	if err != nil {
		return nil, err
	}
	actionConfig := config["action"].(map[string]interface{})
	action, err := MakeAction(actionConfig["name"].(string), actionConfig["description"].(string), actionConfig["type"].(string), actionConfig["id"].([]byte), actionConfig["config"].(map[string]interface{}), definitions)
	if err != nil {
		return nil, err
	}
	parameterGroupConfig := config["parameter_group"].(map[string]interface{})
	return &Parameters{
		action:         action,
		id:             config["id"].([]byte),
		parameters:     config["parameters"],
		parameterStore: parameterStore,
		parameterGroup: &ParameterGroup{
			data: parameterGroupConfig["data"].(map[string]interface{}),
			hash: parameterGroupConfig["hash"].([]byte),
		},
	}, nil
}

func (a Parameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"parameters": a.parameters,
		"id":         hex.EncodeToString(a.id),
		"action":     a.action,
		"parameter_group": map[string]interface{}{
			"hash": hex.EncodeToString(a.parameterGroup.Hash()),
			"data": a.parameterGroup.Data(),
		},
	})
}

func MakeParameters(action Action, parameterStore ParameterStore, parameters interface{}, parameterGroup *ParameterGroup) *Parameters {
	return &Parameters{
		id:             RandomID(),
		action:         action,
		parameterStore: parameterStore,
		parameters:     parameters,
		parameterGroup: parameterGroup,
	}
}

// Saves the parameter set
func (p *Parameters) Save() error {
	if p.parameterStore == nil {
		return fmt.Errorf("no parameter store given")
	}
	_, err := p.parameterStore.SaveParameters(p)
	return err
}

func (p *Parameters) SetParameterStore(parameterStore ParameterStore) {
	p.parameterStore = parameterStore
}

func (p *Parameters) ParameterStore() ParameterStore {
	return p.parameterStore
}

func (p *Parameters) ParameterGroup() *ParameterGroup {
	return p.parameterGroup
}

// Returns the parameters
func (p *Parameters) Parameters() interface{} {
	return p.parameters
}

// Returns whether the parameters are valid for a given parameter group
func (p *Parameters) Valid(action Action, parameterGroup *ParameterGroup) (bool, error) {
	configHashA, err := action.ConfigHash()
	if err != nil {
		return false, err
	}
	configHashB, err := p.action.ConfigHash()
	if err != nil {
		return false, err
	}
	return bytes.Equal(action.ID(), p.Action().ID()) && bytes.Equal(configHashA, configHashB) && p.parameterGroup != nil && bytes.Equal(p.parameterGroup.Hash(), parameterGroup.Hash()), nil
}

// Returns the associated action config
func (p *Parameters) Action() Action {
	return p.action
}

// Returns the ID
func (p *Parameters) ID() []byte {
	return p.id
}
