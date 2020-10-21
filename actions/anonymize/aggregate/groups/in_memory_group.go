package groups

import (
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
	"sync"
)

type InMemoryGroup struct {
	mutex         sync.RWMutex
	state         aggregate.State
	shard         *InMemoryShard
	groupByValues map[string]interface{}
	hash          []byte
	expiration    int64
}

func MakeInMemoryGroup(hash []byte,
	groupByValues map[string]interface{},
	expiration int64,
	shard *InMemoryShard) *InMemoryGroup {
	return &InMemoryGroup{
		hash:          hash,
		shard:         shard,
		groupByValues: groupByValues,
		expiration:    expiration,
	}
}

func (g *InMemoryGroup) Lock() {
	g.mutex.Lock()
}

func (g *InMemoryGroup) Unlock() {
	g.mutex.Unlock()
}

func (g *InMemoryGroup) RLock() {
	g.mutex.RLock()
}

func (g *InMemoryGroup) RUnlock() {
	g.mutex.RUnlock()
}

func (g *InMemoryGroup) Clone() (aggregate.Group, error) {
	clonedState, err := g.state.Clone()
	if err != nil {
		return nil, err
	}
	return &InMemoryGroup{
		state:         clonedState,
		shard:         g.shard,
		groupByValues: g.groupByValues,
		hash:          g.hash,
		expiration:    g.expiration,
	}, nil
}

// Returns whether a given group is initialized
func (g *InMemoryGroup) Initialized() bool {
	return g.state == nil
}

// Initialize the group
func (g *InMemoryGroup) Initialize(state aggregate.State) error {
	if g.state != nil {
		return aggregate.AlreadyInitialized
	}
	g.state = state
	return nil
}

// Return the state of the group
func (g *InMemoryGroup) State() aggregate.State {
	return g.state
}

// Return the group-by fields of a given group
func (g *InMemoryGroup) GroupByValues() map[string]interface{} {
	return g.groupByValues
}

// Get the expiration value for the group
func (g *InMemoryGroup) Expiration() int64 {
	return g.expiration
}

// Get the hash for the group
func (g *InMemoryGroup) Hash() []byte {
	return g.hash
}
