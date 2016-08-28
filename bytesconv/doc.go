// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bytesconv implements conversions to and from string representations
// of basic data types.
//
// The functions are equivalents to the ones found in the standard bytesconv
// package but accept a byte slice instead of strings to avoid the memory
// allocation and copy that currently occurs when converting the byte slice to
// a string.
//
// Issue: https://github.com/golang/go/issues/2632
package bytesconv
