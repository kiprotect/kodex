// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

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
		return nil, fmt.Errorf("No meter type given")
	}
	switch meterType {
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
