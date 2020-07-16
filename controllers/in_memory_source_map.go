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

package controllers

import (
	"fmt"
	"github.com/kiprotect/kiprotect"
	"time"
)

type InMemorySourceMap struct {
	kiprotect.BaseSourceMap
	name    string
	status  kiprotect.SourceStatus
	session interface{}
	source  *InMemorySource
	stream  *InMemoryStream
	id      []byte
}

func MakeInMemorySourceMap(id []byte, stream *InMemoryStream, source *InMemorySource, status kiprotect.SourceStatus) *InMemorySourceMap {
	sourceMap := &InMemorySourceMap{
		id:            id,
		source:        source,
		stream:        stream,
		status:        status,
		BaseSourceMap: kiprotect.BaseSourceMap{},
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

func (i *InMemorySourceMap) Source() kiprotect.Source {
	return i.source
}

func (i *InMemorySourceMap) Stream() kiprotect.Stream {
	return i.stream
}

func (i *InMemorySourceMap) Status() kiprotect.SourceStatus {
	return i.status
}

func (i *InMemorySourceMap) SetStatus(status kiprotect.SourceStatus) error {
	i.status = status
	return nil
}

func (i *InMemorySourceMap) SetStream(stream kiprotect.Stream) error {
	inMemoryStream, ok := stream.(*InMemoryStream)
	if !ok {
		return fmt.Errorf("not a inMemory stream")
	}
	i.stream = inMemoryStream
	return nil
}

func (i *InMemorySourceMap) SetSource(source kiprotect.Source) error {
	inMemorySource, ok := source.(*InMemorySource)
	if !ok {
		return fmt.Errorf("not a inMemory source")
	}
	i.source = inMemorySource
	return nil
}

func (i *InMemorySourceMap) CreatedAt() time.Time {
	return time.Now()
}

func (i *InMemorySourceMap) DeletedAt() *time.Time {
	return nil
}

func (i *InMemorySourceMap) UpdatedAt() time.Time {
	return time.Now()
}

func (i *InMemorySourceMap) Save() error {
	return nil
}

func (i *InMemorySourceMap) Refresh() error {
	return nil
}
