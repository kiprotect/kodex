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

package writers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/streadway/amqp"
	"io"
	"time"
)

type AMQPBase struct {
	Connection          *amqp.Connection
	Channel             *amqp.Channel
	Queue               amqp.Queue
	QueueExpiresAfterMs int64
	URL                 string
	Format              string
	Compress            bool
	QueueName           string
	RoutingKey          string
	BaseRoutingKey      string
	BaseQueueName       string
	Exchange            string
	ExchangeType        string
	Model               kodex.Model
}

type AMQPWriter struct {
	ConfirmationTimeout float64
	Confirmations       chan amqp.Confirmation
	AMQPBase
}

func MakeAMQPBase(params map[string]interface{}) (AMQPBase, error) {
	return AMQPBase{
		URL:                 params["url"].(string),
		Compress:            params["compress"].(bool),
		BaseRoutingKey:      params["routing_key"].(string),
		BaseQueueName:       params["queue"].(string),
		QueueExpiresAfterMs: params["queue_expires_after_ms"].(int64),
		Exchange:            params["exchange"].(string),
		ExchangeType:        params["exchange_type"].(string),
		Format:              params["format"].(string),
	}, nil
}

func MakeAMQPWriter(config map[string]interface{}) (kodex.Writer, error) {
	if params, err := AMQPWriterForm.Validate(config); err != nil {
		return nil, err
	} else {
		base, err := MakeAMQPBase(params)
		if err != nil {
			return nil, err
		}
		return &AMQPWriter{
			ConfirmationTimeout: params["confirmation_timeout"].(float64),
			AMQPBase:            base,
		}, nil
	}
}

func (a *AMQPWriter) Write(payload kodex.Payload) error {
	var buf *bytes.Buffer

	var writer io.Writer
	var gzWriter *gzip.Writer
	var bufioWriter *bufio.Writer

	b := make([]byte, 0)
	buf = bytes.NewBuffer(b)

	if a.Compress {
		gzWriter = gzip.NewWriter(buf)
		writer = gzWriter
	} else {
		bufioWriter = bufio.NewWriter(buf)
		writer = bufioWriter
	}

	for _, item := range payload.Items() {
		serializedItem, err := item.Serialize(a.Format)
		if err != nil {
			return err
		}
		_, err = writer.Write(serializedItem)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	if a.Compress {
		gzWriter.Flush()
		gzWriter.Close()
	} else {
		bufioWriter.Flush()
	}

	bs := buf.Bytes()

	headers := amqp.Table{
		"format":      a.Format,
		"compress":    a.Compress,
		"endOfStream": payload.EndOfStream(),
	}

	// we add the config ID as a header, if it is defined
	if a.Model != nil {
		headers["modelID"] = hex.EncodeToString(a.Model.ID())
		headers["modelType"] = a.Model.Type()
	}

	if err := a.Channel.Publish(
		a.Exchange,
		a.RoutingKey,
		true,  // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/kiprotect",
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
			Body:         bs,
		},
	); err != nil {
		return err
	}

	select {
	case _ = <-a.Confirmations:
		break
	case <-time.After(time.Nanosecond * time.Duration(a.ConfirmationTimeout*1e9)):
		return fmt.Errorf("timeout while waiting for confirmation")
	}

	return nil
}

func (a *AMQPWriter) Setup(config kodex.Config) error {
	return a.setup(nil)
}

func (a *AMQPWriter) SetupWithModel(model kodex.Model) error {
	return a.setup(model)
}

func (a *AMQPWriter) setup(model kodex.Model) error {
	if err := a.AMQPBase.SetupWithModel(model); err != nil {
		return err
	}

	a.Confirmations = make(chan amqp.Confirmation, 1000)

	if err := a.Channel.Confirm(false); err != nil {
		return err
	}

	// we create a confirmation channel for publications
	a.Channel.NotifyPublish(a.Confirmations)

	return nil
}

func (a *AMQPBase) SetupWithModel(model kodex.Model) error {
	return a.setup(model)
}

func (a *AMQPBase) Setup(config kodex.Config) error {
	return a.setup(nil)
}

func (a *AMQPBase) setup(model kodex.Model) error {
	var err error

	a.Model = model

	if a.Model != nil {
		strId := hex.EncodeToString(model.ID())
		typeName := model.Type()
		a.RoutingKey = a.BaseRoutingKey + fmt.Sprintf(".%s.%s", typeName, strId)
		a.QueueName = a.BaseQueueName + fmt.Sprintf("-%s-%s", typeName, strId)
	} else {
		a.RoutingKey = a.BaseRoutingKey
		a.QueueName = a.BaseQueueName
	}

	if a.Connection, err = amqp.Dial(a.URL); err != nil {
		return err
	}
	if a.Channel, err = a.Connection.Channel(); err != nil {
		return err
	}
	// we do not prefetch any messages
	if err = a.Channel.Qos(1, 0, true); err != nil {
		return err
	}
	// we declare an exchange
	if err := a.Channel.ExchangeDeclare(
		a.Exchange,
		a.ExchangeType,
		true,  // durable
		true,  // autodelete
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return err
	}

	args := map[string]interface{}{}
	if a.QueueExpiresAfterMs != 0 {
		args["x-expires"] = a.QueueExpiresAfterMs
	}

	// we declare a new queue
	if a.Queue, err = a.Channel.QueueDeclare(
		a.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		args,        // arguments
	); err != nil {
		return err
	}

	// we make sure our queue is bound to the exchange
	if err = a.Channel.QueueBind(
		a.QueueName,  // queue name
		a.RoutingKey, // routing key
		a.Exchange,   // exchange name
		false,        // no-wait
		nil,          // arguments
	); err != nil {
		return err
	}

	return nil
}

func (a *AMQPBase) Teardown() error {
	var err error
	if a.Channel != nil {
		err = a.Channel.Close()
		a.Channel = nil
	}
	if a.Connection != nil {
		if !a.Connection.IsClosed() {
			err = a.Connection.Close()
		}
		a.Connection = nil
	}
	return err
}
