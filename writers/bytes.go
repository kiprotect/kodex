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
	"github.com/kiprotect/kodex"
	"io"
	"sync"
)

type BytesWriter struct {
	Output   []byte
	Format   string
	Compress bool
	mutex    *sync.Mutex
}

func (s *BytesWriter) Teardown() error {
	return nil
}

func (s *BytesWriter) Setup(config kodex.Config) error {
	return nil
}

func (s *BytesWriter) Write(payload kodex.Payload) error {

	var buf *bytes.Buffer

	var writer io.Writer
	var gzWriter *gzip.Writer
	var bufioWriter *bufio.Writer

	b := make([]byte, 0)
	buf = bytes.NewBuffer(b)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.Compress {
		gzWriter = gzip.NewWriter(buf)
		writer = gzWriter
	} else {
		bufioWriter = bufio.NewWriter(buf)
		writer = bufioWriter
	}

	for _, item := range payload.Items() {
		serializedItem, err := item.Serialize(s.Format)
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

	if s.Compress {
		gzWriter.Flush()
		gzWriter.Close()
	} else {
		bufioWriter.Flush()
	}

	s.Output = append(s.Output, buf.Bytes()...)

	return nil
}

func MakeBytesWriter(config map[string]interface{}) (kodex.Writer, error) {
	if params, err := BytesWriterForm.Validate(config); err != nil {
		return nil, err
	} else {
		return &BytesWriter{
			Format:   params["format"].(string),
			Compress: params["compress"].(bool),
			Output:   make([]byte, 0),
			mutex:    &sync.Mutex{},
		}, nil
	}
}
