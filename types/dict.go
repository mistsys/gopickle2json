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
type Dict strings.Builder

var _ DictSetter = &Dict{}

// NewDict makes and returns a new empty Dict.
func NewDict() *Dict {
	var d strings.Builder
	d.WriteByte('{')
	return (*Dict)(&d)
}

// Set sets into the Dict the given key/value pair.
func (d *Dict) Set(key, value Object) {
	b := (*strings.Builder)(d)
	if b.Len() != 1 {
		b.WriteByte(',')
	}
	b.WriteString(key.String())
	b.WriteByte(':')
	b.WriteString(value.String())
}

func (d *Dict) String() string {
	return (*strings.Builder)(d).String() + "}"
}
