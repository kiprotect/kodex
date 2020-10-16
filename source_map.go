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

import (
	"encoding/json"
	"fmt"
)

type SourceMap interface {
	Processable
	Source() Source
	Stream() Stream
	Status() SourceStatus
	SetStatus(SourceStatus) error
	// Return the current session (only applicable for batch sources)
	Session() interface{}
	// Update the current session (only applicable for batch sources)
	SetSession(interface{}) error
}

/* Base Functionality */

type BaseSourceMap struct {
	Self SourceMap
}

func (b *BaseSourceMap) Type() string {
	return "source_map"
}

func (b *BaseSourceMap) Update(values map[string]interface{}) error {
	return fmt.Errorf("BaseSourceMap.Update not implemented")
}

func (b *BaseSourceMap) Create(values map[string]interface{}) error {
	return fmt.Errorf("BaseSourceMap.Create not implemented")
}

func (b *BaseSourceMap) MarshalJSON() ([]byte, error) {

	data := map[string]interface{}{
		"status": b.Self.Status(),
		"source": b.Self.Source(),
		"stream": b.Self.Stream(),
	}

	for k, v := range JSONData(b.Self) {
		data[k] = v
	}

	return json.Marshal(data)
}
