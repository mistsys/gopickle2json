// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "strings"

// FrozenSet represents a Python "frozenset" (builtin type).
type FrozenSet string

// NewFrozenSetFromSlice makes and returns a new FrozenSet initialized
// with the elements of the given slice.
func NewFrozenSetFromSlice(slice []Object) FrozenSet {
	// we represent a FrozenSet as a list in JSON
	var b strings.Builder
	b.WriteByte('[')
	for i, o := range slice {
		if i != 0 {
			b.WriteByte(',')
		}
		b.WriteString(o.String())
	}
	b.WriteByte(']')
	return FrozenSet(b.String())
}

func (f FrozenSet) String() string { return string(f) }
