// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "fmt"

type String string

func (s String) JSON() string {
	// escape the string if needed
	for _, r := range s {
		if r >= 0x80 || r == '"' || r == '\'' || r < ' ' {
			return fmt.Sprintf("%q", string(s)) // TODO check that upper unicode chars are ok
		}
	}
	return string(s)
}
