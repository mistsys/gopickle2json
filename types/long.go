// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "math/big"

type Long big.Int

func (l *Long) String() string {
	return ((*big.Int)(l)).String()
}

func (l *Long) JSON() string {
	return l.String()
}
