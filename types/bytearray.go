// Copyright 2020 NLP Odyssey Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "encoding/base64"

// ByteArray represents a Python "bytearray" (builtin type).
type ByteArray string

// NewByteArray makes and returns a new ByteArray initialized with
// the elements of the given slice.
func NewByteArray(bytes []byte) ByteArray {
	return ByteArray(base64.StdEncoding.EncodeToString(bytes))
}
