// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package metering

import (
	"encoding/hex"
	"fmt"
	"github.com/kiprotect/kodex"
)

type BaseModelMeter struct {
	Self kodex.Meter
}

func (r *BaseModelMeter) AddToModel(model kodex.Model,
	name string,
	tw kodex.TimeWindow,
	value int64) error {
	modelId := r.ModelID(model)
	statsModel, ok := model.(kodex.StatsModel)
	if ok {
		if err := statsModel.AddToStat(name, value); err != nil {
			return err
		}
	}
	return r.Self.Add(modelId, name, map[string]string{}, tw, value)
}

func (r *BaseModelMeter) ModelID(model kodex.Model) string {
	return fmt.Sprintf("%s:%s", model.Type(), hex.EncodeToString(model.ID()))
}
