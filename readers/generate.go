// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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
	"github.com/kiprotect/kodex"
	"time"
)

type GenerateReader struct {
	generators map[string]Generator
	t          time.Time
	residue    float64
	frequency  float64
}

type GeneratePayload struct {
	items   []*kodex.Item
	headers map[string]interface{}
}

func (f *GeneratePayload) Items() []*kodex.Item {
	return f.items
}

func (f *GeneratePayload) Acknowledge() error {
	return nil
}

func (f *GeneratePayload) Reject() error {
	return nil
}

func (s *GeneratePayload) Teardown() error {
	return nil
}

func (f *GeneratePayload) Headers() map[string]interface{} {
	return f.headers
}

func MakeGeneratePayload() *BytesPayload {
	payload := BytesPayload{
		items:   make([]*kodex.Item, 0),
		headers: map[string]interface{}{},
	}
	return &payload
}

type Generator func() interface{}

type GeneratorMaker func(map[string]interface{}) (Generator, error)

func Literal(config map[string]interface{}) (Generator, error) {
	params, err := LiteralForm.Validate(config)
	if err != nil {
		return nil, err
	}
	literal := params["value"]
	return func() interface{} {
		return literal
	}, nil
}

func Timestamp(config map[string]interface{}) (Generator, error) {
	params, err := TimestampForm.Validate(config)
	if err != nil {
		return nil, err
	}
	format := params["format"].(string)
	start := time.Now().UTC()
	inc := time.Duration(time.Second)
	return func() interface{} {
		start = start.Add(inc)
		switch format {
		case "rfc3339":
			return start.Format(time.RFC3339)
		case "unix":
			return start.Unix()
		default:
			panic("invalid format")
		}
	}, nil
}

var generators = map[string]GeneratorMaker{
	"timestamp": Timestamp,
	"literal":   Literal,
}

func MakeGenerateReader(config map[string]interface{}) (kodex.Reader, error) {
	params, err := GenerateForm.Validate(config)
	if err != nil {
		return nil, err
	}

	generators := make(map[string]Generator)

	for key, generator := range params["fields"].(map[string]interface{}) {
		generators[key] = generator.(Generator)
	}

	return &GenerateReader{
		frequency:  params["frequency"].(float64),
		t:          time.Now().UTC(),
		generators: generators,
	}, nil
}

func (f *GenerateReader) Purge() error {
	return nil
}

func (g *GenerateReader) generateItem() *kodex.Item {
	d := make(map[string]interface{})
	for key, generator := range g.generators {
		d[key] = generator()
	}
	return kodex.MakeItem(d)
}

func (g *GenerateReader) Read() (kodex.Payload, error) {
	et := time.Since(g.t)
	nf := g.frequency*float64(et)/1e9 + g.residue
	n := int64(nf)
	// we save the remainder for later
	g.residue = nf - float64(n)
	if n == 0 {
		return nil, nil
	}
	g.t = time.Now().UTC()
	payload := MakeGeneratePayload()
	for i := int64(0); i < n; i++ {
		item := g.generateItem()
		payload.items = append(payload.items, item)
	}
	return payload, nil
}

func (g *GenerateReader) Setup(stream kodex.Stream) error {
	return nil
}

func (g *GenerateReader) Teardown() error {
	return nil
}
