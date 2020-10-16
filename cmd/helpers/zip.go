// Kodex (Community Edition - CE) - Privacy & Security Engineering Platform
// Copyright (C) 2020  KIProtect GmbH (HRB 208395B) - Germany
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
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/kiprotect/go-helpers/yaml"
	"github.com/kiprotect/kodex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Download(url string) ([]byte, error) {

	if resp, err := http.Get(url); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}
}
func ExtractBlueprints(data []byte, dest string) error {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	if err != nil {
		return err
	}

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File, name string) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		fullPath := filepath.Join(dest, name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fullPath, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(fullPath), f.Mode())
			f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	var blueprintsPath string
	var blueprintsSettings map[string]interface{}

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/.blueprints.yml") {
			if rc, err := f.Open(); err != nil {
				return err
			} else {
				defer rc.Close()
				if data, err := ioutil.ReadAll(rc); err != nil {
					return err
				} else if err := yaml.Unmarshal(data, &blueprintsSettings); err != nil {
					return err
				}
			}
			blueprintsPath = path.Dir(f.Name)
			break
		}
	}

	if blueprintsPath == "" {
		return fmt.Errorf("no '.blueprints.yml' file found in archive")
	}

	var packageName, versionName string
	var ok bool

	if packageName, ok = blueprintsSettings["package"].(string); !ok {
		return fmt.Errorf("package name is missing")
	}

	if versionName, ok = blueprintsSettings["version"].(string); !ok {
		versionName = "default"
	}

	kodex.Log.Infof("Extracting blueprints package '%s-%s' to path '%s'", packageName, versionName, dest)

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, blueprintsPath) {
			name := filepath.Join(fmt.Sprintf("%s-%s", packageName, versionName), f.Name[len(blueprintsPath):len(f.Name)])
			if err := extractAndWriteFile(f, name); err != nil {
				return err
			}
		}
	}

	return nil
}
