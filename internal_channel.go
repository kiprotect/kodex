// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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
	"fmt"
	"github.com/kiprotect/go-helpers/maps"
	"sync"
)

type InternalChannel struct {
	InternalReader Reader
	InternalWriter Writer
}

// Adaptor for using an internal channel as a writer
type InternalWriter struct {
	*InternalChannel
}

// Adaptor for using an internal channel as a reader
type InternalReader struct {
	*InternalChannel
}

func MakeInternalWriter(channel *InternalChannel) *InternalWriter {
	return &InternalWriter{
		InternalChannel: channel,
	}
}

func MakeInternalReader(channel *InternalChannel) *InternalReader {
	return &InternalReader{
		InternalChannel: channel,
	}
}

func (s *InternalReader) Setup(stream Stream) error {
	return nil
}

func (s *InternalWriter) Setup(config Config) error {
	return nil
}

func (s *InternalChannel) Teardown() error {
	var err error
	if s.InternalReader != nil {
		err = s.InternalReader.Teardown()
		if err != nil {
			// we log this error as it might get overwritten
			Log.Error(err)
		}
		s.InternalReader = nil
	}
	if s.InternalWriter != nil {
		err = s.InternalWriter.Teardown()
		s.InternalWriter = nil
	}
	return err
}

func (a *InternalChannel) Purge() error {
	if a.InternalReader == nil {
		return nil
	}
	return a.InternalReader.Purge()
}

type BasicInternalReader struct {
	Store *ItemStore
	Model Model
}

type BasicInternalPayload struct {
	items       []*Item
	store       *ItemStore
	endOfStream bool
}

func (p *BasicInternalPayload) EndOfStream() bool {
	return p.endOfStream
}

func (p *BasicInternalPayload) Items() []*Item {
	return p.items
}

func (p *BasicInternalPayload) Headers() map[string]interface{} {
	return map[string]interface{}{}
}

func (p *BasicInternalPayload) Acknowledge() error {
	return nil
}

func (p *BasicInternalPayload) Reject() error {
	return nil
}

func (i *BasicInternalReader) Read() (Payload, error) {
	modelType := i.Model.Type()
	i.Store.Lock()
	defer i.Store.Unlock()
	modelChannels, ok := i.Store.Items[modelType]
	if !ok {
		return nil, nil
	}
	modelID := hex.EncodeToString(i.Model.ID())
	modelPayloads, ok := modelChannels[modelID]
	if !ok {
		// no payloads for this channel
		return nil, nil
	}

	payload := modelPayloads[0]

	if len(modelPayloads) == 1 {
		// we delete the payloads
		delete(modelChannels, modelID)
	} else {
		modelChannels[modelID] = modelPayloads[1:len(modelPayloads)]
	}

	return payload, nil
}

func (i *BasicInternalReader) Setup(Stream) error {
	return fmt.Errorf("setup with stream not supported")
}

func (i *BasicInternalReader) Purge() error {
	return nil
}

func (i *BasicInternalReader) Teardown() error {
	return nil
}

func (i *BasicInternalReader) SetupWithModel(model Model) error {
	Log.Debugf("Setting up reader for model type %s", model.Type())
	i.Model = model
	return nil
}

// Create the internal reader from which source received via the Internal is read
func makeReader(controller Controller, model Model) (ModelReader, error) {
	config, err := controller.Settings().Get("internal-channel")
	if err != nil {
		config = map[string]interface{}{
			"type":   "basic",
			"config": map[string]interface{}{},
		}
	}
	configErr := fmt.Errorf("invalid configuration")
	configMap, ok := maps.ToStringMap(config)
	if !ok {
		return nil, configErr
	}
	readerType, ok := configMap["type"].(string)
	if !ok {
		return nil, configErr
	}
	readerConfig, ok := maps.ToStringMap(configMap["config"])
	if !ok {
		return nil, configErr
	}

	var reader Reader

	if readerType == "basic" {
		itemStoreObj, ok := controller.GetVar("item-store")
		if !ok {
			itemStoreObj = MakeItemStore()
			controller.SetVar("item-store", itemStoreObj)
		}
		itemStore, ok := itemStoreObj.(*ItemStore)
		if !ok {
			return nil, fmt.Errorf("not a valid item store")
		}
		reader = &BasicInternalReader{
			Store: itemStore,
		}
	} else {
		definition, ok := controller.Definitions().ReaderDefinitions[readerType]

		if !ok {
			return nil, configErr
		}

		reader, err = definition.Maker(readerConfig)

		if err != nil {
			return nil, err
		}

	}

	modelReader, ok := reader.(ModelReader)

	if !ok {
		return nil, fmt.Errorf("not a model reader")
	}

	if err := modelReader.SetupWithModel(model); err != nil {
		return nil, err
	}

	return modelReader, nil

}

type ItemStore struct {
	mutex sync.Mutex
	Items map[string]map[string][]Payload
}

func (i *ItemStore) Lock() {
	i.mutex.Lock()
}

func (i *ItemStore) Unlock() {
	i.mutex.Unlock()
}

func MakeItemStore() *ItemStore {
	return &ItemStore{
		Items: make(map[string]map[string][]Payload),
		mutex: sync.Mutex{},
	}
}

type BasicInternalWriter struct {
	Store *ItemStore
	Model Model
}

func (i *BasicInternalWriter) Setup(Config) error {
	return fmt.Errorf("Setup with config not supported")
}

func (i *BasicInternalWriter) Teardown() error {
	return nil
}

func (i *BasicInternalWriter) Close() error {
	Log.Info("Closing internal writer...")
	return nil
}

func (i *BasicInternalWriter) Write(payload Payload) error {
	modelType := i.Model.Type()
	i.Store.Lock()
	defer i.Store.Unlock()
	modelChannels, ok := i.Store.Items[modelType]
	if !ok {
		modelChannels = make(map[string][]Payload)
		i.Store.Items[modelType] = modelChannels
	}
	modelID := hex.EncodeToString(i.Model.ID())
	modelPayloads, ok := modelChannels[modelID]
	if !ok {
		modelPayloads = make([]Payload, 0)
	}
	modelPayloads = append(modelPayloads, payload)
	modelChannels[modelID] = modelPayloads
	return nil
}

func (i *BasicInternalWriter) SetupWithModel(model Model) error {
	i.Model = model
	return nil
}

// Create the internal writer to which source received via the Internal is sent
func makeWriter(controller Controller, model Model) (ModelWriter, error) {
	config, err := controller.Settings().Get("internal-channel")
	if err != nil {
		config = map[string]interface{}{
			"type":   "basic",
			"config": map[string]interface{}{},
		}
	}
	configErr := fmt.Errorf("invalid configuration")
	configMap, ok := maps.ToStringMap(config)
	if !ok {
		return nil, configErr
	}
	writerType, ok := configMap["type"].(string)
	if !ok {
		return nil, configErr
	}
	writerConfig, ok := maps.ToStringMap(configMap["config"])
	if !ok {
		return nil, configErr
	}

	var writer Writer

	if writerType == "basic" {
		// we create a basic internal writer
		itemStoreObj, ok := controller.GetVar("item-store")
		if !ok {
			itemStoreObj = MakeItemStore()
			controller.SetVar("item-store", itemStoreObj)
		}
		itemStore, ok := itemStoreObj.(*ItemStore)
		if !ok {
			return nil, fmt.Errorf("not a valid item store")
		}
		writer = &BasicInternalWriter{
			Store: itemStore,
		}
	} else {
		definition, ok := controller.Definitions().WriterDefinitions[writerType]

		if !ok {
			return nil, configErr
		}

		writer, err = definition.Maker(writerConfig)

		if err != nil {
			return nil, err
		}
	}

	modelWriter, ok := writer.(ModelWriter)

	if !ok {
		return nil, fmt.Errorf("not a model writer")
	}

	if err := modelWriter.SetupWithModel(model); err != nil {
		return nil, err
	}

	return modelWriter, nil
}

func (s *InternalChannel) teardownInternalReader() error {
	return nil
}

// Sets up the Internal reader, which relies on an internal queue to process items.
func (s *InternalChannel) Setup(controller Controller, model Model) error {

	var err error

	if s.InternalReader, err = makeReader(controller, model); err != nil {
		return err
	}

	if s.InternalWriter, err = makeWriter(controller, model); err != nil {
		return err
	}

	return nil
}

func (s *InternalChannel) Read() (Payload, error) {
	return s.InternalReader.Read()
}

// We write items to the internal Internal writer.
func (s *InternalChannel) Write(payload Payload) error {
	return s.InternalWriter.Write(payload)
}

func MakeInternalChannel() *InternalChannel {
	return &InternalChannel{}
}
