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

package writers

import (
	"compress/gzip"
	"fmt"
	"github.com/kiprotect/kodex"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileWriter struct {
	BasePath string
	Name     string
	Format   string
	Compress bool
	AddTime  bool
	mutex    *sync.Mutex
}

func (s *FileWriter) Teardown() error {
	return nil
}

func (s *FileWriter) Setup(config kodex.Config) error {
	if s.BasePath == "" {
		return nil
	}
	info, err := os.Stat(s.BasePath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("path must be a directory")
	}

	return nil
}

func (s *FileWriter) Write(payload kodex.Payload) error {

	var fileName, extension string

	if s.Compress {
		extension = fmt.Sprintf("%s.gz", s.Format)
	} else {
		extension = s.Format
	}
	if s.AddTime {
		tn := time.Now().UTC().Unix()
		//we rotate the files every 60 seconds
		ts := tn - (tn % 60)
		fileName = fmt.Sprintf("%s-%d.%s", s.Name, ts, extension)
	} else {
		fileName = fmt.Sprintf("%s.%s", s.Name, extension)
	}

	fullPath := filepath.Join(s.BasePath, fileName)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer f.Close()
	defer f.Sync()

	var writer io.Writer = f

	if s.Compress {
		gzWriter := gzip.NewWriter(writer)
		writer = gzWriter
		defer gzWriter.Close()
		defer gzWriter.Flush()
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

	return nil
}

func MakeFileWriter(config map[string]interface{}) (kodex.Writer, error) {

	if params, err := FileWriterForm.Validate(config); err != nil {
		return nil, err
	} else {
		return &FileWriter{
			BasePath: params["path"].(string),
			Name:     params["base-name"].(string),
			AddTime:  params["add-time"].(bool),
			Compress: params["compress"].(bool),
			Format:   params["format"].(string),
			mutex:    &sync.Mutex{},
		}, nil
	}
}
