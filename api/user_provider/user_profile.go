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

package provider

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
	"sync"
	"time"
)

type ProfileWithRetrievalTime struct {
	Profile     api.UserProfile
	Error       error
	RetrievedAt time.Time
	Refreshing  bool
	mutex       *sync.Mutex
}

type UserProfileProviderMaker func(kodex.Settings) (UserProfileProvider, error)

type UserProfileProvider interface {
	Get(string) (api.UserProfile, error)
	Start()
	Stop()
}

var providers = map[string]UserProfileProviderMaker{}

func MakeUserProfileProvider(settings kodex.Settings) (UserProfileProvider, error) {
	providerType, ok := settings.String("user-profile-provider.type")
	if !ok {
		return nil, fmt.Errorf("Provider type config missing (user-profile-provider.type)")
	}

	maker, ok := providers[providerType]

	if !ok {
		return nil, fmt.Errorf("Unknown provider type: %s", providerType)
	}

	return maker(settings)
}
