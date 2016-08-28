// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// export access to bytesconv internals for tests

package bytesconv

func NewDecimal(i uint64) *decimal {
	d := new(decimal)
	d.Assign(i)
	return d
}

func SetOptimize(b bool) bool {
	old := optimize
	optimize = b
	return old
}
