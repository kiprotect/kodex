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

package anonymize

import (
	"encoding/base64"
	"github.com/kiprotect/go-helpers/errors"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate/filter_functions"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate/group_by_functions"
	"github.com/kiprotect/kodex/actions/anonymize/aggregate/groups"
	"sync"
)

type AggregateAnonymizer struct {
	channels             []string
	resultName           string
	function             Function
	finalizeAfter        int64
	id                   []byte
	name                 string
	filterFunctions      []filterFunctions.FilterFunction
	groupByFunctions     []groupByFunctions.GroupByFunction
	alwaysIncludedGroups int
	groupStore           aggregate.GroupStore
	mutex                sync.Mutex
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
		alwaysIncludedGroups := 0
		gbf := make([]groupByFunctions.GroupByFunction, 0)
		for _, groupByParams := range params["group-by"].([]interface{}) {
			groupByParamsMap := groupByParams.(map[string]interface{})
			functionConfig := groupByParamsMap["config"].(map[string]interface{})
			functionName := groupByParamsMap["function"].(string)
			if functionMaker, ok := groupByFunctions.Functions[functionName]; !ok {
				panic("should never happen")
			} else if groupByFunction, err := functionMaker(functionConfig); err != nil {
				return nil, err
			} else {
				alwaysIncluded := groupByParamsMap["always-included"].(bool)
				if alwaysIncluded {
					alwaysIncludedGroups++
					// we prepend the group so that the always included groups always come first
					gbf = append([]groupByFunctions.GroupByFunction{groupByFunction}, gbf...)
				} else {
					gbf = append(gbf, groupByFunction)
				}
			}
		}
		resultName, ok := params["result-name"].(string)
		if !ok {
			resultName = name
		}
		return &AggregateAnonymizer{
			function:             params["function"].(Function),
			channels:             params["channels"].([]string),
			finalizeAfter:        params["finalize-after"].(int64),
			groupByFunctions:     gbf,
			alwaysIncludedGroups: alwaysIncludedGroups,
			resultName:           resultName,
			name:                 name,
			id:                   id,
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

func (a *AggregateAnonymizer) Reset() error {
	return a.groupStore.Reset()
}

func (a *AggregateAnonymizer) Advance(writer kodex.ChannelWriter) ([]*kodex.Item, error) {
	kodex.Log.Info("Advancing aggregate anonymizer...")
	return nil, nil
}

func (a *AggregateAnonymizer) Finalize(writer kodex.ChannelWriter) ([]*kodex.Item, error) {
	if items, err := a.finalizeAllGroups(); err != nil {
		return nil, err
	} else {
		if err := a.submitResults(items, writer); err != nil {
			return nil, err
		}
	}
	return nil, nil
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (a *AggregateAnonymizer) getGroupByValues(item *kodex.Item) ([]*groupByFunctions.GroupByValue, error) {
	/*
		We calculate the power set of all group-by function values. Each group-by function can
		produce one or more values. We
	*/
	// Returns all unique group by value combinations of the given item
	groupByValues := make([][]*groupByFunctions.GroupByValue, 0, len(a.groupByFunctions))
	for _, groupByFunction := range a.groupByFunctions {
		if functionGroupByValues, err := groupByFunction(item); err != nil {
			return nil, err
		} else if functionGroupByValues != nil && len(functionGroupByValues) > 0 {
			groupByValues = append(groupByValues, functionGroupByValues)
		}
	}

	combinedGroupByValues := make([]*groupByFunctions.GroupByValue, 0)
	// we generate all combinations from 1 to n elements (up to a maximum number of combinations)
	for n := min(1+max(0, a.alwaysIncludedGroups-1), len(groupByValues)); n <= len(groupByValues); n++ {
		// the indices of the group-by functions that we add to the given combined group
		groupIndices := make([]int, n)
		// we specific groups within each group-by function that we add to the combined group
		elementIndices := make([]int, n)
		for i := 0; i < n; i++ {
			// we start with the first groups
			groupIndices[i] = i
			// we start with the first element in each group
			elementIndices[i] = 0
		}
		for {
			combinedGroupByValue := &groupByFunctions.GroupByValue{
				Expiration: 0,
				Values:     make(map[string]interface{}),
			}
			for i := 0; i < n; i++ {
				groupByValue := groupByValues[groupIndices[i]][elementIndices[i]]
				for k, v := range groupByValue.Values {
					if existingValue, ok := combinedGroupByValue.Values[k]; ok {
						// if there's an existing value for the given key already,
						// we append the new value to it
						if l, ok := existingValue.([]interface{}); ok {
							combinedGroupByValue.Values[k] = append(l, v)
						} else {
							combinedGroupByValue.Values[k] = []interface{}{existingValue, v}
						}
					} else {
						combinedGroupByValue.Values[k] = v
					}
				}
				// we update the expiration value
				if combinedGroupByValue.Expiration == 0 || groupByValue.Expiration > combinedGroupByValue.Expiration {
					combinedGroupByValue.Expiration = groupByValue.Expiration
				}
			}
			combinedGroupByValues = append(combinedGroupByValues, combinedGroupByValue)
			found := false
			// now we increase the index
			for i := n - 1; i >= 0; i-- {
				// we can still increase the value within this group
				if elementIndices[i] < len(groupByValues[groupIndices[i]])-1 {
					elementIndices[i] += 1
					found = true
					// we set all larger indices to 0
					for j := i + 1; j < n; j++ {
						elementIndices[j] = 0
					}
					break
				}
			}
			if !found {
				// we can't increase any individual index, so we need to increase
				// the group instead. We will not increase the groups that always
				// need to be included
				for i := n - 1; i >= a.alwaysIncludedGroups; i-- {
					if groupIndices[i] < len(groupByValues)-n+i {
						groupIndices[i] += 1
						found = true
						// we reset all element indices now as this is a new
						// combination of groups
						for j := 0; j < n; j++ {
							elementIndices[j] = 0
						}
						// we reset the higher group indices
						for j := i + 1; j < n; j++ {
							groupIndices[j] = groupIndices[i] + j - i
						}
						break
					}
				}
			}
			if !found {
				// we've exhausted all combinations of groups for this given
				// value n, so we increase n instead and continue
				break
			}
			// we stop at 100 groups
			if len(combinedGroupByValues) >= 100 {
				break
			}
		}
		// we stop at 100 groups
		if len(combinedGroupByValues) >= 100 {
			break
		}
	}
	return combinedGroupByValues, nil
}

func (a *AggregateAnonymizer) getGroups(item *kodex.Item, function aggregate.Function, shard aggregate.Shard) ([]aggregate.Group, error) {
	groupByValuesList, err := a.getGroupByValues(item)
	if err != nil {
		return nil, errors.MakeExternalError("error getting group-by values",
			"GET-GROUP-BY-VALUES",
			nil,
			err)
	}

	itemGroups := make([]aggregate.Group, 0)
	for _, groupByValue := range groupByValuesList {
		hash, err := kodex.StructuredHash(groupByValue.Values)
		if err != nil {
			return nil, err
		}
		group, err := shard.GroupByHash(hash)
		if err != nil && err != aggregate.NotFound {
			return nil, err
		}
		if group == nil {
			group, err = shard.CreateGroup(hash, groupByValue.Values, groupByValue.Expiration)
			if err != nil {
				return nil, err
			}
			// we initialize the group
			if err := function.Initialize(group); err != nil {
				return nil, err
			}
		}
		itemGroups = append(itemGroups, group)
	}
	return itemGroups, nil
}

func (a *AggregateAnonymizer) finalizeAllGroups() ([]*kodex.Item, error) {
	allGroups, err := a.groupStore.ExpireAllGroups()
	if err != nil {
		return nil, err
	}
	return a.finalizeGroups(allGroups)
}

func (a *AggregateAnonymizer) getMinimumExpiration(groups []aggregate.Group) int64 {
	var exp int64 = -1
	for _, group := range groups {
		if exp < 0 || group.Expiration() < exp {
			exp = group.Expiration()
		}
	}
	return exp
}

func (a *AggregateAnonymizer) finalizeExpiredGroups(shard aggregate.Shard, expiration int64) ([]*kodex.Item, error) {
	if a.finalizeAfter == -1 {
		return nil, nil
	}
	expiredGroups, err := a.groupStore.ExpireGroups(expiration)
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
		aggregations, err := a.finalizeExpiredGroups(shard, a.getMinimumExpiration(groups))
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
