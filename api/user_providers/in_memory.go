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
	users map[string]*InMemoryUser
}

func MakeInMemoryUserProvider(settings kodex.Settings) (api.UserProvider, error) {
	kodex.Log.Info("Making in-memory user provider")
	return &InMemoryUserProvider{}, nil
}

type InMemoryAccessToken struct {
	scopes []string `json:"scopes"`
	token  []byte
}

func (i *InMemoryAccessToken) Scopes() []string {
	return i.scopes
}

func (i *InMemoryAccessToken) Token() []byte {
	return i.token
}

type InMemoryOrganizationRoles struct {
	roles        []string                  `json:"roles"`
	organization *InMemoryUserOrganization `json:"organization"`
}

func (i *InMemoryOrganizationRoles) Roles() []string {
	return i.roles
}

func (i *InMemoryOrganizationRoles) Organization() api.UserOrganization {
	return i.organization
}

type InMemoryUserOrganization struct {
	name        string `json:"name"`
	isDefault   bool   `json:"default"`
	description string `json:"description"`
	id          []byte `json:"id"`
}

func (i *InMemoryUserOrganization) Name() string {
	return i.name
}

func (i *InMemoryUserOrganization) Source() string {
	return "inMemory"
}

func (i *InMemoryUserOrganization) Default() bool {
	return i.isDefault
}

func (i *InMemoryUserOrganization) Description() string {
	return i.description
}

func (i *InMemoryUserOrganization) ID() []byte {
	return i.id
}

func (i *InMemoryUserOrganization) ApiOrganization(controller api.Controller) (api.Organization, error) {
	return controller.Organization("inMemory", i.id)
}

type InMemoryUser struct {
	sourceID    []byte                       `json:"sourceID"`
	email       string                       `json:"email"`
	displayName string                       `json:"displayName"`
	superuser   bool                         `json:"superuser"`
	accessToken *InMemoryAccessToken         `json:"accessToken"`
	roles       []*InMemoryOrganizationRoles `json:"roles"`
	limits      map[string]interface{}       `json:"limits"`
}

func (i *InMemoryUser) Source() string {
	return "inMemory"
}

func (i *InMemoryUser) SourceID() []byte {
	return i.sourceID
}

func (i *InMemoryUser) EMail() string {
	return i.email
}

func (i *InMemoryUser) SuperUser() bool {
	return i.superuser
}

func (i *InMemoryUser) DisplayName() string {
	return i.displayName
}

func (i *InMemoryUser) AccessToken() api.AccessToken {
	return i.accessToken
}

func (i *InMemoryUser) Roles() []api.OrganizationRoles {
	roles := make([]api.OrganizationRoles, len(i.roles))
	for i, role := range i.roles {
		roles[i] = role
	}
	return roles
}

func (i *InMemoryUser) Limits() map[string]interface{} {
	return i.limits
}

// Return a user with the given access token
func (i *InMemoryUserProvider) Get(string) (api.User, error) {
	return nil, fmt.Errorf("access token missing")
}
func (i *InMemoryUserProvider) Start() {
	kodex.Log.Info("Starting in-memory user provider...")
}

func (i *InMemoryUserProvider) Stop() {

}
