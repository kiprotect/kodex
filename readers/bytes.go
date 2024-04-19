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

package readers

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/kiprotect/kodex"
	"io"
)

type BytesReader struct {
	Input      []byte
	Reader     *bufio.Reader
	Format     string
	Compressed bool
	Headers    map[string]interface{}
	ChunkSize  int
}

type BytesPayload struct {
	items       []*kodex.Item
	endOfStream bool
	headers     map[string]interface{}
}

func (f *BytesPayload) EndOfStream() bool {
	return f.endOfStream
}

func (f *BytesPayload) Items() []*kodex.Item {
	return f.items
}

func (f *BytesPayload) Acknowledge() error {
	return nil
}

func (f *BytesPayload) Reject() error {
	return nil
}

func (f *BytesReader) Purge() error {
	return nil
}

func (s *BytesReader) Teardown() error {
	return nil
}

func (s *BytesReader) MakeBytesPayload() (*BytesPayload, error) {
	payload := BytesPayload{
		items:   make([]*kodex.Item, 0),
		headers: s.Headers,
	}
	return &payload, nil
}

func (f *BytesPayload) Headers() map[string]interface{} {
	return f.headers
}

func (s *BytesReader) Read() (kodex.Payload, error) {

	payload, err := s.MakeBytesPayload()

	if err != nil {
		return nil, err
	}

	items := make([]*kodex.Item, 0)

	for i := 0; i < s.ChunkSize; i++ {
		item := make(map[string]interface{})
		line, err := s.Reader.ReadBytes('\n')
		if err != nil && !(err == io.EOF && len(line) > 0) {
			break
		}
		switch s.Format {
		case "json":
			err := json.Unmarshal(line, &item)
			if err != nil {
				return nil, err
			}
			item := kodex.MakeItem(item)
			items = append(items, item)
			break
		}
	}

	if len(items) == 0 {
		return nil, nil
	}

	payload.items = items

	return payload, nil

}

func (b *BytesReader) Setup(stream kodex.Stream) error {

	bytesReader := bytes.NewReader(b.Input)

	if b.Compressed {
		gzReader, err := gzip.NewReader(bytesReader)
		if err != nil {
			return err
		}
		b.Reader = bufio.NewReader(gzReader)
	} else {
		b.Reader = bufio.NewReader(bytesReader)
	}
	return nil
}

func MakeBytesReader(config map[string]interface{}) (kodex.Reader, error) {

	if params, err := BytesReaderForm.Validate(config); err != nil {
		return nil, err
	} else {
		return &BytesReader{
			Input:     params["input"].([]byte),
			ChunkSize: int(params["chunk-size"].(int64)),
			Headers:   params["headers"].(map[string]interface{}),
			Format:    params["format"].(string),
		}, nil
	}
}
