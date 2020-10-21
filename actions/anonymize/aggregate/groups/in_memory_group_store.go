// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
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

package groups

import (
	"encoding/hex"
	"fmt"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
	"sync"
)

type InMemoryGroupStore struct {
	shards     map[int]*InMemoryShard
	usedShards map[int]bool
	id         []byte
	mutex      sync.RWMutex
	shardCount int
}

var Store map[string]*InMemoryGroupStore
var mutex sync.Mutex

// Create a new InMemoryGroupStore object for the given config
func MakeInMemoryGroupStore(id []byte) (aggregate.GroupStore, error) {
	if Store == nil {
		mutex.Lock()
		if Store == nil {
			Store = make(map[string]*InMemoryGroupStore)
		}
		mutex.Unlock()
	}
	strId := hex.EncodeToString(id)
	if Store[strId] == nil {
		mutex.Lock()
		if Store[strId] == nil {
			Store[strId] = &InMemoryGroupStore{
				shards:     make(map[int]*InMemoryShard),
				usedShards: make(map[int]bool),
				id:         id,
			}
		}
		mutex.Unlock()
	}
	return Store[strId], nil
}

func (g *InMemoryGroupStore) Return(id int) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	used, ok := g.usedShards[id]
	if !ok {
		return fmt.Errorf("Shard does not exist")
	}
	if !used {
		return fmt.Errorf("Shard is not used")
	}
	g.usedShards[id] = false
	return nil
}

// Reset the store
func (g *InMemoryGroupStore) Reset() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.shards = make(map[int]*InMemoryShard)
	g.usedShards = make(map[int]bool)
	return nil
}

func (g *InMemoryGroupStore) Shard() (aggregate.Shard, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	for i, shard := range g.shards {
		if used, ok := g.usedShards[i]; !ok || !used {
			g.usedShards[i] = true
			return shard, nil
		}
	}
	g.shardCount++
	newShard := MakeInMemoryShard(g.shardCount, g)
	g.shards[g.shardCount] = newShard
	g.usedShards[g.shardCount] = true
	return newShard, nil
}

func (g *InMemoryGroupStore) ExpireGroups(expiration int64) (map[string][]aggregate.Group, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return nil, nil
}

func (g *InMemoryGroupStore) ExpireAllGroups() (map[string][]aggregate.Group, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return nil, nil
}
