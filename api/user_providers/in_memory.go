package providers

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex"
	"github.com/kiprotect/kodex/api"
)

var InMemoryUserForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "displayName",
			Validators: []forms.Validator{
				forms.IsOptional{Default: ""},
				forms.IsString{},
			},
		},
	},
}

type InMemoryUserProviderSettings struct {
}

func ValidateInMemoryUserProviderSettings(settings map[string]interface{}) (interface{}, error) {
	if params, err := InMemoryUserForm.Validate(settings); err != nil {
		return nil, err
	} else {
		providerSettings := &InMemoryUserProviderSettings{}
		if err := InMemoryUserForm.Coerce(providerSettings, params); err != nil {
			return nil, err
		}
		return providerSettings, nil
	}
}

type InMemoryUserProvider struct {
	users []*api.User
}

func MakeInMemoryUserProvider(settings kodex.Settings) (api.UserProvider, error) {
	return &InMemoryUserProvider{
		users: make([]*api.User, 0),
	}, nil
}

func (i *InMemoryUserProvider) Create(user *api.User) error {
	i.users = append(i.users, user)
	return nil
}

// Return a user with the given access token
func (i *InMemoryUserProvider) Get(stringToken string) (*api.User, error) {
	if token, err := hex.DecodeString(stringToken); err != nil {
		return nil, fmt.Errorf("malformed access token")
	} else {
		for _, user := range i.users {
			if bytes.Equal(user.AccessToken.Token, token) {
				return user, nil
			}
		}
	}
	return nil, fmt.Errorf("invalid access token")
}
func (i *InMemoryUserProvider) Start() {
}

func (i *InMemoryUserProvider) Stop() {

}
