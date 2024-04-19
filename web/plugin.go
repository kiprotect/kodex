// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2024  KIProtect GmbH (HRB 208395B) - Germany
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

package web

import (
	"github.com/gospel-sh/gospel"
	"github.com/kiprotect/kodex/api"
	"io/fs"
)

type WebPluginMaker interface {
	InitializeWebPlugin(controller api.Controller) (WebPlugin, error)
}

type PluginRoutes struct {
	Main       []*gospel.RouteConfig
	Authorized []*gospel.RouteConfig
}

type WebPlugin interface {
	Routes(gospel.Context) PluginRoutes
}

type StaticFilesPlugin interface {
	StaticFiles() fs.FS
}

type AppLink struct {
	Name      string
	Path      string
	Icon      string
	Superuser bool
}

type AppLinkPlugin interface {
	AppLink() AppLink
}

type UserProviderPlugin interface {
	LoginPath() string
}
