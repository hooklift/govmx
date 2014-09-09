// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package vmx

import "bytes"

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
