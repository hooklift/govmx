// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package vmx

import (
	"bytes"
	"fmt"
	"strings"
)

// Marshal traverses the value v recursively.
// If an encountered value implements the Marshaler interface
// and is not a nil pointer, Marshal calls its MarshalVMX method
// to produce VMX.  The nil pointer exception is not strictly necessary
// but mimics a similar, necessary exception in the behavior of
// UnmarshalVMX.
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := NewEncoder(&b).Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Takes VMX data and binds it to the Go value pointed by v
func Unmarshal(data []byte, v interface{}) error {
	return NewDecoder(bytes.NewReader(data), false).Decode(v)
}

// Parses struct tag
func parseTag(tag string) (string, bool, error) {
	omitempty := false

	// Takes out first colon found
	parts := strings.Split(tag, ":")
	if len(parts) < 2 || parts[1] == "" {
		return "", omitempty, fmt.Errorf("Invalid tag: %s", tag)
	}

	if parts[1] == `""` {
		return "", omitempty, fmt.Errorf("Tag name is missing: %s", tag)
	}

	// Takes out double quotes
	parts2 := strings.Split(parts[1], `"`)
	if len(parts2) < 2 {
		return "", omitempty, fmt.Errorf("Tag name has to be enclosed in double quotes: %s", tag)
	}

	values := strings.Split(parts2[1], ",")
	if len(values) > 1 && values[1] == "omitempty" {
		omitempty = true

	}

	return values[0], omitempty, nil
}
