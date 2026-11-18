// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"iter"
	"strings"
)

// DictSetter is implemented by any value that exhibits a dict-like behaviour,
// allowing arbitrary key/value pairs to be set.
type DictSetter interface {
	Set(key, value Object) error
	SetMany([]Object) error // pairs of key, value
	Object
}

// Dict represents a Python "dict" (builtin type) stored as successive
// pairs of key, value.
type Dict []Object

var _ DictSetter = &Dict{}

// NewDict makes and returns a new empty Dict with room for n items
func NewDict(n int, ram *[]Object) *Dict {
	d := (*Dict)(ram)
	if n > 0 {
		*d = make(Dict, 0, 2*n)
	}
	return d
}

// ValidKey returns true if key can be a dict key in JSON.
// Only some types can be keys. We match python's json.dump(), which permits only str, int, float, bool or None
func (*Dict) IsValidKey(key Object) bool {
	switch key.(type) {
	case String, Int, *Long, Float, Bool, None:
		return true
	default:
		return false
	}
}

// Set sets into the Dict the given key/value pair.
func (d *Dict) Set(key, value Object) error {
	if !d.IsValidKey(key) {
		return fmt.Errorf("type %T can't be a dict key", key)
	}
	*d = append(*d, key, value)
	return nil
}

func (d *Dict) SetMany(kv []Object) error {
	// validate the keys
	for i := 0; i < len(kv); i += 2 {
		if !d.IsValidKey(kv[i]) {
			return fmt.Errorf("type %T can't be a dict key", kv[i])
		}
	}

	// it's very common to have a single SetMany call which loads all items into an empty dict
	// so treat that as a special case not requiring any copies
	if len(*d) == 0 {
		*d = kv
		return nil
	}
	*d = append(*d, kv...)
	return nil
}

// Range generates an iterator.
//
//	for key, val := range d.Range() { ...
func (d *Dict) Range() iter.Seq2[Object, Object] {
	return func(yield func(Object, Object) bool) {
		for i := 0; i+1 < len(*d); i += 2 {
			if !yield((*d)[i], (*d)[i+1]) {
				return
			}
		}
	}
}

func (d *Dict) JSON(b *strings.Builder) {
	b.WriteByte('{')
	for i, x := range *d {
		quoted := false
		if i&1 == 0 {
			// key
			if i != 0 {
				b.WriteByte(',')
			}
			switch xx := x.(type) {
			case Int, *Long, Bool, None:
				// JSON requires quotes around non-string keys
				b.WriteByte('"')
				quoted = true
			case Float:
				if !xx.NeedsQuotes() {
					b.WriteByte('"')
					quoted = true
				} // else it's going to be quoted already
			}
		} else {
			// value
			b.WriteByte(':')
		}
		x.JSON(b)
		if quoted {
			b.WriteByte('"')
		}
	}
	b.WriteByte('}')
}
