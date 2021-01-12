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

import (
	"github.com/kiprotect/kodex/actions/pseudonymize/merengue"
)

func PS(c CompositeType, key []byte) (CompositeType, error) {
	cc := c.Copy()
	pv, err := cc.Encode()
	if err != nil {
		return nil, err
	}
	i := 0
	for {
		i += 1
		pvb := merengue.PseudonymizeBidirectional(pv.Bytes(), pv.Length(), key, key, merengue.Sha256)
		pv, err = MakeBitArrayFromBytes(pvb, pv.Length())
		if err != nil {
			return nil, err
		}
		if err = cc.Decode(pv); err != nil {
			return nil, err
		}
		valid := cc.IsValid()
		found := false
		for _, v := range valid {
			if !v {
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return cc, nil
}

func DPS(c CompositeType, key []byte) (CompositeType, error) {
	var err error
	cc := c.Copy()
	pv, err := cc.Encode()
	if err != nil {
		return nil, err
	}
	i := 0
	for {
		i += 1
		pvb := merengue.DepseudonymizeBidirectional(pv.Bytes(), pv.Length(), key, key, merengue.Sha256)
		pv, err = MakeBitArrayFromBytes(pvb, pv.Length())
		if err != nil {
			return nil, err
		}
		if err = cc.Decode(pv); err != nil {
			return nil, err
		}
		valid := cc.IsValid()
		found := false
		for _, v := range valid {
			if !v {
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return cc, nil
}

func PSH(c CompositeType, key []byte) (CompositeType, error) {

	cc := c.Copy()

	ba, err := cc.Encode()

	if err != nil {
		return nil, err
	}

	var pvs *BitArray
	for {
		pvsb := merengue.Pseudonymize(ba.Bytes(), ba.Length(), key, merengue.Sha256)
		pvs, err = MakeBitArrayFromBytes(pvsb, ba.Length())
		if err != nil {
			return nil, err
		}
		if err = cc.Decode(pvs); err != nil {
			return nil, err
		}
		valid := cc.IsValid()
		found := false
		for i, v := range valid {
			if !v {

				off, err := cc.Offset(i)
				if err != nil {
					return nil, err
				}

				length, err := cc.Length(i)
				if err != nil {
					return nil, err
				}

				pvi, err := pvs.Extract(off, length)
				if err != nil {
					return nil, err
				}
				ba.Update(pvi, off)
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return cc, nil
}

func DPSH(c CompositeType, key []byte) (CompositeType, error) {

	cc := c.Copy()

	ba, err := cc.Encode()

	if err != nil {
		return nil, err
	}

	pvsb := merengue.Depseudonymize(ba.Bytes(), ba.Length(), key, merengue.Sha256)
	pvs, err := MakeBitArrayFromBytes(pvsb, ba.Length())

	if err != nil {
		return nil, err
	}

	for {

		if err = cc.Decode(pvs); err != nil {
			return nil, err
		}

		valid := cc.IsValid()

		found := false
		for i := len(valid) - 1; i >= 0; i-- {
			if !valid[i] {
				off, err := cc.Offset(i)
				if err != nil {
					return nil, err
				}
				length, err := cc.Length(i)
				if err != nil {
					return nil, err
				}
				pvi, err := pvs.Extract(off, length)
				if err != nil {
					return nil, err
				}
				pvk := ba.Copy()
				pvk.Update(pvi, off)
				pvsdb := merengue.Depseudonymize(pvk.Bytes(), pvk.Length(), key, merengue.Sha256)
				pvsd, err := MakeBitArrayFromBytes(pvsdb, pvk.Length())
				if err != nil {
					return nil, err
				}
				pvid, err := pvsd.Extract(off, length)
				if err != nil {
					return nil, err
				}
				pvs.Update(pvid, off)
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	return cc, nil

}
