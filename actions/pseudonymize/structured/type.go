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

package structured

type Type interface {
	IsValid() bool
	Decode(*BitArray) error
	Encode() (*BitArray, error)
	Length() uint
	Copy() Type
	Equals(Type) bool
}

type CompositeType interface {

	// helper functions to copy, check equality and validity etc...
	IsValid() []bool
	Equals(CompositeType) bool
	Copy() CompositeType

	// get info about subtypes
	Offset(int) (uint, error)
	Length(int) (uint, error)
	Get(int) (Type, error)

	// marshal and unmarshal to/from a given format
	Marshal(format string) (interface{}, error)
	Unmarshal(format string, value interface{}) error

	// encode and decode a composite type to/from bitarrays
	Decode(*BitArray) error
	Encode() (*BitArray, error)
	EncodeSubtype(int) (*BitArray, error)
	DecodeSubtype(int, *BitArray) error
}
