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

package writers

import (
	"github.com/kiprotect/kiprotect"
	"sync"
)

type CountWriter struct {
	count     int64
	max       int
	lastItems []map[string]interface{}
	mutex     *sync.Mutex
}

func (s *CountWriter) Teardown() error {
	return nil
}

func (s *CountWriter) Setup(config kiprotect.Config) error {
	return nil
}

func (s *CountWriter) LastItems() []map[string]interface{} {
	return s.lastItems
}

func (s *CountWriter) Count() int64 {
	return s.count
}

func (s *CountWriter) Write(payload kiprotect.Payload) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	var from int

	items := payload.Items()

	if len(items) > s.max {
		from = len(items) - s.max
	} else {
		from = 0
	}

	s.count += int64(len(items))

	for _, item := range items[from:] {
		s.lastItems = append(s.lastItems, item.All())
	}

	if len(s.lastItems) > s.max {
		// we truncate the lastItems array
		s.lastItems = s.lastItems[:s.max]
	}

	return nil
}

func MakeCountWriter(config map[string]interface{}) (kiprotect.Writer, error) {
	if params, err := CountWriterForm.Validate(config); err != nil {
		return nil, err
	} else {
		return &CountWriter{
			count:     0,
			max:       int(params["max"].(int64)),
			lastItems: make([]map[string]interface{}, 0),
			mutex:     &sync.Mutex{},
		}, nil
	}
}
