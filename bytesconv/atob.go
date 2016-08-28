// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bytesconv

import "bytes"

// ParseBool returns the boolean value represented by the string.
// It accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
// Any other value returns an error.
func ParseBool(s []byte) (bool, error) {
	for _, trueBytes := range validTrueBytes {
		if bytes.Equal(s, trueBytes) {
			return true, nil
		}
	}

	for _, falseBytes := range validFalseBytes {
		if bytes.Equal(s, falseBytes) {
			return false, nil
		}
	}

	return false, syntaxError("ParseBool", s)
}

var (
	validTrueBytes = [...][]byte{
		[]byte("1"),
		[]byte("t"),
		[]byte("T"),
		[]byte("TRUE"),
		[]byte("true"),
		[]byte("True"),
	}

	validFalseBytes = [...][]byte{
		[]byte("0"),
		[]byte("f"),
		[]byte("F"),
		[]byte("FALSE"),
		[]byte("false"),
		[]byte("False"),
	}
)
