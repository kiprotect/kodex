package api

type UserProviderSettings struct {
	Type     string      `json:"type"`
	Settings interface{} `json:"settings"`
}

type Settings struct {
	UserProvider *UserProviderSettings `json:"userProvider"`
}
