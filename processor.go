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
	"encoding/hex"
	"github.com/kiprotect/go-helpers/errors"
)

type Processor struct {
	errorPolicy   ErrorPolicy
	parameterSet  *ParameterSet
	channelWriter ChannelWriter
	config        Config
	key, salt     []byte
	id            string
}

func (p *Processor) ParameterSet() *ParameterSet {
	return p.parameterSet
}

func (p *Processor) updateParams(item *Item, undo bool) error {

	updated := false

	for _, action := range p.parameterSet.Actions() {
		// we get the parameter group for the specific item
		parameterGroup, err := action.ParameterGroup(item)
		if err != nil {
			return err
		}
		i := 0
		for {
			i += 1
			var actionParams *Parameters
			var loaded bool
			var err error
			// if a key is specified, we generate all parameters from it and
			// do not persist anything to the parameter store
			if p.key == nil {
				actionParams, loaded, err = p.parameterSet.ParametersFor(action, parameterGroup)
				if err != nil {
					// this might be a race condition with another processor
					if i > 2 {
						return errors.MakeExternalError("(1) error getting action params", "GET-ACTION-PARAMS", nil, err)
					}
					continue
				}
			}
			if actionParams == nil {
				if undo && p.key == nil {
					return errors.MakeExternalError("error getting parameters", "GET-ACTION-PARAMS", nil, nil)
				}
				// there are no action parameters, we generate some and try
				// to save them
				if err = action.GenerateParams(p.key, p.salt); err != nil {
					return errors.MakeExternalError("error generating params", "GEN-PARAMS", nil, err)
				}
				if p.key == nil && action.HasParams() {
					// we update the action parametrs
					if err := p.parameterSet.UpdateParameters(action, action.Params(), parameterGroup); err != nil {
						// this might be a race condition with another processor
						if i > 2 {
							return errors.MakeExternalError("error setting params", "GEN-PARAMS", nil, err)
						}
						// we try again, maybe another processor was simply faster
						continue
					}
					updated = true
				}
			} else if action.HasParams() {
				// if the parameters were loaded we need to potentially save the parameter set
				if loaded {
					updated = true
				}
				// the action might have changed
				if err = actionParams.Action().SetParams(actionParams.Parameters()); err != nil {
					return errors.MakeExternalError("error setting params", "SET-PARAMS", nil, err)
				}
			}
			break
		}
	}
	if updated && !undo && p.key == nil && !p.parameterSet.Empty() {
		// we try to save the new parameter set as well
		return p.parameterSet.Save()
	}
	return nil
}

func MakeProcessor(parameterSet *ParameterSet, channelWriter ChannelWriter, config Config) (*Processor, error) {

	processor := Processor{
		parameterSet:  parameterSet,
		channelWriter: channelWriter,
		errorPolicy:   ReportErrors,
		config:        config,
		id:            hex.EncodeToString(RandomID()),
	}
	return &processor, nil
}

func (p *Processor) SetWriter(channelWriter ChannelWriter) {
	p.channelWriter = channelWriter
}

func (p *Processor) Writer() ChannelWriter {
	return p.channelWriter
}

func (p *Processor) SetSalt(salt []byte) {
	p.salt = salt
}

func (p *Processor) SetKey(key []byte) {
	p.key = key
}

func (p *Processor) SetErrorPolicy(policy ErrorPolicy) {
	p.errorPolicy = policy
}

func (p *Processor) ErrorPolicy() ErrorPolicy {
	return p.errorPolicy
}

func (p *Processor) Setup() error {
	for i, action := range p.parameterSet.Actions() {
		if err := action.Setup(); err != nil {
			// we tear down the actions that were already set up
			for j, otherAction := range p.parameterSet.Actions() {
				if j >= i {
					break
				}
				if err := otherAction.Teardown(); err != nil {
					Log.Error(err)
				}
			}
			return err
		}
	}
	return nil
}

func (p *Processor) Teardown() error {
	// to do: make a proper error list here
	var lastErr error
	for _, action := range p.parameterSet.Actions() {
		if err := action.Teardown(); err != nil {
			Log.Error(err)
			lastErr = err
		}
	}
	return lastErr
}

func (p *Processor) processItem(item *Item, paramsMap map[string]interface{}, undo bool) (*Item, error) {
	var err error
	if err = p.updateParams(item, undo); err != nil {
		return nil, errors.MakeExternalError("error setting action params", "SET-ACTION-PARAMS", nil, err)
	}
	newItem := item
	for _, action := range p.parameterSet.Actions() {
		err = nil
		if undo {
			if undoableAction, ok := action.(UndoableAction); ok {
				// not all actions that have an Undo function are always
				// undoable (e.g. some pseudonymization methods are one-way)
				if undoableAction.Undoable(newItem) {
					newItem, err = undoableAction.Undo(item, p.channelWriter)
				}
			}
		} else {
			if configurableAction, ok := action.(ConfigurableAction); p.config != nil && ok {
				newItem, err = configurableAction.DoWithConfig(newItem, p.channelWriter, p.config)
			} else if doableAction, ok := action.(DoableAction); ok {
				newItem, err = doableAction.Do(newItem, p.channelWriter)
			}
		}
		if err != nil {
			itemError := errors.MakeExternalError("error processing action", "PROCESS-ACTION", action.Name(), err)
			return nil, itemError
		}
		if newItem == nil {
			break
		}
	}
	if !undo && p.key == nil && !p.parameterSet.Empty() {
		hashStr := hex.EncodeToString(p.parameterSet.Hash())
		if paramsMap != nil {
			if _, ok := paramsMap[hashStr]; !ok {
				paramsMap[hashStr] = p.parameterSet
			}
		}
		if newItem != nil {
			newItem.Set("_kip", hashStr)
		}
	}
	if undo {
		newItem.Delete("_kip")
	}
	return newItem, nil
}
func (p *Processor) Undo(items []*Item, paramsMap map[string]interface{}) ([]*Item, error) {
	return p.process(items, paramsMap, true)
}

func (p *Processor) Process(items []*Item, paramsMap map[string]interface{}) ([]*Item, error) {
	return p.process(items, paramsMap, false)
}

func (p *Processor) process(items []*Item, paramsMap map[string]interface{}, undo bool) ([]*Item, error) {
	Log.Debugf("Processing %d items with error policy '%s'", len(items), p.errorPolicy)
	newItems := make([]*Item, 0)
	for _, item := range items {
		newItem, err := p.processItem(item, paramsMap, undo)
		if err != nil {
			switch p.errorPolicy {
			case ReportErrors:
				itemError := errors.MakeExternalError("error processing item", "PROCESS-ITEM", nil, err)
				// if we encounter an error during error reporting (i.e.
				// too many errors received) we abort the processing.
				Log.Error(itemError)
				if err := p.channelWriter.Error(item, itemError); err != nil {
					return newItems, err
				}
				continue
			case AbortOnError:
				return newItems, errors.MakeExternalError("error processing item", "PROCESS-ITEM", map[string]interface{}{"item": item.All()}, err)
			}
		}
		if newItem != nil {
			newItems = append(newItems, newItem)
		}
	}
	return newItems, nil
}
