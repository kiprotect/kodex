// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
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
	"github.com/kiprotect/kodex"
	"sync"
	"time"
)

type ConfigContext struct {
	Config       kodex.Config
	Processor    *kodex.Processor
	Destinations map[string][]kodex.DestinationMap
}

type LocalStreamWorker struct {
	pool              chan chan kodex.Payload
	acknowledgeFailed bool
	started           bool
	mutex             sync.Mutex
	contexts          []*ConfigContext
	executor          StreamExecutor
	payloadChannel    chan kodex.Payload
	stop              chan bool
}

func MakeLocalStreamWorker(pool chan chan kodex.Payload,
	contexts []*ConfigContext,
	acknowledgeFailed bool,
	executor StreamExecutor) (*LocalStreamWorker, error) {
	// todo: proper error handling

	return &LocalStreamWorker{
		pool:              pool,
		acknowledgeFailed: acknowledgeFailed,
		payloadChannel:    make(chan kodex.Payload),
		stop:              make(chan bool),
		contexts:          contexts,
		started:           false,
		executor:          executor,
	}, nil
}

func (w *LocalStreamWorker) Start() {
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
				break
			case <-w.stop:
				stop = true
			case <-time.After(time.Millisecond):
				if stop && len(w.payloadChannel) == 0 {
					w.stop <- true
					return
				}
			}
		}
	}()
}

func (w *LocalStreamWorker) Stop() {

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

	w.started = false
}

func (w *LocalStreamWorker) ProcessPayload(payload kodex.Payload) error {

	handleError := func(err error) error {
		kodex.Log.Error(err)
		if w.acknowledgeFailed {
			kodex.Log.Warning("Acknowledging failed payload...")
			payload.Acknowledge()
		} else {
			kodex.Log.Warning("Rejecting failed payload...")
			payload.Reject()
		}
		return err
	}

	var items, newItems []*kodex.Item
	var err error

	items = payload.Items()

	kodex.Log.Debugf("Received %d items for payload...", len(items))

	for _, context := range w.contexts {

		if newItems, err = context.Processor.Process(items, nil, payload.EndOfStream()); err != nil {
			kodex.Log.Error("an error occurred")
			return handleError(err)
		}

		for _, destinationMaps := range context.Destinations {

			for _, destinationMap := range destinationMaps {

				// we only write items to active destinations
				if destinationMap.Status() != kodex.ActiveDestination {
					continue
				}

				// we do not perform any writer setup as we already did this before
				writer, err := destinationMap.InternalWriter()

				if err != nil {
					kodex.Log.Error("error writing destination items...")
					return handleError(err)
				}

				if err := writer.Write(kodex.MakeBasicPayload(newItems, payload.Headers(), payload.EndOfStream())); err != nil {
					kodex.Log.Error("error writing items...")
					return handleError(err)
				}
				kodex.Log.Debugf("Wrote %d items", len(newItems))

			}
		}

	}

	payload.Acknowledge()

	return nil

}
