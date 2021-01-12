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

package kodex

import (
	"encoding/hex"
	"time"
)

type Model interface {
	ID() []byte
	Type() string
	Delete() error
	UpdatedAt() time.Time
	CreatedAt() time.Time
	DeletedAt() *time.Time
	Refresh() error
	Save() error
	Create(values map[string]interface{}) error
	Update(values map[string]interface{}) error
}

// A model that has an associated priority
type PriorityModel interface {
	AddToPriority(value float64) error
	Priority() float64
	PriorityTime() time.Time
	ResetPriority() error
}

// A model that allows storing/retrieving statistics
type StatsModel interface {
	Stats() (map[string]int64, error)
	Stat(string) (int64, error)
	// Set a given statistic
	SetStat(string, int64) error
	// Add to a given statistic
	AddToStat(string, int64) error
}

func JSONData(model Model) map[string]interface{} {
	return map[string]interface{}{
		"id":         hex.EncodeToString(model.ID()),
		"created_at": model.CreatedAt(),
		"updated_at": model.UpdatedAt(),
	}
}
