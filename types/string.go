// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"bytes"
	"strings"
)

type String []byte

func NewString(s []byte) String {
	// do we need to escape this string?
	for _, r := range s {
		if r < 0x20 || r == '"' || r == '\\' || r >= 0x80 {
			return escapedString(s)
		}
	}

	// we only need quotes around the string
	return String(s)
}

func escapedString(s []byte) String {
	// the rule in JSON in that JSON text must be UTF-8, or if you must, use unicode \uxxxx notation.
	// only ascii control chars (<0x20), \ and " need to be escaped, and some control chars can use \[bfnrt/] instead of \u00xx encoding.
	// (why \/ might be used I don't know, but it's there on json.org's flow chart
	// We used to use "%q" printf formatting, but it's a huge hot spot (30% of CPU time) since every single field name gets encoded this way
	var buf = bytes.NewBuffer(make([]byte, 0, len(s)*2)) // *2 is a rough guess
	buf.WriteByte('x')
	buf.WriteByte('"')
	for _, r := range string(s) {
		switch r {
		case '"':
			buf.WriteByte('\\')
			buf.WriteByte('"')
		case '\\':
			buf.WriteByte('\\')
			buf.WriteByte('\\')
		case '\u2028':
			buf.WriteString("\\u2028")
		case '\u2029':
			buf.WriteString("\\u2029")
		case '\b':
			buf.WriteByte('\\')
			buf.WriteByte('b')
		case '\f':
			buf.WriteByte('\\')
			buf.WriteByte('f')
		case '\n':
			buf.WriteByte('\\')
			buf.WriteByte('n')
		case '\r':
			buf.WriteByte('\\')
			buf.WriteByte('r')
		case '\t':
			buf.WriteByte('\\')
			buf.WriteByte('t')
		default:
			if r < 0x20 {
				buf.WriteString("\\u00")
				buf.WriteByte('0' + byte(r>>4))
				buf.WriteByte(hex[r&0xf])
			} else {
				// all other runes, even those >= 0x80, can be themselves. JSON source is UTF-8
				buf.WriteRune(r)
			}
		}
	}
	buf.WriteByte('"')
	return String(buf.Bytes())
}

const hex = "0123456789abcdef"

// return the escaped string
func (s String) JSON(b *strings.Builder) {
	b.WriteByte('"')
	b.Write(([]byte)(s))
	b.WriteByte('"')
}

// return the unescaped string (useful every once in a while)
func (s String) String() string {
	return string(s)
	// NOTE: in git history is an un-escaper, should that ever be useful again
}
