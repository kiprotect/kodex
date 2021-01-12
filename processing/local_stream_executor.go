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

package processing

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"sync"
	"time"
)

type LocalStreamExecutor struct {
	maxStreamWorkers int
	workers          []*LocalStreamWorker
	id               []byte
	pool             chan chan kodex.Payload
	stream           kodex.Stream
	channel          *kodex.InternalChannel
	contexts         []*ConfigContext
	stopChannel      chan bool
	mutex            sync.Mutex
	supervisor       StreamSupervisor
	stopped          bool
	stopping         bool
	payloadChannel   chan kodex.Payload
}

func MakeLocalStreamExecutor(maxStreamWorkers int,
	id []byte) *LocalStreamExecutor {
	return &LocalStreamExecutor{
		stopChannel:      make(chan bool),
		stopped:          true,
		id:               id,
		payloadChannel:   make(chan kodex.Payload, maxStreamWorkers*8),
		maxStreamWorkers: maxStreamWorkers,
	}
}

func (d *LocalStreamExecutor) ID() []byte {
	return d.id
}

func (d *LocalStreamExecutor) Start(supervisor StreamSupervisor, stream kodex.Stream) error {

	kodex.Log.Debugf("Executing stream %s", string(stream.ID()))

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.stopped {
		return fmt.Errorf("busy")
	}

	d.stream = stream
	d.supervisor = supervisor
	d.channel = kodex.MakeInternalChannel()

	if err := d.channel.Setup(stream.Project().Controller(), stream); err != nil {
		return err
	}

	return d.run()
}

func (d *LocalStreamExecutor) Stop(graceful bool) error {
	return d.stop(graceful)
}

func (d *LocalStreamExecutor) run() error {

	d.workers = make([]*LocalStreamWorker, 0)

	if d.stream == nil {
		return fmt.Errorf("no stream defined")
	}

	// we retrieve all active configs for this stream
	configs, err := d.stream.Configs()

	if err != nil {
		return err
	}

	activeConfigs := make([]kodex.Config, 0)

	for _, config := range configs {
		if config.Status() == kodex.ActiveConfig {
			activeConfigs = append(activeConfigs, config)
		}
	}

	d.pool = make(chan chan kodex.Payload, d.maxStreamWorkers)

	d.contexts, err = makeContexts(activeConfigs)

	if err != nil {
		return err
	}

	for i := 0; i < d.maxStreamWorkers; i++ {
		worker, err := MakeLocalStreamWorker(d.pool, d.contexts, false, d)
		if err != nil {
			return err
		}
		worker.Start()
		d.workers = append(d.workers, worker)
	}

	d.stopped = false

	go d.read()

	return nil
}

func (d *LocalStreamExecutor) Stopped() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.stopped
}

func (d *LocalStreamExecutor) Stream() kodex.Stream {
	return d.stream
}

func makeContexts(configs []kodex.Config) ([]*ConfigContext, error) {
	contexts := make([]*ConfigContext, 0)

	for _, config := range configs {
		var processor *kodex.Processor
		var destinations map[string][]kodex.DestinationMap
		var err error
		if processor, err = config.Processor(); err != nil {
			return nil, err
		}
		if err = processor.Setup(); err != nil {
			return nil, err
		}
		if destinations, err = config.Destinations(); err != nil {
			return nil, err
		}

		context := &ConfigContext{
			Config:       config,
			Processor:    processor,
			Destinations: destinations,
		}
		contexts = append(contexts, context)
	}
	return contexts, nil
}

func (d *LocalStreamExecutor) stop(graceful bool) error {

	if d.stopping || d.stopped {
		return nil
	}

	d.mutex.Lock()

	stream := d.stream
	supervisor := d.supervisor

	defer func() {

		d.stream = nil
		d.contexts = nil
		d.channel = nil
		d.stopped = true
		d.stopping = false
		d.supervisor = nil

		if supervisor != nil {
			supervisor.ExecutorStopped(d, stream)
		}

		d.mutex.Unlock()
	}()

	d.stopping = true

	d.stopChannel <- true
	<-d.stopChannel

	// then we stop the workers...
	for _, worker := range d.workers {
		worker.Stop()
	}

	// then we tear down the stream channel
	if err := d.channel.Teardown(); err != nil {
		kodex.Log.Error(err)
	}

	// we tear down all processors and writers...
	for _, context := range d.contexts {

		if err := context.Processor.Teardown(); err != nil {
			kodex.Log.Error(err)
		}

		for _, destinationMaps := range context.Destinations {
			for _, destinationMap := range destinationMaps {
				writer, err := destinationMap.InternalWriter()
				if err != nil {
					kodex.Log.Error(err)
					continue
				}
				if err := writer.Teardown(); err != nil {
					kodex.Log.Error(err)
				}
			}
		}
	}

	return nil

}

func (d *LocalStreamExecutor) read() {
	var stopping bool
	for {
		var payload kodex.Payload
		var err error

		select {
		case <-d.stopChannel:
			// we stop reading any more payloads and return...
			d.stopChannel <- true
			return
		case <-time.After(time.Second):
			break
		}

		// to do: check if the stream was updated and if yes break out of
		// the loop (to reload configuration)

		if payload, err = d.channel.Read(); err != nil {
			if !stopping {
				stopping = true
				go d.stop(true)
			}
			continue
		}

		// we didn't receive any new items...
		if payload == nil {
			// we stop processing any further items...
			if !stopping {
				stopping = true
				go d.stop(true)
			}
			continue
		}

		workerChannel := <-d.pool
		workerChannel <- payload

		if payload.EndOfStream() {
			if !stopping {
				stopping = true
				go d.stop(true)
			}
		}
	}
}
