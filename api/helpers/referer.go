// Kodex (Enterprise Edition - EE) - Privacy & Security Engineering Platform
// Copyright (C) 2019-2021  KIProtect GmbH (HRB 208395B) - All Rights Reserved

package helpers

import (
	"fmt"
	"github.com/kiprotect/kodex"
	"net/url"
)

// Removes the query (which is potentially personally identifiable information) from the URL
func SanitizeReferer(referer string) string {
	if parsedUrl, err := url.Parse(referer); err != nil {
		kodex.Log.Errorf("Invalid referer: %s", referer)
		return ""
	} else {
		return fmt.Sprintf("%s://%s%s", parsedUrl.Scheme, parsedUrl.Host, parsedUrl.Path)
	}
}
