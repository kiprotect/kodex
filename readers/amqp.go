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

package readers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/writers"
	"github.com/streadway/amqp"
	"io"
	"time"
)

type AMQPReader struct {
	writers.AMQPBase
	deliveries   <-chan amqp.Delivery
	ConsumerName string
}

func MakeAMQPReader(config map[string]interface{}) (kodex.Reader, error) {
	if params, err := AMQPReaderForm.Validate(config); err != nil {
		return nil, err
	} else {
		base, err := writers.MakeAMQPBase(params)
		if err != nil {
			return nil, err
		}
		return &AMQPReader{
			AMQPBase:     base,
			ConsumerName: params["consumer"].(string),
		}, nil
	}
}

type AMQPPayload struct {
	delivery     amqp.Delivery
	compressed   bool
	rejected     bool
	acknowledged bool
	endOfStream  bool
	format       string
	items        []*kodex.Item
	headers      map[string]interface{}
}

func (f *AMQPPayload) EndOfStream() bool {
	return f.endOfStream
}

func (f *AMQPPayload) Items() []*kodex.Item {
	return f.items
}

func (f *AMQPPayload) Acknowledge() error {
	if f.rejected {
		return fmt.Errorf("payload was already rejected")
	}
	f.acknowledged = true
	return f.delivery.Ack(false)
}

func (f *AMQPPayload) Headers() map[string]interface{} {
	return f.headers
}

func (f *AMQPPayload) Reject() error {
	if f.rejected {
		return nil
	}
	if f.acknowledged {
		return fmt.Errorf("payload was already acknowledged")
	}
	f.rejected = true
	return f.delivery.Reject(false)
}

func (a *AMQPReader) Purge() error {
	if a.Channel == nil {
		return nil
	}
	_, err := a.Channel.QueuePurge(a.QueueName, false)
	return err
}

func (a *AMQPReader) MakeAMQPPayload(delivery amqp.Delivery) (*AMQPPayload, error) {

	// we ensure that this message is really designated for the specific
	// object, providing safety against AMQP routing errors.
	if a.Model != nil {
		modelIDStr, ok := delivery.Headers["modelID"].(string)
		if !ok {
			return nil, fmt.Errorf("model ID missing")
		}
		modelTypeStr, ok := delivery.Headers["modelType"].(string)
		if !ok {
			return nil, fmt.Errorf("model type missing")
		}
		modelID, err := hex.DecodeString(modelIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid stream ID")
		}
		if !bytes.Equal(a.Model.ID(), modelID) || a.Model.Type() != modelTypeStr {
			return nil, fmt.Errorf("models do not match")
		}
	}

	var endOfStream bool

	if eos, ok := delivery.Headers["endOfStream"]; ok {
		if eosBool, ok := eos.(bool); ok {
			endOfStream = eosBool
		}
	}

	payload := AMQPPayload{
		delivery:    delivery,
		compressed:  a.Compress,
		format:      a.Format,
		endOfStream: endOfStream,
		items:       make([]*kodex.Item, 0),
		headers:     delivery.Headers,
	}
	if err := payload.readItems(); err != nil {
		return nil, err
	}
	return &payload, nil
}

func (a *AMQPPayload) getReader() (*bufio.Reader, error) {
	bytesReader := bytes.NewReader(a.delivery.Body)

	var reader *bufio.Reader

	if a.compressed {
		gzReader, err := gzip.NewReader(bytesReader)
		if err != nil {
			return nil, err
		}
		reader = bufio.NewReader(gzReader)
	} else {
		reader = bufio.NewReader(bytesReader)
	}
	return reader, nil

}

func (a *AMQPReader) Peek() (kodex.Payload, error) {
	payload, err := a.Read()
	if err != nil {
		return payload, err
	}
	// we reject the payload
	if err := payload.Reject(); err != nil {
		return payload, err
	}
	return payload, err
}

func (a *AMQPReader) Setup(stream kodex.Stream) error {
	return a.AMQPBase.SetupWithModel(nil)
}

func (a *AMQPPayload) readItems() error {

	reader, err := a.getReader()

	if err != nil {
		return err
	}

	items := make([]*kodex.Item, 0)

	var lastErr error

	for {
		item := make(map[string]interface{})
		line, err := reader.ReadBytes('\n')
		if err != nil && !(err == io.EOF && len(line) > 0) {
			break
		}
		switch a.format {
		case "json":
			err := json.Unmarshal(line, &item)
			if err != nil {
				kodex.Log.Errorf("Error unmarshaling item with format JSON.")
				kodex.Log.Error(err)
				lastErr = err
				continue
			}
			item := kodex.MakeItem(item)
			items = append(items, item)
			break
		}
	}

	a.items = items
	return lastErr

}

func (a *AMQPReader) consume() error {
	var err error
	a.deliveries, err = a.Channel.Consume(
		a.QueueName,    // queue
		a.ConsumerName, // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	return err
}

func (a *AMQPReader) Read() (kodex.Payload, error) {

	if a.deliveries == nil {
		if err := a.consume(); err != nil {
			return nil, err
		}
	}

	var delivery amqp.Delivery

	found := false
	select {
	case delivery = <-a.deliveries:
		found = true
	// notice: this time is really important
	case <-time.After(time.Second * 1):
		break
	}

	if !found {
		return nil, nil
	}

	return a.MakeAMQPPayload(delivery)

}
