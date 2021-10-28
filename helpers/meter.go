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

package helpers

import (
	"fmt"
	"github.com/kiprotect/go-helpers/maps"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/metering"
)

func Meter(settings kodex.Settings) (kodex.Meter, error) {
	if disabled, _ := settings.Bool("meter.disable"); disabled {
		return nil, nil
	}
	meterType, ok := settings.String("meter.type")
	if !ok {
		meterType = "in-memory"
	}
	switch meterType {
	case "in-memory":
		config, err := settings.Get("meter.config")
		if err != nil {
			config = map[string]interface{}{}
		}
		configMap, ok := maps.ToStringMap(config)
		if !ok {
			return nil, fmt.Errorf("not a string map")
		}
		meter, err := metering.MakeInMemoryMeter(configMap)
		if err != nil {
			return nil, err
		}
		return meter, nil
	case "redis":
		config, err := settings.Get("meter.config")
		if err != nil {
			return nil, err
		}
		configMap, ok := maps.ToStringMap(config)
		if !ok {
			return nil, fmt.Errorf("not a string map")
		}
		meter, err := metering.MakeRedisMeter(configMap)
		if err != nil {
			return nil, err
		}
		return meter, nil
	}
	return nil, fmt.Errorf("invalid meter type: %s", meterType)

}
