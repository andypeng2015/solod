// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package errors implements functions to manipulate errors.
package errors

// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
//
//so:extern
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
//
//so:extern
type errorString struct {
	s string
}

//so:extern
func (e *errorString) Error() string {
	return e.s
}
