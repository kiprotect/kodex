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

type LocalDestinationWriter struct {
	maxDestinationWorkers int
	workers               []*LocalDestinationWorker
	id                    []byte
	pool                  chan chan kodex.Payload
	destinationMap        kodex.DestinationMap
	writer                kodex.Writer
	endOfStream           bool
	channel               *kodex.InternalChannel
	stopChannel           chan bool
	mutex                 sync.Mutex
	supervisor            Supervisor
	stopped               bool
	stopping              bool
	payloadChannel        chan kodex.Payload
}

func MakeLocalDestinationWriter(maxDestinationWorkers int,
	id []byte) *LocalDestinationWriter {
	return &LocalDestinationWriter{
		stopChannel:           make(chan bool),
		stopped:               true,
		id:                    id,
		payloadChannel:        make(chan kodex.Payload, maxDestinationWorkers*8),
		maxDestinationWorkers: maxDestinationWorkers,
	}
}

func (d *LocalDestinationWriter) ID() []byte {
	return d.id
}

func (d *LocalDestinationWriter) Start(supervisor Supervisor, processable kodex.Processable) error {

	destinationMap, ok := processable.(kodex.DestinationMap)

	if !ok {
		return fmt.Errorf("not a destination map")
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.stopped {
		return fmt.Errorf("busy")
	}

	d.endOfStream = false
	d.destinationMap = destinationMap
	d.supervisor = supervisor

	var err error

	if d.writer, err = d.destinationMap.Destination().Writer(); err != nil {
		return err
	}

	if err := d.writer.Setup(d.destinationMap.Config()); err != nil {
		return err
	}

	d.channel = kodex.MakeInternalChannel()

	if err := d.channel.Setup(destinationMap.Destination().Project().Controller(), destinationMap); err != nil {
		return err
	}

	return d.run()
}

func (d *LocalDestinationWriter) Stop(graceful bool) error {
	return d.stop(graceful)
}

func (d *LocalDestinationWriter) run() error {

	d.workers = make([]*LocalDestinationWorker, 0)

	if d.destinationMap == nil {
		return fmt.Errorf("no destination map defined")
	}

	d.pool = make(chan chan kodex.Payload, d.maxDestinationWorkers)

	for i := 0; i < d.maxDestinationWorkers; i++ {
		worker, err := MakeLocalDestinationWorker(d.pool, d.writer, d)
		if err != nil {
			return err
		}
		worker.Start()
		d.workers = append(d.workers, worker)
	}

	d.stopped = false

	go d.write()

	return nil
}

func (d *LocalDestinationWriter) Stopped() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.stopped
}

func (d *LocalDestinationWriter) stop(graceful bool) error {

	if d.stopping || d.stopped {
		return nil
	}

	d.mutex.Lock()

	destinationMap := d.destinationMap
	supervisor := d.supervisor

	d.stopping = true

	defer func() {
		d.mutex.Unlock()
		d.destinationMap = nil
		d.writer = nil
		d.stopped = true
		d.stopping = false
		d.supervisor = nil
	}()

	// first we stop the destination writer to stop reading more payloads..
	d.stopChannel <- true
	<-d.stopChannel

	itemsProcessed := 0

	// then we stop the workers...
	for i, worker := range d.workers {
		// we submit the "end of stream" payload to the last active worker
		// to ensure it will be processed as the last payload
		if d.endOfStream && i == len(d.workers)-1 {
			endOfStreamPayload := kodex.MakeBasicPayload([]*kodex.Item{}, map[string]interface{}{}, true)
			workerChannel := <-d.pool
			workerChannel <- endOfStreamPayload
		}
		worker.Stop()
		itemsProcessed += worker.ItemsProcessed
	}

	kodex.Log.Debugf("%d items processed by destination workers", itemsProcessed)

	// then we tear down the destination writer
	if err := d.writer.Teardown(); err != nil {
		kodex.Log.Error(err)
	}

	if err := d.channel.Teardown(); err != nil {
		kodex.Log.Error(err)
	}

	if supervisor != nil {
		supervisor.ExecutorStopped(d, destinationMap)
	}

	return nil

}

func (d *LocalDestinationWriter) write() {

	stopping := false

	stop := func() {
		if !stopping {
			stopping = true
			go d.stop(true)
		}
	}

	itemsProcessed := 0

	for {
		var payload kodex.Payload
		var err error

		select {
		case <-d.stopChannel:
			// we stop reading any more payloads and return...
			stopping = true
		case <-time.After(time.Millisecond):
			break
		}

		// to do: check if the destination was updated and if yes break out of
		// the loop (to reload configuration)

		if payload, err = d.channel.Read(); err != nil {
			kodex.Log.Error(err)
			stop()
			continue
		}

		// we didn't receive any new items...
		if payload == nil {
			if stopping {
				d.stopChannel <- true
				kodex.Log.Debugf("%d items processed in stream", itemsProcessed)
				return
			}
			continue
		}

		itemsProcessed += len(payload.Items())

		if payload.EndOfStream() {
			// we replace the "end of stream payload" and instead send a replacement
			// payload during the stop process to ensure that it will be processed last
			replacedPayload := kodex.MakeBasicPayload(payload.Items(), payload.Headers(), false)
			workerChannel := <-d.pool
			workerChannel <- replacedPayload
			d.mutex.Lock()
			d.endOfStream = true
			d.mutex.Unlock()
		} else {
			workerChannel := <-d.pool
			workerChannel <- payload
		}

		if payload.EndOfStream() {
			stop()
			continue
		}
	}
}
