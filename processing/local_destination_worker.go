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
	"github.com/kiprotect/kodex"
	"sync"
	"time"
)

type LocalDestinationWorker struct {
	pool              chan chan kodex.Payload
	started           bool
	writer            kodex.Writer
	channels          []*kodex.InternalChannel
	destinationWriter DestinationWriter
	mutex             sync.Mutex
	payloadChannel    chan kodex.Payload
	stop              chan bool
}

func MakeLocalDestinationWorker(pool chan chan kodex.Payload,
	writer kodex.Writer,
	destinationWriter DestinationWriter) (*LocalDestinationWorker, error) {
	return &LocalDestinationWorker{
		pool:              pool,
		payloadChannel:    make(chan kodex.Payload, 100),
		stop:              make(chan bool),
		destinationWriter: destinationWriter,
		started:           false,
		writer:            writer,
	}, nil
}

func (w *LocalDestinationWorker) Start() {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	stop := false
	w.started = true

	go func() {
		w.pool <- w.payloadChannel
		for {
			select {
			case payload := <-w.payloadChannel:
				w.ProcessPayload(payload)
				w.pool <- w.payloadChannel
			case <-w.stop:
				stop = true
			case <-time.After(time.Millisecond):
				if stop && len(w.payloadChannel) == 0 {
					w.started = false
					w.stop <- true
					return
				}
			}
		}
	}()
}

func (w *LocalDestinationWorker) Stop() {

	if !w.started {
		return
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.started {
		return
	}

	w.stop <- true
	<-w.stop

	for _, channel := range w.channels {
		if err := channel.Teardown(); err != nil {
			kodex.Log.Error(err)
		}
	}
	w.channels = nil
	w.started = false
}

func (w *LocalDestinationWorker) ProcessPayload(payload kodex.Payload) error {

	// we send the items from the payload to the designated internal queues

	handleError := func(err error) error {
		if err := payload.Reject(); err != nil {
			kodex.Log.Error(err)
		}
		kodex.Log.Error(err)
		return err
	}

	if err := w.writer.Write(payload); err != nil {
		return handleError(err)
	}

	return payload.Acknowledge()

}
