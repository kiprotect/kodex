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

package writers

import (
	"github.com/kiprotect/kodex"
	"sync"
)

type InMemoryWriter struct {
	writer []map[string]interface{}
	mutex  *sync.Mutex
}

func (s *InMemoryWriter) Teardown() error {
	return nil
}

func (s *InMemoryWriter) Setup(config kodex.Config) error {
	return nil
}

func (s *InMemoryWriter) Result() []map[string]interface{} {
	return s.writer
}

func (s *InMemoryWriter) Write(payload kodex.Payload) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, item := range payload.Items() {
		s.writer = append(s.writer, item.All())
	}

	return nil
}

func MakeInMemoryWriter(params map[string]interface{}) (kodex.Writer, error) {

	return &InMemoryWriter{
		writer: make([]map[string]interface{}, 0),
		mutex:  &sync.Mutex{},
	}, nil
}
