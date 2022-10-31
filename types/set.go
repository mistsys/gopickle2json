// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "strings"

// SetAdder is implemented by any value that exhibits a set-like behaviour,
// allowing arbitrary values to be added.
type SetAdder interface {
	Add(v Object)
	Object
}

// Set represents a Python "set" (builtin type).
type Set strings.Builder

var _ SetAdder = &Set{}

// NewSet makes and returns a new empty Set.
func NewSet() *Set {
	// we represent a Set as a list in JSON
	var b strings.Builder
	b.WriteByte('[')
	return (*Set)(&b)
}

// Add adds one element to the Set.
func (s *Set) Add(v Object) {
	var b = (*strings.Builder)(s)
	if b.Len() != 1 {
		b.WriteByte(',')
	}
	b.WriteString(v.String())
}

func (s *Set) String() string {
	return (*strings.Builder)(s).String() + "]"
}
