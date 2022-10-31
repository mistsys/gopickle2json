// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "fmt"

type String string

func NewString(s string) String {
	return String(s) // don't quote and escape the string yet --- we might need it to form the name of a class or something
}

func (s String) String() string {
	return fmt.Sprintf("%q", string(s)) // TODO check that upper unicode chars are ok
}
