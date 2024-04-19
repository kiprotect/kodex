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
	"github.com/kiprotect/kodex"
	"sync"
	"time"
)

type LocalSourceWorker struct {
	pool           chan chan kodex.Payload
	started        bool
	ItemsProcessed int
	streams        []kodex.Stream
	channels       []*kodex.InternalChannel
	executor       Executor
	mutex          sync.Mutex
	payloadChannel chan kodex.Payload
	stop           chan bool
}

func makeChannels(streams []kodex.Stream) ([]*kodex.InternalChannel, error) {
	channels := make([]*kodex.InternalChannel, 0)
	for _, stream := range streams {
		channel := kodex.MakeInternalChannel()
		if err := channel.Setup(stream.Project().Controller(), stream); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func MakeLocalSourceWorker(pool chan chan kodex.Payload,
	streams []kodex.Stream,
	executor Executor) (*LocalSourceWorker, error) {
	if channels, err := makeChannels(streams); err != nil {
		return nil, err
	} else {
		return &LocalSourceWorker{
			pool:           pool,
			payloadChannel: make(chan kodex.Payload, 100),
			stop:           make(chan bool),
			streams:        streams,
			channels:       channels,
			started:        false,
			executor:       executor,
		}, nil
	}
}

func (w *LocalSourceWorker) Start() {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	stop := false
	w.started = true

	go func() {
		// we submit our payload channel to the worker pool
		// the source reader will fetch it from the pool and submit items
		// to it. If we want to accept more work we will have to submit the
		// channel again.
		w.pool <- w.payloadChannel
		for {
			select {
			case payload := <-w.payloadChannel:
				w.ItemsProcessed += len(payload.Items())
				w.ProcessPayload(payload)
				w.pool <- w.payloadChannel
			case <-w.stop:
				stop = true
			case <-time.After(time.Millisecond):
				if stop && len(w.payloadChannel) == 0 {
					// we remove the worker channel from the pool of channels
					channels := make([]chan kodex.Payload, 0)
				loop:
					for {
						select {
						case wc := <-w.pool:
							if wc != w.payloadChannel {
								channels = append(channels, wc)
							} else {
								// we have found the channel, we break
								break loop
							}
						default: // no more worker channels
							break loop
						}
					}
					// we resubmit the other channels
					for _, channel := range channels {
						w.pool <- channel
					}
					// we close the payload channel
					close(w.payloadChannel)
					w.started = false
					w.stop <- true
					return
				}
			}
		}
	}()
}

func (w *LocalSourceWorker) Stop() {

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
	w.started = false
	w.channels = nil
}

func (w *LocalSourceWorker) ProcessPayload(payload kodex.Payload) error {

	// we send the items from the payload to the designated internal queues

	handleError := func(err error) error {
		kodex.Log.Error(err)
		return err
	}

	for _, channel := range w.channels {
		if err := channel.Write(payload); err != nil {
			return handleError(err)
		}
	}

	return nil

}
