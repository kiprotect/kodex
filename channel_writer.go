package kiprotect

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
