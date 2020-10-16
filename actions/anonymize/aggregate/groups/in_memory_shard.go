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
	"github.com/google/btree"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
	"sync"
)

type InMemoryShard struct {
	id                int
	returned          bool
	groupsByHash      map[string]aggregate.Group
	numericalTriggers map[string]*NumericalTrigger
	valueTriggers     map[string]*ValueTrigger
	deletedGroups     map[string]int64
	store             *InMemoryGroupStore
	deleteMutex       sync.RWMutex
	hashMutex         sync.RWMutex
	toMutex           sync.RWMutex
	minToMutex        sync.RWMutex
}

type ValueTrigger struct {
	Mutex  sync.RWMutex
	Values map[string]map[string]aggregate.Group
}

type NumericalTrigger struct {
	Mutex sync.RWMutex
	Tree  *btree.BTree
}

type NumericalValue struct {
	Value  int64
	Groups map[string]aggregate.Group
}

func (g *NumericalValue) Less(b btree.Item) bool {
	gb, ok := b.(*NumericalValue)
	if !ok {
		return false
	}
	return g.Value < gb.Value
}

// Create a new InMemoryShard object
func MakeInMemoryShard(id int, store *InMemoryGroupStore) *InMemoryShard {
	return &InMemoryShard{
		groupsByHash: make(map[string]aggregate.Group),
		store:        store,
		id:           id,
	}
}

// Return a group based on its hash
func (g *InMemoryShard) GroupByHash(hash []byte) (aggregate.Group, error) {

	g.deleteMutex.Lock()
	_, ok := g.deletedGroups[string(hash)]
	g.deleteMutex.Unlock()

	if ok {
		return nil, aggregate.AlreadyDeleted
	}

	g.hashMutex.Lock()
	group, ok := g.groupsByHash[string(hash)]
	g.hashMutex.Unlock()

	if !ok {
		return nil, aggregate.NotFound
	} else {
		return group, nil
	}
}

// Commit the state of the shard to the store
func (g *InMemoryShard) Commit() error {
	/*
		Committing a shard incurs no extra cost for the in-memory shard.
	*/
	return nil
}

// Return the shard to the store
func (g *InMemoryShard) Return() error {
	g.returned = true
	return g.store.Return(g.id)
}

func (g *InMemoryShard) Groups() (
	[]aggregate.Group, error) {
	groups := make([]aggregate.Group, 0, len(g.groupsByHash))
	for _, group := range g.groupsByHash {
		groups = append(groups, group)
	}
	return groups, nil
}

func (g *InMemoryShard) CreateGroup(hash []byte, groupByValues map[string]interface{}, triggers []*aggregate.Trigger) (aggregate.Group, error) {
	return nil, nil
}

// Return all expired groups from all shards for a given config hash,
// grouped by full hash.
func (g *InMemoryShard) ExpireGroups(triggers []*aggregate.Trigger) (map[string][]aggregate.Group, error) {
	return nil, nil
}

func (g *InMemoryShard) ExpireAllGroups() (map[string][]aggregate.Group, error) {
	return nil, nil
}

func (g *InMemoryShard) ID() interface{} {
	return g.id
}

func (g *InMemoryShard) FinalizeExpiredGroups(time int64) (map[string]aggregate.Group, error) {
	return nil, nil
	/*
		expiredGroups := make(map[string]aggregate.Group)
		g.toMutex.Lock()
		bt, ok := g.groupsByTo[string(hash)]
		g.toMutex.Unlock()
		if ok {
			bt.Mutex.Lock()
			min := bt.Tree.Min()
			bt.Mutex.Unlock()
			if min == nil {
				return expiredGroups, nil
			}
			mg := min.(*GroupsByTo)
			if mg.To >= time {
				return expiredGroups, nil
			}
			bt.Mutex.Lock()
			defer bt.Mutex.Unlock()
			deleteList := make([]btree.Item, 0)
			iterator := func(i btree.Item) bool {
				gi, ok := i.(*GroupsByTo)
				if !ok || gi.To >= time {
					panic("should not happen")
				}
				for key, group := range gi.Groups {
					sh := string(group.Hash())
					if err := g.markGroupAsDeleted(group); err != nil {
						expiredGroups[sh] = nil
						continue
					}
					g.deleteGroupFromConfigMap(group)
					g.deleteGroupFromFullMap(group)
					delete(gi.Groups, key)
					if len(gi.Groups) == 0 {
						deleteList = append(deleteList, i)
					}
					expiredGroups[sh] = group
				}
				return true
			}
			gi := &GroupsByTo{
				To:     time,
				Groups: nil,
			}
			bt.Tree.AscendLessThan(gi, iterator)
			for _, i := range deleteList {
				bt.Tree.Delete(i)
			}
		}
		return expiredGroups, nil
	*/
}
