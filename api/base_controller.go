// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package api

type BaseController struct {
	Definitions_ *Definitions
	Self         Controller
}

func (b *BaseController) APIDefinitions() *Definitions {
	return b.Definitions_
}

func (b *BaseController) RegisterAPIPlugin(plugin APIPlugin) error {
	b.Definitions_.Routes = append(b.Definitions_.Routes, plugin.InitializeAPI)
	if err := plugin.InitializeAdaptors(b.Definitions_.ObjectAdaptors); err != nil {
		return err
	}
	return nil
}
