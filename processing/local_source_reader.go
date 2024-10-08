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

package processing

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"sync"
	"time"
)

type LocalSourceReader struct {
	maxSourceWorkers int
	workers          []*LocalSourceWorker
	id               []byte
	pool             chan chan kodex.Payload
	sourceMap        kodex.SourceMap
	reader           kodex.Reader
	configs          []kodex.Config
	stopChannel      chan bool
	endOfStream      bool
	mutex            sync.Mutex
	supervisor       Supervisor
	stopped          bool
	stopping         bool
	payloadChannel   chan kodex.Payload
}

func MakeLocalSourceReader(maxSourceWorkers int,
	id []byte) *LocalSourceReader {
	return &LocalSourceReader{
		stopChannel:      make(chan bool, 1),
		stopped:          true,
		id:               id,
		payloadChannel:   make(chan kodex.Payload, maxSourceWorkers*8),
		maxSourceWorkers: maxSourceWorkers,
	}
}

func (d *LocalSourceReader) ID() []byte {
	return d.id
}

func (d *LocalSourceReader) Start(supervisor Supervisor, processable kodex.Processable) error {

	sourceMap, ok := processable.(kodex.SourceMap)

	if !ok {
		return fmt.Errorf("not a source map")
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.stopped {
		return fmt.Errorf("busy")
	}

	d.endOfStream = false
	d.sourceMap = sourceMap
	d.supervisor = supervisor

	var err error
	if d.reader, err = d.sourceMap.Source().Reader(); err != nil {
		return err
	}

	if err := d.reader.Setup(nil); err != nil {
		return err
	}

	return d.run()
}

func (d *LocalSourceReader) Stop(graceful bool) error {
	return d.stop(graceful)
}

func (d *LocalSourceReader) run() error {

	d.workers = make([]*LocalSourceWorker, 0)

	if d.sourceMap == nil {
		return fmt.Errorf("no source defined")
	}

	// we retrieve all active streams for this source
	streams, err := d.sourceMap.Source().Streams(kodex.ActiveSource)

	if err != nil {
		return err
	}

	d.pool = make(chan chan kodex.Payload, d.maxSourceWorkers)

	for i := 0; i < d.maxSourceWorkers; i++ {
		worker, err := MakeLocalSourceWorker(d.pool, streams, d)
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

func (d *LocalSourceReader) Stopped() bool {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.stopped
}

func (d *LocalSourceReader) SourceMap() kodex.SourceMap {
	return d.sourceMap
}

func (d *LocalSourceReader) stop(graceful bool) error {

	if d.stopping || d.stopped {
		return nil
	}

	d.mutex.Lock()

	sourceMap := d.sourceMap
	supervisor := d.supervisor
	defer func() {
		d.sourceMap = nil
		d.reader = nil
		d.stopped = true
		d.stopping = false
		d.supervisor = nil
		d.mutex.Unlock()
		if supervisor != nil {
			supervisor.ExecutorStopped(d, sourceMap)
		}
	}()

	d.stopping = true

	// first we stop the source reader to stop reading more payloads..
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

	kodex.Log.Debugf("%d items processed by reader workers", itemsProcessed)

	// then we tear down the source reader
	if err := d.reader.Teardown(); err != nil {
		kodex.Log.Error(err)
	}

	return nil

}

func (d *LocalSourceReader) read() {
	stopping := false
	stopRequested := false

	stop := func() {
		if !stopRequested && !stopping {
			stopRequested = true
			go d.stop(true)
		}
	}

	for {
		var payload kodex.Payload
		var err error

		select {
		case <-d.stopChannel:
			// we stop...
			d.stopChannel <- true
			return
		case <-time.After(time.Millisecond):
			break
		}

		// to do: check if the source was updated and if yes break out of
		// the loop (to reload configuration)

		if payload, err = d.reader.Read(); err != nil {
			stop()
			continue
		}

		// we didn't receive any new items...
		if payload == nil || len(payload.Items()) == 0 {
			if stopping {

				return
			}
		}

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
