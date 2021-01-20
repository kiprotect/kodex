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
	channel               *kodex.InternalChannel
	stopWriter            chan bool
	mutex                 sync.Mutex
	supervisor            Supervisor
	stopped               bool
	stopping              bool
	payloadChannel        chan kodex.Payload
}

func MakeLocalDestinationWriter(maxDestinationWorkers int,
	id []byte) *LocalDestinationWriter {
	return &LocalDestinationWriter{
		stopWriter:            make(chan bool),
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
	return d.stop(true, false)
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

func (d *LocalDestinationWriter) stop(gracefully bool, fromReader bool) error {

	if d.stopping || d.stopped {
		return nil
	}

	d.mutex.Lock()

	destinationMap := d.destinationMap
	supervisor := d.supervisor

	d.stopping = true

	if !fromReader {
		// first we stop the destination writer to stop reading more payloads..
		d.stopWriter <- true
		<-d.stopWriter
	}

	// then we stop the workers...
	for _, worker := range d.workers {
		worker.Stop()
	}

	// then we tear down the destination writer
	if err := d.writer.Teardown(); err != nil {
		kodex.Log.Error(err)
	}

	if err := d.channel.Teardown(); err != nil {
		kodex.Log.Error(err)
	}

	d.destinationMap = nil
	d.writer = nil
	d.stopped = true
	d.stopping = false
	d.supervisor = nil

	if supervisor != nil {
		supervisor.ExecutorStopped(d, destinationMap)
	}

	d.mutex.Unlock()

	return nil

}

func (d *LocalDestinationWriter) write() {
	for {
		var payload kodex.Payload
		var err error

		select {
		case <-d.stopWriter:
			// we stop reading any more payloads and return...
			d.stopWriter <- true
			return
		case <-time.After(time.Second):
			break
		}

		// to do: check if the destination was updated and if yes break out of
		// the loop (to reload configuration)

		if payload, err = d.channel.Read(); err != nil {
			kodex.Log.Error(err)
			break
		}

		// we didn't receive any new items...
		if payload == nil {
			break
		}

		workerChannel := <-d.pool
		workerChannel <- payload

		if payload.EndOfStream() {
			break
		}
	}
	d.stop(true, true)
}
