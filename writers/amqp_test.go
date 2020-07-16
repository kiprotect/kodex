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

package writers_test

import (
	"bytes"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kiprotect"
	pt "github.com/kiprotect/kiprotect/helpers/testing"
	pf "github.com/kiprotect/kiprotect/helpers/testing/fixtures"
	"github.com/kiprotect/kiprotect/writers"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func TestAMQPWriter(t *testing.T) {

	var fixtureConfig = []pt.FC{
		pt.FC{&pf.Settings{}, "settings"},
	}

	fixtures, err := pt.SetupFixtures(fixtureConfig)
	defer pt.TeardownFixtures(fixtureConfig, fixtures)

	if err != nil {
		t.Fatal(err)
	}

	st, _ := fixtures["settings"].(kiprotect.Settings)

	config, err := st.Get("testing.amqp")

	if err != nil {
		kiprotect.Log.Info("Skipping test, no AMQP URL specified...")
		return
	}

	mapConfig, ok := maps.ToStringMap(config)

	if !ok {
		t.Fatal("invalid config")
	}

	writer, err := writers.MakeAMQPWriter(mapConfig)

	if err != nil {
		t.Fatal(err)
	}

	amqpWriter := writer.(*writers.AMQPWriter)

	if err := writer.Setup(nil); err != nil {
		t.Fatal(err)
	}

	// we purge all existing messages on this queue
	if _, err := amqpWriter.Channel.QueuePurge(amqpWriter.QueueName, false); err != nil {
		t.Fatal(err)
	}

	items := []*kiprotect.Item{
		kiprotect.MakeItem(map[string]interface{}{"foo": "bar"}),
	}

	if err := writer.Write(kiprotect.MakeBasicPayload(items, map[string]interface{}{}, false)); err != nil {
		t.Fatal(err)
	}

	deliveries, err := amqpWriter.Channel.Consume(
		amqpWriter.QueueName, // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)

	if err != nil {
		t.Fatal(err)
	}

	// we see if we can receive the message
	i := 0
	var delivery amqp.Delivery
	for {
		received := false
		select {
		case delivery = <-deliveries:
			received = true
		}
		if received {
			break
		}
		i++
		if i > 1000 {
			t.Fatal("no message received")
		}
		time.Sleep(time.Millisecond)

	}

	if !bytes.Equal(delivery.Body, []byte("{\"foo\":\"bar\"}\n")) {
		t.Fatal("wrong body", string(delivery.Body))
	}

}
