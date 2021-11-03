package providers

import (
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
	users map[string]*api.User
}

func MakeInMemoryUserProvider(settings kodex.Settings) (api.UserProvider, error) {
	kodex.Log.Info("Making in-memory user provider")
	return &InMemoryUserProvider{}, nil
}

// Return a user with the given access token
func (i *InMemoryUserProvider) Get(string) (*api.User, error) {
	return nil, fmt.Errorf("access token missing")
}
func (i *InMemoryUserProvider) Start() {
	kodex.Log.Info("Starting in-memory user provider...")
}

func (i *InMemoryUserProvider) Stop() {

}
