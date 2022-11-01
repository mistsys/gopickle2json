// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"strings"
)

// DictSetter is implemented by any value that exhibits a dict-like behaviour,
// allowing arbitrary key/value pairs to be set.
type DictSetter interface {
	Set(key, value Object)
	Object
}

// Dict represents a Python "dict" (builtin type).
type Dict []Object

var _ DictSetter = &Dict{}

// NewDict makes and returns a new empty Dict.
func NewDict() *Dict {
	var d Dict
	return &d
}

// Set sets into the Dict the given key/value pair.
func (d *Dict) Set(key, value Object) {
	*d = append(*d, key, value)
}

func (d *Dict) JSON(b *strings.Builder) {
	b.WriteByte('{')
	for i, x := range *d {
		if i&1 == 0 {
			if i != 0 {
				b.WriteByte(',')
			}
		} else {
			b.WriteByte(':')
		}
		x.JSON(b)
	}
	b.WriteByte('}')
}
