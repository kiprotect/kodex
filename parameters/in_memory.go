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

package parameters

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/kiprotect/kiprotect"
	"sync"
)

type InMemoryParameterStore struct {
	mutex         sync.Mutex
	definitions   kiprotect.Definitions
	config        map[string]interface{}
	parameterSets map[string]*kiprotect.ParameterSet
	// stores parameters based on the action ID
	parameters     map[string]map[string][]*kiprotect.Parameters
	parametersById map[string]*kiprotect.Parameters
}

func MakeInMemoryParameterStore(config map[string]interface{}, definitions kiprotect.Definitions) (kiprotect.ParameterStore, error) {
	return &InMemoryParameterStore{
		config:         config,
		definitions:    definitions,
		mutex:          sync.Mutex{},
		parameterSets:  make(map[string]*kiprotect.ParameterSet),
		parameters:     make(map[string]map[string][]*kiprotect.Parameters),
		parametersById: make(map[string]*kiprotect.Parameters),
	}, nil
}

func (p *InMemoryParameterStore) ParametersById(id []byte) (*kiprotect.Parameters, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	parameters, _ := p.parametersById[hex.EncodeToString(id)]
	return parameters, nil
}

func (p *InMemoryParameterStore) RestoreParameters(data map[string]interface{}) (*kiprotect.Parameters, error) {
	return kiprotect.RestoreParameters(data, p, p.definitions)
}

func (p *InMemoryParameterStore) RestoreParameterSet(data map[string]interface{}) (*kiprotect.ParameterSet, error) {
	return kiprotect.RestoreParameterSet(data, p)
}

func (p *InMemoryParameterStore) Parameters(action kiprotect.Action, parameterGroup *kiprotect.ParameterGroup) (*kiprotect.Parameters, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.getParameters(action, parameterGroup)
}

func (p *InMemoryParameterStore) getParameters(action kiprotect.Action, parameterGroup *kiprotect.ParameterGroup) (*kiprotect.Parameters, error) {
	if action.ID() == nil {
		return nil, fmt.Errorf("action has no ID")
	}
	id := hex.EncodeToString(action.ID())
	actionParameters, ok := p.parameters[id]
	if !ok {
		return nil, nil
	}
	configHash, err := action.ConfigHash()
	if err != nil {
		return nil, err
	}
	configHashStr := hex.EncodeToString(configHash)
	for actionConfigHash, configHashParameters := range actionParameters {
		if actionConfigHash == configHashStr {
			for _, parameters := range configHashParameters {
				if bytes.Equal(parameters.ParameterGroup().Hash(), parameterGroup.Hash()) {
					return parameters, nil
				}
			}
		}
	}
	return nil, nil
}

func (p *InMemoryParameterStore) AllParameters() ([]*kiprotect.Parameters, error) {
	parametersList := make([]*kiprotect.Parameters, 0, len(p.parametersById))
	for _, parameters := range p.parametersById {
		parametersList = append(parametersList, parameters)
	}
	return parametersList, nil
}

func (p *InMemoryParameterStore) AllParameterSets() ([]*kiprotect.ParameterSet, error) {
	parameterSets := make([]*kiprotect.ParameterSet, 0, len(p.parameterSets))
	for _, parameterSet := range p.parameterSets {
		parameterSets = append(parameterSets, parameterSet)
	}
	return parameterSets, nil
}

func (p *InMemoryParameterStore) DeleteParameterSet(parameterSet *kiprotect.ParameterSet) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.parameterSets, hex.EncodeToString(parameterSet.Hash()))
	return nil
}

func (p *InMemoryParameterStore) DeleteParameters(parameters *kiprotect.Parameters) error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	actionID := parameters.Action().ID()

	if actionID == nil {
		return fmt.Errorf("action config has no ID")
	}
	id := hex.EncodeToString(actionID)
	actionParameters, ok := p.parameters[id]
	if !ok {
		return fmt.Errorf("parameters do not exist")
	}
	configHash, err := parameters.Action().ConfigHash()
	if err != nil {
		return err
	}
	configHashStr := hex.EncodeToString(configHash)
	for actionConfigHash, configHashParameters := range actionParameters {
		if actionConfigHash == configHashStr {
			newActionParameters := make([]*kiprotect.Parameters, 0, len(actionParameters)-1)
			for _, existingParameters := range configHashParameters {
				if bytes.Equal(existingParameters.ParameterGroup().Hash(), parameters.ParameterGroup().Hash()) {
					continue
				}
				newActionParameters = append(newActionParameters, parameters)
			}
			p.parameters[id][configHashStr] = newActionParameters
			break
		}
	}
	delete(p.parametersById, hex.EncodeToString(parameters.ID()))
	return nil

}

func (p *InMemoryParameterStore) SaveParameters(parameters *kiprotect.Parameters) (bool, error) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	// we call the unlocked function as we have created the lock above and want
	// this whole function call to be atomic
	existingParameters, err := p.getParameters(parameters.Action(), parameters.ParameterGroup())

	if err != nil {
		return false, err
	}

	if existingParameters != nil {
		return false, fmt.Errorf("parameters already exist for this parameter group")
	}

	if parameters.Action().ID() == nil {
		return false, fmt.Errorf("action config has no ID")
	}
	id := hex.EncodeToString(parameters.Action().ID())
	configHash, err := parameters.Action().ConfigHash()
	if err != nil {
		return false, err
	}
	configHashStr := hex.EncodeToString(configHash)
	actionParameters, ok := p.parameters[id]

	if !ok {
		actionParameters = make(map[string][]*kiprotect.Parameters)
		p.parameters[id] = actionParameters
	}

	configHashParameters, ok := actionParameters[configHashStr]

	if !ok {
		configHashParameters = make([]*kiprotect.Parameters, 0, 10)
	}

	for _, existingParameters := range configHashParameters {
		if bytes.Equal(existingParameters.ID(), parameters.ID()) {
			return false, nil
		}
	}

	configHashParameters = append(configHashParameters, parameters)
	p.parameters[id][configHashStr] = configHashParameters
	p.parametersById[hex.EncodeToString(parameters.ID())] = parameters

	return true, nil

}

func (p *InMemoryParameterStore) ParameterSet(hash []byte) (*kiprotect.ParameterSet, error) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	parameterSet, ok := p.parameterSets[hex.EncodeToString(hash)]
	if !ok {
		return nil, nil
	}
	return parameterSet, nil
}

func (p *InMemoryParameterStore) SaveParameterSet(parameterSet *kiprotect.ParameterSet) (bool, error) {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	hashStr := hex.EncodeToString(parameterSet.Hash())
	if _, ok := p.parameterSets[hashStr]; ok {
		return false, nil
	}
	p.parameterSets[hashStr] = parameterSet
	return true, nil
}
