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

package functions

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
)

type Int64 struct {
	I int64
}

func (i *Int64) Serialize() ([]byte, error) {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, i.I)
	return buf, nil
}

func (i *Int64) Deserialize(buf []byte) error {
	var err error
	i.I, err = binary.ReadVarint(bytes.NewBuffer(buf))
	return err
}

func (i *Int64) Clone() (aggregate.State, error) {
	return &Int64{I: i.I}, nil
}

type StringBoolMap struct {
	M map[string]bool
}

func (m *StringBoolMap) Clone() (aggregate.State, error) {
	newMap := make(map[string]bool)
	for key, value := range m.M {
		newMap[key] = value
	}
	return &StringBoolMap{M: newMap}, nil
}

func (m *StringBoolMap) Serialize() ([]byte, error) {
	var o bytes.Buffer
	enc := gob.NewEncoder(&o)
	if err := enc.Encode(m.M); err != nil {
		return nil, err
	}
	return o.Bytes(), nil
}

func (m *StringBoolMap) Deserialize(buf []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(buf))
	return dec.Decode(&m.M)
}
