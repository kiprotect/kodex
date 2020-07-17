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

package processing

import (
	"fmt"
	"github.com/kiprotect/kiprotect"
	"sync"
	"time"
)

type LocalSourceReader struct {
	maxSourceWorkers int
	workers          []*LocalSourceWorker
	id               []byte
	pool             chan chan kiprotect.Payload
	sourceMap        kiprotect.SourceMap
	reader           kiprotect.Reader
	configs          []kiprotect.Config
	stopReader       chan bool
	mutex            sync.Mutex
	supervisor       SourceSupervisor
	stopped          bool
	stopping         bool
	payloadChannel   chan kiprotect.Payload
}

func MakeLocalSourceReader(maxSourceWorkers int,
	id []byte) *LocalSourceReader {
	return &LocalSourceReader{
		stopReader:       make(chan bool),
		stopped:          true,
		id:               id,
		payloadChannel:   make(chan kiprotect.Payload, maxSourceWorkers*8),
		maxSourceWorkers: maxSourceWorkers,
	}
}

func (d *LocalSourceReader) ID() []byte {
	return d.id
}

func (d *LocalSourceReader) Start(supervisor SourceSupervisor, sourceMap kiprotect.SourceMap) error {

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if !d.stopped {
		return fmt.Errorf("busy")
	}

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
	return d.stop(true, false)
}

func (d *LocalSourceReader) run() error {

	d.workers = make([]*LocalSourceWorker, 0)

	if d.sourceMap == nil {
		return fmt.Errorf("no source defined")
	}

	// we retrieve all active streams for this source
	streams, err := d.sourceMap.Source().Streams(kiprotect.ActiveSource)

	if err != nil {
		return err
	}

	d.pool = make(chan chan kiprotect.Payload, d.maxSourceWorkers)

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

func (d *LocalSourceReader) SourceMap() kiprotect.SourceMap {
	return d.sourceMap
}

func (d *LocalSourceReader) stop(gracefully bool, fromReader bool) error {

	if d.stopping || d.stopped {
		return nil
	}

	d.mutex.Lock()

	sourceMap := d.sourceMap
	supervisor := d.supervisor
	defer func() {
		d.mutex.Unlock()
		if supervisor != nil {
			supervisor.ReaderStopped(d, sourceMap)
		}
	}()

	d.stopping = true

	if !fromReader {
		// first we stop the source reader to stop reading more payloads..
		d.stopReader <- true
		<-d.stopReader
	}

	// then we stop the workers...
	for _, worker := range d.workers {
		worker.Stop()
	}

	// then we tear down the source reader
	if err := d.reader.Teardown(); err != nil {
		kiprotect.Log.Error(err)
	}

	d.sourceMap = nil
	d.reader = nil
	d.stopped = true
	d.stopping = false
	d.supervisor = nil

	return nil

}

func (d *LocalSourceReader) read() {
Loop:
	for {
		var payload kiprotect.Payload
		var err error

		select {
		case <-d.stopReader:
			// we stop reading any more payloads and return...
			d.stopReader <- true
			break Loop
		case <-time.After(1 * time.Millisecond):
			break
		}

		// to do: check if the source was updated and if yes break out of
		// the loop (to reload configuration)

		if payload, err = d.reader.Read(); err != nil {
			kiprotect.Log.Error(err)
			break
		}

		// we didn't receive any new items...
		if payload == nil {
			continue
		}

		workerChannel := <-d.pool
		workerChannel <- payload

		if payload.EndOfStream() {
			break
		}

	}
	d.stop(true, true)
}
