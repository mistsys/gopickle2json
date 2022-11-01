// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"strconv"
	"strings"
)

type Float float64

func NewFloat(f float64) Float {
	return Float(f)
}

func (f Float) JSON(b *strings.Builder) {
	dst := make([]byte, 0, 32) // TODO sync.Pool if this is a hot spot
	dst = strconv.AppendFloat(dst, float64(f), 'G', -1, 64)
	b.Write(dst)
}
