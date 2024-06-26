// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"strings"
)

type ObjectClass struct{}

var _ PyNewable = &ObjectClass{}
var _ PyNewable = &ObjectClass{}

func (o *ObjectClass) PyNew(args ...Object) (Object, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("ObjectClass.PyNew called with no arguments")
	}
	switch class := args[0].(type) {
	case PyNewable:
		return class.PyNew()
	default:
		return nil, fmt.Errorf(
			"ObjectClass.PyNew unprocessable args: %#v", args)
	}
}

func (*ObjectClass) JSON(*strings.Builder) {
	panic("can't serialize ObjectClass to JSON")
}
