// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "strings"

// FrozenSet represents a Python "frozenset" (builtin type).
//
// It is implemented in Go as a map with empty struct values; the actual set
// of generic "interface{}" items is thus represented by all the keys.
type FrozenSet map[Object]frozenSetEmptyStruct

type frozenSetEmptyStruct struct{}

// NewFrozenSetFromSlice makes and returns a new FrozenSet initialized
// with the elements of the given slice.
func NewFrozenSetFromSlice(slice []Object) FrozenSet {
	f := make(FrozenSet, len(slice))
	for _, item := range slice {
		f[item] = frozenSetEmptyStruct{}
	}
	return f
}

// Len returns the length of the FrozenSet.
func (f FrozenSet) Len() int {
	return len(f)
}

// Has returns whether the given value is present in the FrozenSet (true)
// or not (false).
func (f FrozenSet) Has(v Object) bool {
	_, ok := f[v]
	return ok
}

func (f FrozenSet) JSON() string {
	// we represent a FrozenSet as a list in JSON
	var b strings.Builder
	b.WriteByte('[')
	for o := range f {
		if b.Len() != 1 {
			b.WriteByte(',')
		}
		b.WriteString(o.JSON())
	}
	b.WriteByte(']')
	return b.String()
}
