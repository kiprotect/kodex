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
	"compress/gzip"
	"encoding/json"
	"github.com/kiprotect/kodex"
	"io"
	"os"
)

type StdinReader struct {
	Reader     *bufio.Reader
	GzReader   *gzip.Reader
	Format     string
	Compressed bool
	Headers    map[string]interface{}
	ChunkSize  int
}

type StdinPayload struct {
	items       []*kodex.Item
	endOfStream bool
	headers     map[string]interface{}
}

func (f *StdinPayload) EndOfStream() bool {
	return f.endOfStream
}

func (f *StdinPayload) Items() []*kodex.Item {
	return f.items
}

func (f *StdinPayload) Acknowledge() error {
	return nil
}

func (f *StdinPayload) Reject() error {
	return nil
}

func (f *StdinReader) Purge() error {
	return nil
}

func (s *StdinReader) Setup(stream kodex.Stream) error {

	var reader io.Reader
	var err error

	if s.Compressed {
		s.GzReader, err = gzip.NewReader(os.Stdin)
		if err != nil {
			return err
		}
		reader = s.GzReader
	} else {
		reader = os.Stdin
	}

	s.Reader = bufio.NewReader(reader)

	return nil

}

func (s *StdinReader) Teardown() error {
	if s.GzReader != nil {
		if err := s.GzReader.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (s *StdinReader) MakeStdinPayload() (*StdinPayload, error) {

	payload := StdinPayload{
		items:   make([]*kodex.Item, 0),
		headers: s.Headers,
	}
	return &payload, nil
}

func (f *StdinPayload) Headers() map[string]interface{} {
	return f.headers
}

func (s *StdinReader) Read() (kodex.Payload, error) {

	payload, err := s.MakeStdinPayload()

	if err != nil {
		return nil, err
	}

	items := make([]*kodex.Item, 0)

	endOfStdin := false

	for i := 0; i < s.ChunkSize; i++ {
		item := make(map[string]interface{})
		line, err := s.Reader.ReadBytes('\n')
		if err == io.EOF {
			endOfStdin = true
			if len(line) == 0 {
				break
			}
		} else if err != nil {
			return nil, err
		}
		if len(line) <= 1 {
			continue
		}
		switch s.Format {
		case "json":
			err := json.Unmarshal(line, &item)
			if err != nil {
				return nil, err
			}
			item := kodex.MakeItem(item)
			items = append(items, item)
		}
	}

	kodex.Log.Debugf("Read %d items...", len(items))

	if len(items) == 0 {
		return nil, nil
	}

	payload.items = items
	payload.endOfStream = endOfStdin

	return payload, nil

}

func MakeStdinReader(config map[string]interface{}) (kodex.Reader, error) {
	if params, err := StdinReaderForm.Validate(config); err != nil {
		return nil, err
	} else {
		return &StdinReader{
			ChunkSize:  int(params["chunk-size"].(int64)),
			Headers:    params["headers"].(map[string]interface{}),
			Format:     params["format"].(string),
			Compressed: params["compressed"].(bool),
		}, nil
	}
}
