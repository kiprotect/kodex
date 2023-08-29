// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2022  KIProtect GmbH (HRB 208395B) - Germany
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
	"fmt"
	"github.com/kiprotect/go-helpers/settings"
	"github.com/kiprotect/kodex"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var EnvSettingsName = "KODEX_SETTINGS"

func SettingsPaths() ([]string, fs.FS, error) {
	envValue := os.Getenv(EnvSettingsName)
	if envValue == "" {
		return []string{"settings"}, kodex.DefaultSettings, nil
	}
	values := strings.Split(envValue, ":")
	sanitizedValues := make([]string, 0, len(values))

	mainRoot := ""

	for _, value := range values {
		if value == "" {
			continue
		}
		var err error
		if value, err = filepath.Abs(value); err != nil {
			return nil, nil, err
		}

		root := filepath.VolumeName(value)

		if mainRoot != "" && root != mainRoot {
			return nil, nil, fmt.Errorf("cannot load settings from multiple volumes on Windows, sorry...")
		}

		// we set the main root from the path root
		mainRoot = root

		// we replace the slashes
		value = filepath.ToSlash(value)

		// we remove the volume or root part (e.g. 'c:/' on Windows or '/' on Linux)
		value = value[len(root)+1:]

		// we append the value to the sanitized paths
		sanitizedValues = append(sanitizedValues, value)
	}

	return sanitizedValues, os.DirFS(mainRoot + "/"), nil
}

func Settings(settingsPaths []string, fS fs.FS) (kodex.Settings, error) {
	return settings.MakeSettings(settingsPaths, fS)
}
