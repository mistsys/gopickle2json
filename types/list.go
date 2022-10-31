// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"strings"
)

// ListAppender is implemented by any value that exhibits a list-like
// behaviour, allowing arbitrary values to be appended.
type ListAppender interface {
	Append(v Object)
	Object
}

// List represents a Python "list" (builtin type).
type List strings.Builder

var _ ListAppender = &List{}

// NewList makes and returns a new empty List.
func NewList() *List {
	var l strings.Builder
	l.WriteString("[")
	return (*List)(&l)
}

// NewListFromSlice makes and returns a new List initialized with the elements
// of the given slice.
func NewListFromSlice(slice []Object) *List {
	var b strings.Builder
	b.WriteByte('[')
	for i, obj := range slice {
		if i != 0 {
			b.WriteByte(',')
		}
		b.WriteString(toString(obj))
	}
	return (*List)(&b)
}

// Append appends one element to the end of the List.
func (l *List) Append(obj Object) {
	b := (*strings.Builder)(l)
	if b.Len() != 1 {
		b.WriteByte(',')
	}
	b.WriteString(toString(obj))
}

func (l *List) String() string {
	return (*strings.Builder)(l).String() + "]"
}
