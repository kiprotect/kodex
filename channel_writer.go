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

package kodex

type MessageType string

const (
	Info  MessageType = "INFO"
	Debug MessageType = "DEBUG"
	Quota MessageType = "QUOTA"
)

type ChannelWriter interface {
	Message(item *Item, data map[string]interface{}, mt MessageType) error
	Write(string, []*Item) error
	Error(*Item, error) error
	Warning(*Item, error) error
}

type Message struct {
	Type MessageType
	Item *Item
	Data map[string]interface{}
}

type Error struct {
	Error error
	Item  *Item
}

type Warning struct {
	Warning error
	Item    *Item
}

type InMemoryChannelWriter struct {
	Items    map[string][]*Item
	Messages []*Message
	Errors   []*Error
	Warnings []*Warning
}

func MakeInMemoryChannelWriter() *InMemoryChannelWriter {
	return &InMemoryChannelWriter{
		Items:    make(map[string][]*Item),
		Messages: make([]*Message, 0),
		Errors:   make([]*Error, 0),
		Warnings: make([]*Warning, 0),
	}
}

func (c *InMemoryChannelWriter) Message(item *Item, data map[string]interface{}, mt MessageType) error {
	c.Messages = append(c.Messages, &Message{
		Item: item,
		Type: mt,
		Data: data,
	})
	return nil
}

func (c *InMemoryChannelWriter) Write(channel string, items []*Item) error {
	if _, ok := c.Items[channel]; !ok {
		c.Items[channel] = make([]*Item, 0)
	}
	c.Items[channel] = append(c.Items[channel], items...)
	return nil
}

func (c *InMemoryChannelWriter) Error(item *Item, err error) error {
	c.Errors = append(c.Errors, &Error{
		Item:  item,
		Error: err,
	})
	return nil
}

func (c *InMemoryChannelWriter) Warning(item *Item, warn error) error {
	c.Warnings = append(c.Warnings, &Warning{
		Item:    item,
		Warning: warn,
	})
	return nil
}
