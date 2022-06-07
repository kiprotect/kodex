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

package controllers

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"time"
)

type InMemorySourceMap struct {
	kodex.BaseSourceMap
	name    string
	prio    float64
	prioT   time.Time
	status  kodex.SourceStatus
	session interface{}
	source  *InMemorySource
	stream  *InMemoryStream
	id      []byte
}

func MakeInMemorySourceMap(id []byte, stream *InMemoryStream, source *InMemorySource, status kodex.SourceStatus) *InMemorySourceMap {
	sourceMap := &InMemorySourceMap{
		id:            id,
		source:        source,
		stream:        stream,
		status:        status,
		prioT:         time.Now().UTC(),
		BaseSourceMap: kodex.BaseSourceMap{},
	}
	sourceMap.Self = sourceMap
	return sourceMap
}

func (i *InMemorySourceMap) ID() []byte {
	return i.id
}

func (i *InMemorySourceMap) Delete() error {
	return nil
}

func (i *InMemorySourceMap) Session() interface{} {
	return i.session
}

func (i *InMemorySourceMap) SetSession(session interface{}) error {
	i.session = session
	return nil
}

func (i *InMemorySourceMap) Source() kodex.Source {
	return i.source
}

func (i *InMemorySourceMap) Stream() kodex.Stream {
	return i.stream
}

func (i *InMemorySourceMap) Status() kodex.SourceStatus {
	return i.status
}

func (i *InMemorySourceMap) SetStatus(status kodex.SourceStatus) error {
	i.status = status
	return nil
}

func (i *InMemorySourceMap) SetStream(stream kodex.Stream) error {
	inMemoryStream, ok := stream.(*InMemoryStream)
	if !ok {
		return fmt.Errorf("not a inMemory stream")
	}
	i.stream = inMemoryStream
	return nil
}

func (i *InMemorySourceMap) SetSource(source kodex.Source) error {
	inMemorySource, ok := source.(*InMemorySource)
	if !ok {
		return fmt.Errorf("not a inMemory source")
	}
	i.source = inMemorySource
	return nil
}

func (i *InMemorySourceMap) CreatedAt() time.Time {
	return time.Now().UTC()
}

func (i *InMemorySourceMap) DeletedAt() *time.Time {
	return nil
}

func (i *InMemorySourceMap) UpdatedAt() time.Time {
	return time.Now().UTC()
}

func (i *InMemorySourceMap) Save() error {
	return nil
}

func (i *InMemorySourceMap) Refresh() error {
	return nil
}

/* Priority Related Functionality */

func (i *InMemorySourceMap) SetPriority(value float64) error {
	i.prio = value
	return nil
}

func (i *InMemorySourceMap) Priority() float64 {
	return i.prio
}

func (i *InMemorySourceMap) PriorityTime() time.Time {
	return i.prioT
}

func (i *InMemorySourceMap) SetPriorityTime(t time.Time) error {
	i.prioT = t
	return nil
}

func (i *InMemorySourceMap) SetPriorityAndTime(value float64, t time.Time) error {
	i.prioT = t
	i.prio = value
	return nil
}
