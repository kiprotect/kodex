// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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

package aggregate

import (
	"github.com/kiprotect/go-helpers/errors"
)

type Group interface {
	Initialized() bool
	Initialize(state State) error
	State() State
	GroupByValues() map[string]interface{}
	Hash() []byte
	Expiration() int64
	Clone() (Group, error)
	Lock()
	Unlock()
}

type Shard interface {
	// A unique ID identifying the shard globally
	ID() interface{}
	// Synchronize the shard state with the backend
	Commit() error
	// Return the shard to the group store so that it can be reused
	Return() error
	// Create a group in the shard
	CreateGroup(hash []byte, groupByValues map[string]interface{}, expiration int64) (Group, error)
	// Return a group by its unique hash value
	GroupByHash(hash []byte) (Group, error)
	// Expire groups in the shard based on an expiration index
	ExpireGroups(expiration int64) ([]Group, error)
	// Expire all groups in the shard
	ExpireAllGroups() ([]Group, error)
}

type GroupStore interface {
	// Tear down the group store
	Teardown() error
	// Reset the entire state of the group store
	Reset() error
	// Get a shard from the store
	Shard() (Shard, error)
	// Expire groups in all shards based on an expiration index
	ExpireGroups(expiration int64) (map[string][]Group, error)
	// Expire all groups in all shards
	ExpireAllGroups() (map[string][]Group, error)
}

var AlreadyInitialized = errors.MakeExternalError("group has already been initialized", "GROUP-STORE", nil, nil)
var AlreadyFinalized = errors.MakeExternalError("group has already been finalized", "GROUP-STORE", nil, nil)
var NotFound = errors.MakeExternalError("group not found", "GROUP-STORE", nil, nil)
