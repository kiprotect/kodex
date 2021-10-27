// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package pcap

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkPseudonymization(b *testing.B) {

	router, err := getRouter()

	if err != nil {
		b.Fatal(err)
	}

	tests := []TestStruct{
		TestStruct{
			Input:        largePCAP,
			ResponseCode: 200,
			URLParams: map[string]string{
				"format": "pcap",
			},
		},
	}

	test := tests[0]

	b.SetBytes(int64(len(test.Input)))

	for n := 0; n < b.N; n++ {
		reader := bytes.NewReader(test.Input)
		req, _ := http.NewRequest("POST", "/pseudonymize", reader)
		q := req.URL.Query()

		for key, value := range test.URLParams {
			q.Add(key, value)
		}

		q.Add("key", "testtest")

		req.URL.RawQuery = q.Encode()
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != 200 {
			b.Fatal("Invalid code")
		}

		_ = resp.Body.Bytes()
	}

}
