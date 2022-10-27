// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"reflect"
	"strings"
)

// DictSetter is implemented by any value that exhibits a dict-like behaviour,
// allowing arbitrary key/value pairs to be set.
type DictSetter interface {
	Set(key, value Object)
	Object
}

// Dict represents a Python "dict" (builtin type).
//
// It is implemented as a slice, instead of a map, because in Go not
// all types can be map's keys (e.g. slices).
type Dict []DictEntry

type DictEntry struct {
	Key   Object
	Value Object
}

var _ DictSetter = &Dict{}

// NewDict makes and returns a new empty Dict.
func NewDict() *Dict {
	d := make(Dict, 0, 4)
	return &d
}

// Set sets into the Dict the given key/value pair.
func (d *Dict) Set(key, value Object) {
	*d = append(*d, DictEntry{
		Key:   key,
		Value: value,
	})
}

// Get returns the value associated with the given key (if any), and whether
// the key is present or not.
func (d *Dict) Get(key Object) (Object, bool) {
	for _, entry := range *d {
		if reflect.DeepEqual(entry.Key, key) {
			return entry.Value, true
		}
	}
	return nil, false
}

// MustGet returns the value associated with the given key, if if it exists,
// otherwise it panics.
func (d *Dict) MustGet(key Object) Object {
	value, ok := d.Get(key)
	if !ok {
		panic(fmt.Errorf("key not found in Dict: %#v", key))
	}
	return value
}

// Len returns the length of the Dict, that is, the amount of key/value pairs
// contained by the Dict.
func (d *Dict) Len() int {
	return len(*d)
}

func (d *Dict) JSON() string {
	var b strings.Builder
	b.WriteByte('{')
	for i, e := range *d {
		if i != 0 {
			b.WriteByte(',')
		}
		b.WriteString(e.Key.JSON())
		b.WriteByte(':')
		b.WriteString(e.Value.JSON())
	}
	b.WriteByte('}')
	return b.String()
}
