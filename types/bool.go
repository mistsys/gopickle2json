// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

type Bool string

func NewBool(b bool) Bool {
	if b {
		return Bool("true")
	}
	return Bool("false")
}
