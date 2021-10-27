// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

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
