// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "strconv"

type Int int64

func (i Int) JSON() string {
	return strconv.FormatInt(int64(i), 10)
}
