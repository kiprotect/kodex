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

package anonymize

import (
	"encoding/base64"
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate/groups"
	"sync"
)

type AggregateAnonymizer struct {
	channels       []string
	resultName     string
	function       Function
	id             []byte
	name           string
	groupByClauses []aggregate.GroupByClause
	groupByConfig  []map[string]interface{}
	groupStore     aggregate.GroupStore
	mutex          sync.Mutex
}

func (a *AggregateAnonymizer) Setup() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var err error
	if a.groupStore, err = groups.GroupStores["in-memory"](a.id); err != nil {
		return errors.MakeExternalError("in-memory store not defined", "IN-MEMORY-STORE", nil, err)
	}
	return nil
}

func (a *AggregateAnonymizer) Teardown() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return nil
}

func MakeAggregateAnonymizer(name string, id []byte, config map[string]interface{}) (Anonymizer, error) {
	if params, err := AggregateForm.Validate(config); err != nil {
		return nil, err
	} else {
		resultName, ok := params["result-name"].(string)
		if !ok {
			resultName = name
		}
		return &AggregateAnonymizer{
			function:      params["function"].(Function),
			groupByConfig: params["group-by"].([]map[string]interface{}),
			channels:      params["channels"].([]string),
			resultName:    resultName,
			name:          name,
			id:            id,
		}, nil
	}
}

func (a *AggregateAnonymizer) Params() interface{} {
	return nil
}

func (a *AggregateAnonymizer) GenerateParams(key, salt []byte) error {
	return nil
}

func (a *AggregateAnonymizer) SetParams(params interface{}) error {
	return nil
}

func (a *AggregateAnonymizer) Anonymize(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	return a.process(item, writer)
}

func (a *AggregateAnonymizer) Finalize() ([]*kodex.Item, error) {
	shard, err := a.groupStore.Shard()
	if err != nil {
		return nil, errors.MakeExternalError("cannot get a shard", "IN-MEMORY-STORE", nil, err)
	}
	defer shard.Return()
	return a.finalizeAllGroups(shard)
}

func (a *AggregateAnonymizer) process(item *kodex.Item, channelWriter kodex.ChannelWriter) (*kodex.Item, error) {
	shard, err := a.groupStore.Shard()
	if err != nil {
		return nil, errors.MakeExternalError("cannot get a shard", "IN-MEMORY-STORE", nil, err)
	}
	defer shard.Return()
	if err := a.aggregate(item, channelWriter, shard); err != nil {
		return nil, err
	}
	return item, nil
}

func (a *AggregateAnonymizer) getGroupByValues(item *kodex.Item) ([]map[string]interface{}, error) {
	return nil, nil
}

func (a *AggregateAnonymizer) getTriggers(item *kodex.Item) ([]*aggregate.Trigger, error) {
	return nil, nil
}

func (a *AggregateAnonymizer) getGroups(item *kodex.Item, function aggregate.Function, shard aggregate.Shard) ([]aggregate.Group, error) {
	groupByValuesList, err := a.getGroupByValues(item)
	if err != nil {
		return nil, errors.MakeExternalError("error getting group-by values",
			"GET-GROUP-BY-VALUES",
			nil,
			err)
	}
	triggers, err := a.getTriggers(item)
	if err != nil {
		return nil, err
	}
	groups := make([]aggregate.Group, 0)
	for _, groupByValues := range groupByValuesList {
		hash, err := kodex.StructuredHash(groupByValues)
		if err != nil {
			return nil, err
		}
		group, err := shard.GroupByHash(hash)
		if err != nil {
			return nil, err
		}
		if group == nil {
			group, err = shard.CreateGroup(hash, groupByValues, triggers)
			if err != nil {
				return nil, err
			}
			// we initialize the group
			if err := function.Initialize(group); err != nil {
				return nil, err
			}
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (a *AggregateAnonymizer) finalizeAllGroups(shard aggregate.Shard) ([]*kodex.Item, error) {
	allGroups, err := shard.ExpireAllGroups()
	if err != nil {
		return nil, err
	}
	return a.finalizeGroups(allGroups)
}

func (a *AggregateAnonymizer) finalizeExpiredGroups(shard aggregate.Shard, triggers []*aggregate.Trigger) ([]*kodex.Item, error) {
	expiredGroups, err := shard.ExpireGroups(triggers)
	if err != nil {
		return nil, err
	}
	return a.finalizeGroups(expiredGroups)
}

func (a *AggregateAnonymizer) finalizeGroups(groups map[string][]aggregate.Group) ([]*kodex.Item, error) {

	encode := func(data []byte) string {
		return base64.StdEncoding.EncodeToString(data)
	}

	items := make([]*kodex.Item, 0)
	for _, hashGroups := range groups {
		group, err := a.function.Function.Merge(hashGroups)
		if err != nil {
			return items, err
		}
		result, err := a.function.Function.Finalize(group)
		if err != nil {
			return items, err
		}
		// we omit reporting null results (to protect privacy)
		if result == nil {
			continue
		}
		item := kodex.MakeItem(map[string]interface{}{
			a.resultName:  result,
			"action_id":   a.id,
			"action_name": a.name,
			"group":       group.GroupByValues(),
			"group_hash":  encode(group.Hash()),
		})
		items = append(items, item)
	}
	return items, nil
}

func (a *AggregateAnonymizer) submitResults(items []*kodex.Item, channelWriter kodex.ChannelWriter) error {
	for _, channel := range a.channels {
		if err := channelWriter.Write(channel, items); err != nil {
			return err
		}
	}
	return nil
}

func (a *AggregateAnonymizer) aggregate(item *kodex.Item, channelWriter kodex.ChannelWriter, shard aggregate.Shard) error {

	/*
		- Generate the groups for the items using the group-by clauses
		- Add the items to the groups
		- Update the finalization triggers
		- Finalize all groups
	*/

	// we retrieve or create the group for the given item
	groups, err := a.getGroups(item, a.function.Function, shard)

	if err != nil {
		return err
	}

	triggers, err := a.getTriggers(item)

	if err != nil {
		return err
	}

	var groupErr error
	// todo: it might be problematic if a single group action fails for an
	// item, as we do not want to retry it too often (as it will exhaust)
	// the privacy budget of the query. We will need some rollback or
	// transaction-based commit mechanism for this.
	for _, group := range groups {
		// we add the item to the group result using the function
		if err := a.function.Function.Add(item, group); err != nil {
			groupErr = err
			continue
		}

		// we finalize all expired groups and return their results
		aggregations, err := a.finalizeExpiredGroups(shard, triggers)
		// we submit the results to the configured destination configs
		if aggregations != nil && len(aggregations) > 0 {
			if err := a.submitResults(aggregations, channelWriter); err != nil {
				kodex.Log.Error(err)
				groupErr = err
				continue
			}
		}
		if err != nil {
			groupErr = err
		}
	}
	return groupErr
}
