// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"strconv"
	"strings"
)

type String interface {
	String() string
	imAString() // since String() is quite common, add a dummy method unlikely to collide
	Object
}

type SimpleString []byte  // a string that doesnt need to be escaped in JSON
type EscapedString []byte // a string that might need to be escaped in JSON
type ExactString []byte   // a string that might need to be escaped and which needs UTF-16 surrogate pairs and other non-UTF-8 garbage encoded in JSON as-is

func NewString(s []byte, ram *[]byte, bitequal bool) Object {
	*ram = s
	// do we need to escape this string?
	for _, r := range s {
		if r < 0x20 || r == '"' || r == '\\' || r >= 0x80 {
			if bitequal {
				return (*ExactString)(ram)
			}
			return (*EscapedString)(ram)
		}
	}
	// the entire string is escape-free (a common case)
	// IDEA: since we return Object interface, and that forces the SimpleString to be on the heap, keep a buffer of read-to-go SimpleString
	// objects, or even []byte objects, we can reuse for this, and also List, Dict and friends. This would reduce calls to alloc a huge amount,
	// and not use much extra RAM since we're going to allocate and hold everything in RAM at the end of Load() anyway.
	return (*SimpleString)(ram)
}

func (s *EscapedString) JSON(b *strings.Builder) {
	escape_string(b, string(*s), false)
}
func (s *ExactString) JSON(b *strings.Builder) {
	escape_string(b, string(*s), true)
}

// if preserve_surrogate_pairs then valid surrogate pairs are converted into a single UTF-8 character
// (invalid or mismatched pairs always encode as \u####)
func escape_string(b *strings.Builder, str string, preserve_surrogate_pairs bool) {
	b.WriteByte('"')
	// the rule in JSON in that JSON text must be UTF-8, or if you must, use unicode \uxxxx notation.
	// only ascii control chars (<0x20), \ and " need to be escaped, and some control chars can use \[bfnrt/] instead of \u00xx encoding.
	// (why / might need to be escaped I don't know, but it's there on json.org's flow chart. The code in stdlib doesn't escape it, so I won't either)
	var surrogate_hi rune
	for _, r := range str {
		if surrogate_hi != 0 && (r < 0xdc00 || 0xe000 <= r) {
			// there isn't going to be the 2nd half. JSON supports this by using \u####
			b.WriteString("\\u")
			b.WriteString(strconv.FormatUint(uint64(surrogate_hi), 16))
			surrogate_hi = 0
		}

		switch r {
		case '"':
			b.WriteByte('\\')
			b.WriteByte('"')
		case '\\':
			b.WriteByte('\\')
			b.WriteByte('\\')
		case '\u2028':
			b.WriteString("\\u2028")
		case '\u2029':
			b.WriteString("\\u2029")
		case '\b':
			b.WriteByte('\\')
			b.WriteByte('b')
		case '\f':
			b.WriteByte('\\')
			b.WriteByte('f')
		case '\n':
			b.WriteByte('\\')
			b.WriteByte('n')
		case '\r':
			b.WriteByte('\\')
			b.WriteByte('r')
		case '\t':
			b.WriteByte('\\')
			b.WriteByte('t')
		default:
			if r < 0x20 {
				b.WriteString("\\u00")
				b.WriteByte('0' + byte(r>>4))
				b.WriteByte(hex[r&0xf])
			} else if preserve_surrogate_pairs && 0xd800 <= r && r < 0xe000 {
				// encode each half of the pair as \u####
				b.WriteString("\\u")
				b.WriteString(strconv.FormatUint(uint64(r), 16))
			} else if 0xd800 <= r && r < 0xdc00 {
				// 1st (high) half of surrogate pair
				surrogate_hi = r
			} else if 0xdc00 <= r && r < 0xe000 {
				// 2nd (low) half of surrogate pair. reassemble
				if surrogate_hi != 0 {
					b.WriteRune(r&0x3ff | (surrogate_hi&0x3ff)<<10)
					surrogate_hi = 0
				} else {
					// else there was no preceeding high half. JSON does support this by using \u####
					b.WriteString("\\u")
					b.WriteString(strconv.FormatUint(uint64(r), 16)) // since r is between dc00 and dfff no 0 prefixes are needed
				}
			} else {
				// all other runes, even those >= 0x80, can be themselves. JSON source is UTF-8
				b.WriteRune(r)
			}
		}
	}
	if surrogate_hi != 0 {
		// truncated before a lower half. JSON does support this by using \u####
		b.WriteString("\\u")
		b.WriteString(strconv.FormatUint(uint64(surrogate_hi), 16))
	}
	b.WriteByte('"')
}

const hex = "0123456789abcdef"

// return the string in quotes
func (s *SimpleString) JSON(b *strings.Builder) {
	b.WriteByte('"')
	b.Write(([]byte)(*s))
	b.WriteByte('"')
}

// return the unescaped string (useful every once in a while)
func (s *SimpleString) String() string {
	return string(*s)
}

func (s *EscapedString) String() string {
	return string(*s)
}

func (s *ExactString) String() string {
	return string(*s)
}

func (*SimpleString) imAString()  {}
func (*EscapedString) imAString() {}
func (*ExactString) imAString()   {}
