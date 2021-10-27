// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package kodex

type Metric struct {
	Name       string
	TimeWindow TimeWindow
	Value      int64
	Data       map[string]string
}

type Meter interface {
	// Add the given value to the metric
	Add(id string, name string, data map[string]string, tw TimeWindow, value int64) error
	// Return the metric and its assigned quota
	Get(id string, name string, data map[string]string, tw TimeWindow) (*Metric, error)
	// Return metrics for a given ID and time interval
	Range(id string, from, to int64, name, twType string) ([]*Metric, error)
	N(id string, to int64, n int64, name, twType string) ([]*Metric, error)
}

type ModelMeter interface {
	Meter
	// Add a given metric to a model (both time-based metrics and totals)
	AddToModel(model Model, name string, tw TimeWindow, value int64) error
	ModelID(model Model) string
}
