// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/kiprotect/kodex"
	"time"
)

var tws = []kodex.TimeWindowFunc{
	kodex.Minute,
	kodex.Hour,
	kodex.Day,
	kodex.Week,
	kodex.Month,
}

func Meter(meter kodex.ModelMeter, c *gin.Context, model kodex.Model, name string, data map[string]string, count int64) error {
	meterId := meter.ModelID(model)
	now := time.Now().UTC().UnixNano()
	for _, twt := range tws {

		tw := twt(now)
		if err := meter.Add(meterId, name, data, tw, count); err != nil {
			return err
		}
	}

	return nil
}
