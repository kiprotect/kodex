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

package helpers

import (
	"github.com/kiprotect/go-helpers/settings"
	"github.com/kiprotect/kodex"
	"os"
	"path/filepath"
	"strings"
)

var EnvConfigName = "KODEX_CONFIG"
var EnvSettingsName = "KODEX_SETTINGS"

func ConfigPaths() []string {
	envValue := os.Getenv(EnvConfigName)
	if envValue == "" {
		wd, err := os.Getwd()
		if err != nil {
			return []string{}
		}
		return []string{filepath.Join(wd, "config")}
	}
	return strings.Split(envValue, ":")
}

func SettingsPaths() []string {
	envValue := os.Getenv(EnvSettingsName)
	if envValue == "" {
		return []string{}
	}
	values := strings.Split(envValue, ":")
	sanitizedValues := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		sanitizedValues = append(sanitizedValues, value)
	}
	return sanitizedValues
}

func Settings(settingsPaths []string) (kodex.Settings, error) {
	return settings.MakeSettings(settingsPaths)
}
