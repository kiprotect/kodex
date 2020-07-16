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

package readers_test

import (
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kiprotect"
	pt "github.com/kiprotect/kiprotect/helpers/testing"
	pf "github.com/kiprotect/kiprotect/helpers/testing/fixtures"
	"github.com/kiprotect/kiprotect/readers"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func TestAMQPReader(t *testing.T) {

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

	reader, err := readers.MakeAMQPReader(mapConfig)

	if err != nil {
		t.Fatal(err)
	}

	amqpReader := reader.(*readers.AMQPReader)

	if err := reader.Setup(nil); err != nil {
		t.Fatal(err)
	}

	// we purge all existing messages on this queue
	if _, err := amqpReader.Channel.QueuePurge(amqpReader.QueueName, false); err != nil {
		t.Fatal(err)
	}

	if err := amqpReader.Channel.Confirm(false); err != nil {
		t.Fatal(err)
	}

	confirmChannel := make(chan amqp.Confirmation)
	amqpReader.Channel.NotifyPublish(confirmChannel)

	if err := amqpReader.Channel.Publish(
		amqpReader.Exchange,   // exchange
		amqpReader.RoutingKey, // routing key
		true,                  // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("{\"foo\": \"bar\"}\n{\"baz\": 144}"),
		}); err != nil {
		t.Fatal(err)
	}

	i := 0
	for {
		time.Sleep(time.Millisecond)
		found := false
		select {
		case confirmation := <-confirmChannel:
			if !confirmation.Ack {
				t.Fatal("not confirmed")
			}
			found = true
		default:
			i++
		}
		if found {
			break
		}
		if i > 1000 {
			t.Fatal("no confirmation received")
		}
	}

	i = 0

	var payload kiprotect.Payload

	// we make sure we can read a message from the queue

	for {

		if payload, err = reader.Read(); err != nil {
			t.Fatal(err)
		}

		if payload != nil {
			break
		}

		time.Sleep(time.Millisecond)

		i++
		if i > 100 {
			t.Fatal("no items received")
		}

	}

	items := payload.Items()

	if len(items) != 2 {
		t.Fatal("expected 2 items")
	}

	itemA := items[0]
	itemB := items[1]

	if v, ok := itemA.Get("foo"); !ok || v != "bar" {
		t.Fatal("item A is wrong")
	}

	if v, ok := itemB.Get("baz"); !ok || v != float64(144) {
		t.Fatal("item B is wrong")
	}

	if err := payload.Acknowledge(); err != nil {
		t.Fatal(err)
	}

	if err := reader.Teardown(); err != nil {
		t.Fatal(err)
	}

}
