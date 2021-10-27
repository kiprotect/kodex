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
