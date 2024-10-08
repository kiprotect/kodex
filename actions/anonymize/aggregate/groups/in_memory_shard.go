// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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
	id           int
	returned     bool
	store        *InMemoryGroupStore
	groupsByHash map[string]aggregate.Group
	groupsByTo   *MutexBTree
	hashMutex    sync.RWMutex
}

type MutexConfigGroup struct {
	Mutex sync.RWMutex
	Map   map[string]aggregate.Group
}

type MutexBTree struct {
	Mutex sync.RWMutex
	Tree  *btree.BTree
}

type GroupsByTo struct {
	To     int64
	Groups map[string]aggregate.Group
}

func (g *GroupsByTo) Less(b btree.Item) bool {
	gb, ok := b.(*GroupsByTo)
	if !ok {
		return false
	}
	return g.To < gb.To
}

// Create a new InMemoryShard object
func MakeInMemoryShard(id int, store *InMemoryGroupStore) *InMemoryShard {
	return &InMemoryShard{
		groupsByHash: make(map[string]aggregate.Group),
		groupsByTo:   &MutexBTree{Tree: btree.New(2)},
		store:        store,
		id:           id,
	}
}

// Return a group based on its hash
func (g *InMemoryShard) GroupByHash(hash []byte) (aggregate.Group, error) {

	g.hashMutex.RLock()
	group, ok := g.groupsByHash[string(hash)]
	g.hashMutex.RUnlock()

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

func (g *InMemoryShard) deleteGroupFromToMap(group aggregate.Group) {
	h := string(group.Hash())

	gr := &GroupsByTo{
		To:     group.Expiration(),
		Groups: nil,
	}

	g.groupsByTo.Mutex.Lock()
	if gb := g.groupsByTo.Tree.Get(gr); gb != nil {
		gbr, ok := gb.(*GroupsByTo)
		if !ok {
			panic("should not happen")
		}
		delete(gbr.Groups, h)
	}
	g.groupsByTo.Mutex.Unlock()
}

func (g *InMemoryShard) deleteGroupFromMap(group aggregate.Group) {
	h := string(group.Hash())
	g.hashMutex.Lock()
	delete(g.groupsByHash, h)
	g.hashMutex.Unlock()
}

func (g *InMemoryShard) ID() interface{} {
	return g.id
}

// Return a group based on its hash.
func (g *InMemoryShard) CreateGroup(hash []byte,
	groupByFields map[string]interface{}, expiration int64) (aggregate.Group, error) {
	h := string(hash)

	group := MakeInMemoryGroup(hash, groupByFields, expiration, g)

	g.hashMutex.Lock()
	g.groupsByHash[h] = group
	g.hashMutex.Unlock()

	gbt := &GroupsByTo{
		To:     expiration,
		Groups: map[string]aggregate.Group{h: group},
	}

	g.groupsByTo.Mutex.Lock()
	if egbt := g.groupsByTo.Tree.Get(gbt); egbt != nil {
		egbtr, ok := egbt.(*GroupsByTo)
		if !ok {
			panic("should not happen")
		}
		egbtr.Groups[h] = group
	} else {
		g.groupsByTo.Tree.ReplaceOrInsert(gbt)
	}
	g.groupsByTo.Mutex.Unlock()
	return group, nil
}

func (g *InMemoryShard) ExpireAllGroups() ([]aggregate.Group, error) {
	expiredGroups := make([]aggregate.Group, 0)
	g.hashMutex.Lock()
	defer g.hashMutex.Unlock()
	for h, group := range g.groupsByHash {
		delete(g.groupsByHash, h)
		// delete the group from the map
		g.deleteGroupFromToMap(group)
		expiredGroups = append(expiredGroups, group)
	}
	return expiredGroups, nil
}

// Get all expired groups for a given config hash and time in the shard
func (g *InMemoryShard) ExpireGroups(expiration int64) ([]aggregate.Group, error) {

	expiredGroups := make([]aggregate.Group, 0)

	g.groupsByTo.Mutex.RLock()
	min := g.groupsByTo.Tree.Min()
	g.groupsByTo.Mutex.RUnlock()

	if min == nil {
		return expiredGroups, nil
	}

	mg := min.(*GroupsByTo)

	if mg.To >= expiration {
		return expiredGroups, nil
	}

	g.groupsByTo.Mutex.Lock()
	defer g.groupsByTo.Mutex.Unlock()

	deleteList := make([]btree.Item, 0)

	iterator := func(i btree.Item) bool {
		gi, ok := i.(*GroupsByTo)
		if !ok || gi.To >= expiration {
			panic("should not happen")
		}
		for key, group := range gi.Groups {
			g.deleteGroupFromMap(group)
			delete(gi.Groups, key)
			if len(gi.Groups) == 0 {
				deleteList = append(deleteList, i)
			}
			expiredGroups = append(expiredGroups, group)
		}
		return true
	}

	gi := &GroupsByTo{
		To:     expiration,
		Groups: nil,
	}

	g.groupsByTo.Tree.AscendLessThan(gi, iterator)

	for _, i := range deleteList {
		g.groupsByTo.Tree.Delete(i)
	}

	return expiredGroups, nil
}
