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
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/kiprotect/kodex"
	"io"
	"os"
)

type FileReader struct {
	Reader     *bufio.Reader
	File       *os.File
	GzReader   *gzip.Reader
	Format     string
	Compressed bool
	Headers    map[string]interface{}
	Path       string
	ChunkSize  int
}

type FilePayload struct {
	items       []*kodex.Item
	endOfStream bool
	headers     map[string]interface{}
}

func (f *FilePayload) EndOfStream() bool {
	return f.endOfStream
}

func (f *FilePayload) Items() []*kodex.Item {
	return f.items
}

func (f *FilePayload) Acknowledge() error {
	return nil
}

func (f *FilePayload) Reject() error {
	return nil
}

func (f *FileReader) Purge() error {
	return nil
}

func (s *FileReader) Setup(stream kodex.Stream) error {
	if s.File != nil {
		return nil
	}
	info, err := os.Stat(s.Path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return fmt.Errorf("reader path is not a file")
	}

	var reader io.Reader

	if s.File, err = os.Open(s.Path); err != nil {
		return err
	}

	if s.Compressed {
		s.GzReader, err = gzip.NewReader(reader)
		if err != nil {
			return err
		}
		reader = s.GzReader
	} else {
		reader = s.File
	}

	s.Reader = bufio.NewReader(reader)

	return nil

}

func (s *FileReader) Teardown() error {
	if s.GzReader != nil {
		if err := s.GzReader.Close(); err != nil {
			return err
		}
	}
	err := s.File.Close()
	s.File = nil
	return err
}

func (s *FileReader) MakeFilePayload() (*FilePayload, error) {

	payload := FilePayload{
		items:   make([]*kodex.Item, 0),
		headers: s.Headers,
	}
	return &payload, nil
}

func (f *FilePayload) Headers() map[string]interface{} {
	return f.headers
}

func (s *FileReader) Read() (kodex.Payload, error) {

	payload, err := s.MakeFilePayload()

	if err != nil {
		return nil, err
	}

	items := make([]*kodex.Item, 0)

	endOfFile := false

	for i := 0; i < s.ChunkSize; i++ {
		item := make(map[string]interface{})
		line, err := s.Reader.ReadBytes('\n')
		if err == io.EOF {
			endOfFile = true
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
	payload.endOfStream = endOfFile

	return payload, nil

}

func MakeFileReader(config map[string]interface{}) (kodex.Reader, error) {
	if params, err := FileReaderForm.Validate(config); err != nil {
		return nil, err
	} else {
		return &FileReader{
			Path:       params["path"].(string),
			ChunkSize:  int(params["chunk-size"].(int64)),
			Headers:    params["headers"].(map[string]interface{}),
			Format:     params["format"].(string),
			Compressed: params["compressed"].(bool),
		}, nil
	}
}
