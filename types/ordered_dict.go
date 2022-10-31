// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
)

// OrderedDictClass represent Python "collections.OrderedDict" class.
//
// This class allows the indirect creation of OrderedDict objects.
type OrderedDictClass struct{}

var _ Callable = &OrderedDictClass{}

// Call returns a new empty OrderedDict. It is equivalent to Python
// constructor "collections.OrderedDict()".
//
// No arguments are supported.
func (*OrderedDictClass) Call(args ...Object) (Object, error) {
	if len(args) != 0 {
		return nil, fmt.Errorf(
			"OrderedDictClass.Call args not supported: %#v", args)
	}
	return NewOrderedDict(), nil
}

func (*OrderedDictClass) String() string {
	panic("can't serialize OrderedDictClass to JSON")
}

// OrderedDict is a minimal and trivial implementation of an ordered map,
// which represent a Python "collections.OrderedDict" object.
//
// It is composed by a simple unordered Map, and a List to keep the order of
// the entries. The former is useful for direct key lookups, the latter for
// iteration.
type OrderedDict Dict

var _ DictSetter = &OrderedDict{}
var _ PyDictSettable = &OrderedDict{}

// NewOrderedDict makes and returns a new empty OrderedDict.
func NewOrderedDict() *OrderedDict {
	return (*OrderedDict)(NewDict())
}

// Set sets into the OrderedDict the given key/value pair. If the key does not
// exist yet, the new pair is positioned at the end (back) of the OrderedDict.
func (o *OrderedDict) Set(k, v Object) {
	(*Dict)(o).Set(k, v)
}

// PyDictSet mimics the setting of a key/value pair on Python "__dict__"
// attribute of the OrderedDict.
func (o *OrderedDict) PyDictSet(key, value Object) error {
	(*Dict)(o).Set(key, value)
	return nil
}

func (o *OrderedDict) String() string {
	return (*Dict)(o).String()
}
