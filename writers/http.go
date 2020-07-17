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
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/kiprotect/kiprotect"
	"net/http"
)

type HTTPWriter struct {
	Format  string
	URL     string
	Config  kiprotect.Config
	Headers map[string]interface{}
}

func (s *HTTPWriter) Teardown() error {
	return nil
}

func (s *HTTPWriter) Setup(config kiprotect.Config) error {
	s.Config = config
	return nil
}

func (s *HTTPWriter) Write(payload kiprotect.Payload) error {

	b := make([]byte, 0)
	buf := bytes.NewBuffer(b)
	writer := bufio.NewWriter(buf)

	for _, item := range payload.Items() {
		serializedItem, err := item.Serialize(s.Format)
		if err != nil {
			return err
		}
		_, err = writer.Write(serializedItem)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	writer.Flush()

	client := &http.Client{}
	req, err := http.NewRequest("POST", s.URL, buf)

	for k, v := range s.Headers {
		req.Header.Add(k, v.(string))
	}

	// if a config is given we add the config ID to the headers, so the endpoint
	// can know which config these items originate from
	if s.Config != nil {
		req.Header.Add("X-KIP-Config", hex.EncodeToString(s.Config.ID()))
	}

	req.Header.Add("Content-Type", fmt.Sprintf("application/%s", s.Format))

	_, err = client.Do(req)

	return err
}

func MakeHTTPWriter(config map[string]interface{}) (kiprotect.Writer, error) {
	if params, err := HTTPWriterForm.Validate(config); err != nil {
		return nil, err
	} else {
		return &HTTPWriter{
			Format:  params["format"].(string),
			URL:     params["url"].(string),
			Headers: params["headers"].(map[string]interface{}),
		}, nil
	}
}
