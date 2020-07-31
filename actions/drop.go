// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
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

package actions

import (
	"github.com/kiprotect/kodex"
)

type DropAction struct {
	kodex.BaseAction
}

func MakeDropAction(name, description string, id []byte, config map[string]interface{}) (kodex.Action, error) {
	if baseAction, err := kodex.MakeBaseAction(name, description, "drop", id, config); err != nil {
		return nil, err
	} else {
		return &DropAction{
			BaseAction: baseAction,
		}, nil
	}
}

func (a *DropAction) Params() interface{} {
	return nil
}

func (a *DropAction) GenerateParams(key, salt []byte) error {
	return nil
}

func (a *DropAction) SetParams(params interface{}) error {
	return nil
}

func (a *DropAction) Do(item *kodex.Item, writer kodex.ChannelWriter) (*kodex.Item, error) {
	return nil, nil
}

func (a *DropAction) Setup() error {
	return nil
}

func (a *DropAction) Teardown() error {
	return nil
}
