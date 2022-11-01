// Copyright 2022 Juniper Networks/Mist Systems. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"strconv"
	"strings"
)

type String string

func NewString(s []byte) String {
	// do we need to escape this string?
	for _, r := range s {
		if r < 0x20 || r == '"' || r == '\\' || r >= 0x80 {
			return escapedString(s)
		}
	}

	// we only need quotes around the string
	return String("\"" + string(s) + "\"")
}

func escapedString(s []byte) String {
	// the rule in JSON in that JSON text must be UTF-8, or if you must, use unicode \uxxxx notation.
	// only ascii control chars (<0x20), \ and " need to be escaped, and some control chars can use \[bfnrt/] instead of \u00xx encoding.
	// (why \/ might be used I don't know, but it's there on json.org's flow chart
	// We used to use "%q" printf formatting, but it's a huge hot spot (30% of CPU time) since every single field name gets encoded this way
	var buf strings.Builder
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
	return String(buf.String())
}

const hex = "0123456789abcdef"

// return the escaped string
func (s String) String() string {
	if s[0] == '"' {
		return string(s)
	}
	// the string has escaped characters in it
	return string(s[1:])
}

// return the unescaped string (useful every once in a while)
func (s String) RawString() string {
	n := len(s)
	if n < 2 { // should never happen
		return string(s)
	}
	if s[0] == '"' {
		// no escaping happened, just strip off the quotes
		return string(s)[1 : n-1]
	}
	// we need to unescape the text
	var buf strings.Builder
	for i := 2; i < n-1; i++ {
		c := s[i]
		if c != '\\' {
			buf.WriteByte(c)
			continue
		}
		i++
		c = s[i]
		switch c {
		case 'b':
			buf.WriteByte('\b')
			continue
		case 'f':
			buf.WriteByte('\f')
			continue
		case 'n':
			buf.WriteByte('\n')
			continue
		case 'r':
			buf.WriteByte('\r')
			continue
		case 't':
			buf.WriteByte('\t')
			continue
		}
		// must be '\uxxxx'
		x, err := strconv.ParseUint(string(s)[1:5], 16, 16)
		if err != nil {
			panic("can't unescape string " + string(s) + ": " + err.Error())
		}
		buf.WriteRune(rune(x))
	}
	return buf.String()
}
