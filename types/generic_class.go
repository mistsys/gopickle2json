// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
)

type GenericClass struct {
	Module string
	Name   string
}

var _ PyNewable = &GenericClass{}
var _ Object = &GenericClass{}

type GenericObject struct {
	Class           *GenericClass
	ConstructorArgs []Object
}

func NewGenericClass(module, name String) *GenericClass {
	return &GenericClass{Module: string(module), Name: string(name)}
}

func (g *GenericClass) PyNew(args ...Object) (Object, error) {
	return &GenericObject{
		Class:           g,
		ConstructorArgs: args,
	}, nil
}

func (g *GenericObject) String() string {
	panic(fmt.Sprintf("can't serialize GenericObject(%s.%s) to JSON", g.Class.Module, g.Class.Name))
}
