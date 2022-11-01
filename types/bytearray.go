// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/base64"
	"strings"
)

// ByteArray represents a Python "bytearray" (builtin type).
type ByteArray []byte

// NewByteArray makes and returns a new ByteArray initialized with
// the elements of the given slice.
func NewByteArray(bytes []byte) ByteArray {
	return ByteArray(bytes)
}

func (a ByteArray) JSON(b *strings.Builder) {
	w := base64.NewEncoder(base64.StdEncoding, b)
	w.Write([]byte(a))
	w.Close()
}
