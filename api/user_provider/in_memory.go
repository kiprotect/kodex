package provider

import (
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/kodex/api"
)

var InMemoryUserProfileForm = forms.Form{
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

type InMemoryUserProfileProvider struct {
	userProfiles map[string]*InMemoryUserProfile
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

type InMemoryUserProfile struct {
	sourceID    []byte                       `json:"sourceID"`
	email       string                       `json:"email"`
	displayName string                       `json:"displayName"`
	superuser   bool                         `json:"superuser"`
	accessToken *InMemoryAccessToken         `json:"accessToken"`
	roles       []*InMemoryOrganizationRoles `json:"roles"`
	limits      map[string]interface{}       `json:"limits"`
}

func (i *InMemoryUserProfile) Source() string {
	return "inMemory"
}

func (i *InMemoryUserProfile) SourceID() []byte {
	return i.sourceID
}

func (i *InMemoryUserProfile) EMail() string {
	return i.email
}

func (i *InMemoryUserProfile) SuperUser() bool {
	return i.superuser
}

func (i *InMemoryUserProfile) DisplayName() string {
	return i.displayName
}

func (i *InMemoryUserProfile) AccessToken() api.AccessToken {
	return i.accessToken
}

func (i *InMemoryUserProfile) Roles() []api.OrganizationRoles {
	roles := make([]api.OrganizationRoles, len(i.roles))
	for i, role := range i.roles {
		roles[i] = role
	}
	return roles
}

func (i *InMemoryUserProfile) Limits() map[string]interface{} {
	return i.limits
}

// Return a user profile with the given access token
func (i *InMemoryUserProfileProvider) Get(string) (api.UserProfile, error) {
	return nil, nil
}
func (i *InMemoryUserProfileProvider) Start() {

}

func (i *InMemoryUserProfileProvider) Stop() {

}
