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

package kodex

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

	for _, action := range p.parameterSet.Actions() {
		// we get the parameter group for the specific item
		parameterGroup, err := action.ParameterGroup(item)
		if err != nil {
			return err
		}
		i := 0
		for {
			i += 1
			var spec *Parameters
			var loaded bool
			var err error
			// if a key is specified, we generate all parameters from it and
			// do not persist anything to the parameter store
			if p.key == nil {
				spec, loaded, err = p.parameterSet.ParametersFor(action, parameterGroup)
				if err != nil {
					// this might be a race condition with another processor
					if i > 2 {
						return errors.MakeExternalError("(1) error getting action params", "GET-ACTION-PARAMS", nil, err)
					}
					continue
				}
			}
			if spec == nil {
				if undo && p.key == nil {
					return errors.MakeExternalError("error getting parameters", "GET-ACTION-PARAMS", nil, nil)
				}
				// there are no action parameters, we generate some and try
				// to save them
				if err = action.GenerateParams(p.key, p.salt); err != nil {
					return errors.MakeExternalError("error generating params", "GEN-PARAMS", nil, err)
				}
				if p.key == nil && action.HasParams() {
					// we update the action parameters
					if err := p.parameterSet.UpdateParameters(action, action.Params(), parameterGroup); err != nil {
						// this might be a race condition with another processor
						if i > 2 {
							return errors.MakeExternalError("error setting params", "GEN-PARAMS", nil, err)
						}
						// we try again, maybe another processor was simply faster
						continue
					}
					if !p.parameterSet.Empty() {
						if err := p.parameterSet.Save(); err != nil {
							return errors.MakeExternalError("error saving parameter set", "SAVE-PARAMS", nil, err)
						}
					}
				}
			} else if action.HasParams() {
				// if the parameters were loaded we need to potentially save the parameter set
				if loaded {
					if err := p.parameterSet.Save(); err != nil {
						return errors.MakeExternalError("error saving parameter set", "SAVE-PARAMS", nil, err)
					}
				}
				// the action might have changed
				if err = spec.Action().SetParams(spec.Parameters()); err != nil {
					return errors.MakeExternalError("error setting params", "SET-PARAMS", nil, err)
				}
			}
			break
		}
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
		if setupAction, ok := action.(SetupAction); ok {
			controller := p.config.Stream().Project().Controller()
			if err := setupAction.Setup(controller.Settings()); err != nil {
				// we tear down the actions that were already set up
				for j, otherAction := range p.parameterSet.Actions() {
					if j >= i {
						break
					}
					if teardownAction, ok := otherAction.(TeardownAction); ok {
						if err := teardownAction.Teardown(); err != nil {
							Log.Error(err)
						}
					}
				}
				return err
			}

		}
	}
	return nil
}

func (p *Processor) Teardown() error {
	// to do: make a proper error list here
	var lastErr error
	for _, action := range p.parameterSet.Actions() {
		if teardownAction, ok := action.(TeardownAction); ok {
			if err := teardownAction.Teardown(); err != nil {
				Log.Error(err)
				lastErr = err
			}
		}
	}
	return lastErr
}

func (p *Processor) Reset() error {
	for _, action := range p.parameterSet.Actions() {
		if statefulAction, ok := action.(StatefulAction); ok {
			if err := statefulAction.Reset(); err != nil {
				switch p.errorPolicy {
				case ReportErrors:
					actionError := errors.MakeExternalError("error resetting action", "RESET-ACTION", nil, err)
					Log.Error(actionError)
					if err := p.channelWriter.Error(nil, actionError); err != nil {
						return err
					}
					continue
				case AbortOnError:
					return errors.MakeExternalError("error resetting action", "RESET-ACTION", map[string]interface{}{"action": action.Name()}, err)
				}
			}
		}
	}
	return nil

}

func (p *Processor) Finalize() ([]*Item, error) {
	finalizedItems := make([]*Item, 0)
	for _, action := range p.parameterSet.Actions() {
		if statefulAction, ok := action.(StatefulAction); ok {
			if newItems, err := statefulAction.Finalize(p.channelWriter); err != nil {
				switch p.errorPolicy {
				case ReportErrors:
					actionError := errors.MakeExternalError("error finalization action", "FINALIZE-ACTION", nil, err)
					// if we encounter an error during error reporting (i.e.
					// too many errors received) we abort the processing.
					Log.Error(actionError)
					if err := p.channelWriter.Error(nil, actionError); err != nil {
						return newItems, err
					}
					continue
				case AbortOnError:
					return finalizedItems, errors.MakeExternalError("error finalizing action", "FINALIZE-ACTION", map[string]interface{}{"action": action.Name()}, err)
				}
			} else {
				finalizedItems = append(finalizedItems, newItems...)
			}
		}
	}
	return finalizedItems, nil
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
					newItem, err = undoableAction.Undo(newItem, p.channelWriter)
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

func (p *Processor) Advance() ([]*Item, error) {
	newItems := make([]*Item, 0)
	for _, action := range p.parameterSet.Actions() {
		if statefulAction, ok := action.(StatefulAction); ok {
			if advanceItems, err := statefulAction.Advance(p.channelWriter); err != nil {
				switch p.errorPolicy {
				case ReportErrors:
					advanceErr := errors.MakeExternalError("error advancing action", "ADVANCE-ACTION", action.Name(), err)
					Log.Error(advanceErr)
					if err := p.channelWriter.Error(nil, advanceErr); err != nil {
						return newItems, err
					}
					continue
				case AbortOnError:
					return newItems, errors.MakeExternalError("error advancing action", "ADVANCE-ACTION", action.Name(), err)
				}
			} else {
				newItems = append(newItems, advanceItems...)
			}
		}
	}
	return newItems, nil
}

func (p *Processor) process(items []*Item, paramsMap map[string]interface{}, undo bool) ([]*Item, error) {
	Log.Tracef("Processing %d items with error policy '%s'", len(items), p.errorPolicy)
	newItems := make([]*Item, 0)
	// we first perform the Advance() method (for stateful actions)
	if advanceItems, err := p.Advance(); err != nil {
		advanceErr := errors.MakeExternalError("error advancing actions", "ADVANCE-ACTIONS", nil, err)
		Log.Error(advanceErr)
		if err := p.channelWriter.Error(nil, advanceErr); err != nil {
			return newItems, err
		}
	} else {
		newItems = append(newItems, advanceItems...)
	}
	for _, item := range items {
		newItem, err := p.processItem(item, paramsMap, undo)
		if err != nil {
			switch p.errorPolicy {
			case ReportErrors:
				itemError := errors.MakeExternalError("error processing item", "PROCESS-ITEM", nil, err)
				// if we encounter an error during error reporting (i.e.
				// too many errors received) we abort the processing.
				if p.config != nil {
					Log.Errorf("Error processing item in config %s (project %s): %v", hex.EncodeToString(p.config.ID()), hex.EncodeToString(p.config.Stream().Project().ID()), itemError)
				} else {
					Log.Errorf("Error processing item: %v", itemError)
				}
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
