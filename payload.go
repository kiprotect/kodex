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

package kiprotect

type Payload interface {
	Items() []*Item
	Headers() map[string]interface{}
	EndOfStream() bool
	Acknowledge() error
	Reject() error
}

type BasicPayload struct {
	items       []*Item
	headers     map[string]interface{}
	endOfStream bool
}

func MakeBasicPayload(items []*Item, headers map[string]interface{}, endOfStream bool) *BasicPayload {
	return &BasicPayload{
		items:       items,
		headers:     headers,
		endOfStream: endOfStream,
	}
}

func (b *BasicPayload) Items() []*Item {
	return b.items
}

func (b *BasicPayload) Headers() map[string]interface{} {
	return b.headers
}

func (b *BasicPayload) EndOfStream() bool {
	return b.endOfStream
}

func (b *BasicPayload) Acknowledge() error {
	return nil
}

func (b *BasicPayload) Reject() error {
	return nil
}
