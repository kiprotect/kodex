// KIProtect (Community Edition - CE) - Privacy & Security Engineering Platform
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

package writers

import (
	"encoding/json"
	"fmt"
	"github.com/kiprotect/kiprotect"
	"os"
)

type StdoutWriter struct {
}

func (s *StdoutWriter) Teardown() error {
	return nil
}

func (s *StdoutWriter) Setup(config kiprotect.Config) error {
	return nil
}

func (s *StdoutWriter) Write(payload kiprotect.Payload) error {
	for _, item := range payload.Items() {
		v, err := json.MarshalIndent(item.All(), "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			fmt.Println(string(v))
		}
	}

	return nil
}

func MakeStdoutWriter(params map[string]interface{}) (kiprotect.Writer, error) {
	return &StdoutWriter{}, nil
}
