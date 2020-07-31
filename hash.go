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
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"sort"
)

// Computes a hash of a structured data type that can contain various types
// like strings or []byte arrays. The hash reflects both the type values and
// the structure of the source.
func StructuredHash(source interface{}) ([]byte, error) {
	h := sha256.New()
	if err := addHash(source, h); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func sortedKeys(h map[string]interface{}) []string {
	s := make([]string, 0)
	for key, _ := range h {
		s = append(s, key)
	}
	sort.Strings(s)
	return s
}

func addHash(source interface{}, h hash.Hash) error {
	switch v := source.(type) {
	case []byte:
		_, err := h.Write(v)
		return err
	case []interface{}:
		for i, entry := range v {
			if err := addHash(i, h); err != nil {
				return err
			}
			if err := addHash(entry, h); err != nil {
				return err
			}
		}
	case []string:
		// we duplicate this code for sake of efficiency
		for i, entry := range v {
			if err := addHash(i, h); err != nil {
				return err
			}
			if err := addHash(entry, h); err != nil {
				return err
			}
		}
	case []int:
		// we duplicate this code for sake of efficiency
		for i, entry := range v {
			if err := addHash(i, h); err != nil {
				return err
			}
			if err := addHash(entry, h); err != nil {
				return err
			}
		}
	case []int64:
		// we duplicate this code for sake of efficiency
		for i, entry := range v {
			if err := addHash(i, h); err != nil {
				return err
			}
			if err := addHash(entry, h); err != nil {
				return err
			}
		}
	case map[string]interface{}:
		for _, key := range sortedKeys(v) {
			value := v[key]
			if err := addHash(key, h); err != nil {
				return err
			}
			if err := addHash(value, h); err != nil {
				return err
			}
		}
	case string:
		if _, err := h.Write([]byte(v)); err != nil {
			return err
		}
	case bool:
		if v {
			return addHash(1, h)
		}
		return addHash(0, h)
	case int:
		return addHash(int64(v), h)
	case int64:
		bs := make([]byte, binary.MaxVarintLen64)
		binary.PutVarint(bs, v)
		if _, err := h.Write(bs); err != nil {
			return err
		}
	case float64:
		bits := math.Float64bits(v)
		bs := make([]byte, binary.MaxVarintLen64)
		binary.LittleEndian.PutUint64(bs, bits)
		if _, err := h.Write(bs); err != nil {
			return err
		}
	case nil:
		h.Write([]byte("magic nil value"))
	default:
		return fmt.Errorf("unknown type, can't hash")
	}
	return nil
}
