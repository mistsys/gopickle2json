// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"container/list"
	"fmt"
	"strings"
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

func (*OrderedDictClass) JSON() string {
	panic("cannot serialize OrderedDictClass to JSON")
}

// OrderedDict is a minimal and trivial implementation of an ordered map,
// which represent a Python "collections.OrderedDict" object.
//
// It is composed by a simple unordered Map, and a List to keep the order of
// the entries. The former is useful for direct key lookups, the latter for
// iteration.
type OrderedDict struct {
	// Map associates a key of any type (interface{}) to OrderedDictEntry
	// pointer values. These values are shared with List.
	Map map[Object]*OrderedDictEntry
	// List is an ordered list of OrderedDictEntry pointers, which are
	// also shared with Map.
	List *list.List
	// PyDict represents Python "object.__dict__" dictionary of attributes.
	PyDict map[string]interface{}
}

var _ DictSetter = &OrderedDict{}
var _ PyDictSettable = &OrderedDict{}

// OrderedDictEntry is a single key/value pair stored in an OrderedDict.
//
// A pointer to an OrderedDictEntry is always shared between OrderedDict's Map
// and List.
type OrderedDictEntry struct {
	// Key of a single OrderedDict's entry.
	Key Object
	// Value of a single OrderedDict's entry.
	Value Object
	// ListElement is a pointer to the OrderedDict's List Element which
	// contains this very OrderedDictEntry.
	ListElement *list.Element
}

// NewOrderedDict makes and returns a new empty OrderedDict.
func NewOrderedDict() *OrderedDict {
	return &OrderedDict{
		Map:    make(map[Object]*OrderedDictEntry),
		List:   list.New(),
		PyDict: make(map[string]interface{}),
	}
}

// Set sets into the OrderedDict the given key/value pair. If the key does not
// exist yet, the new pair is positioned at the end (back) of the OrderedDict.
// If the key already exists, the existing associated value is replaced with the
// new one, and the original position is maintained.
func (o *OrderedDict) Set(k, v Object) {
	if entry, ok := o.Map[k]; ok {
		entry.Value = v
		return
	}

	entry := &OrderedDictEntry{
		Key:   k,
		Value: v,
	}
	entry.ListElement = o.List.PushBack(entry)
	o.Map[k] = entry
}

// Get returns the value associated with the given key (if any), and whether
// the key is present or not.
func (o *OrderedDict) Get(k Object) (Object, bool) {
	entry, ok := o.Map[k]
	if !ok {
		return nil, false
	}
	return entry.Value, true
}

// MustGet returns the value associated with the given key, if if it exists,
// otherwise it panics.
func (o *OrderedDict) MustGet(key Object) Object {
	value, ok := o.Get(key)
	if !ok {
		panic(fmt.Errorf("key not found in OrderedDict: %#v", key))
	}
	return value
}

// Len returns the length of the OrderedDict, that is, the amount of key/value
// pairs contained by the OrderedDict.
func (o *OrderedDict) Len() int {
	return len(o.Map)
}

// PyDictSet mimics the setting of a key/value pair on Python "__dict__"
// attribute of the OrderedDict.
func (o *OrderedDict) PyDictSet(key, value Object) error {
	sKey, keyOk := key.(String)
	if !keyOk {
		return fmt.Errorf(
			"OrderedDict.PyDictSet() requires string key: %#v", key)
	}
	o.PyDict[string(sKey)] = value
	return nil
}

func (o *OrderedDict) JSON() string {
	var b strings.Builder
	b.WriteByte('{')
	for e := o.List.Front(); e != nil; e = e.Next() {
		if b.Len() != 1 {
			b.WriteByte(',')
		}
		entry := e.Value.(*OrderedDictEntry)
		b.WriteString(entry.Key.JSON())
		b.WriteByte(':')
		b.WriteString(entry.Value.JSON())
	}
	b.WriteByte('}')
	return b.String()
}
