package pickle

import (
	"encoding/hex"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/mistsys/gopickle2json/types"
)

func TestOverlongEncoding(t *testing.T) {
	// we rely on utf8.Valid() to detect overload UTF-8 encodings and reject them. check that it in fact does
	// the source code looks like it does, but the doc doesn't spell it out
	if utf8.Valid([]byte{0xc0, 0xbF}) {
		t.Error("utf8.Valid() doesn't reject overlong UTF-8 encoding")
	}
}

func TestDecodeRawUnicodeEscape(t *testing.T) {
	tests := [][2]string{ // UTF-8, raw-unicode-escaped
		{"", ``},
		{"\\", `\\`},
		{"abc 123", "abc 123"},
		{"abc\n123", "abc\n123"},
		{"\n", `\u000a`},
		{"\n", `\U0000000a`},
		{"\n", `\U0000000A`},
		{"\\u000a", `\\u000a`},
		{"\\\n", `\\\u000a`},
		{"\\\n", `\\\u000A`},
		{"\\\\u000a", `\\\\u000a`},
		{"abc\n123\\", `abc\ud800\udc0a123\\`},
		{"\000", `\ud800\udc00`},
	}
	u := NewUnpickler(nil)
	for i, test := range tests {
		out, err := u.decode_raw_unicode_escaped([]byte(test[1]))
		if err != nil {
			t.Errorf("test case %d FAILED: decode_raw_unicode_escaped(%q) failed: %v", i, test[1], err)
		} else if string(out) != test[0] {
			t.Errorf("test case %d FAILED: decode_raw_unicode_escaped(%q) returned %q; expected %q", i, test[1], out, test[0])
		} else {
			t.Logf("test case %d PASSED: decode_raw_unicode_escaped(%q) returned %q", i, test[1], out)
		}
	}
}

func TestInvalidDecodeRawUnicodeEscape(t *testing.T) {
	tests := []string{ // invalid raw-unicode-escaped
		`\`,
		`abc\n123`,
		`\u00a`,
		`\U000000aU`,
		`\ud800\ud900`,
		`\ud800ud900`,
		`\ud800\U0000dc00`,
		`\ud800\ue000`,
		`\ud800\udc`,
		`\ud800\udc0`,
		`\ud800\xdc00`,
		`\udc00\ud800`,
		`\ud800`,
		`\udbff`,
		`\udc00`,
		`\udfff`,
		`\ud836\udc00`,
		`\ud836\udfff`,
		`\ud837\udc00`,
		`\ud837\udfff`,
	}
	u := NewUnpickler(nil)
	for i, test := range tests {
		out, err := u.decode_raw_unicode_escaped([]byte(test))
		if err == nil {
			t.Errorf("test case %d FAILED: decode_raw_unicode_escaped(%q) succeeded and returned %q", i, test, out)
		} else {
			t.Logf("test case %d PASSED: decode_raw_unicode_escaped(%q) failed and returned %v", i, test, err)
		}
	}

	u.PreserveInvalidSurrogates(true)
	for i, test := range tests {
		out, err := u.decode_raw_unicode_escaped([]byte(test))
		if err == nil {
			t.Errorf("test case %d FAILED: decode_raw_unicode_escaped(%q) succeeded and returned %q", i, test, out)
		} else {
			t.Logf("test case %d PASSED: decode_raw_unicode_escaped(%q) failed and returned %v", i, test, err)
		}
	}

}

func TestPickledRawUnicodeEscape(t *testing.T) {
	tests := [][3]string{ // UTF-8, pickle, JSON
		// 8, 16 and 32-bit unicode
		{"\"café'\t🐍…\\\n", "5622636166e927095c5530303031663430645c75323032365c75303035635c75303030610a70300a2e", `"\"café'\t🐍…\\\n"`}, // pickle.dumps('"café\'\t🐍…\\\n', 0).hex()
		// backslashes
		{"\\", "565c75303035630a70300a2e", `"\\"`},                 // pickle.dumps('\\', 0).hex()
		{"\n", "565c75303030610a70300a2e", `"\n"`},                 // pickle.dumps('\n', 0).hex()
		{"\\n", "565c75303035636e0a70300a2e", `"\\n"`},             // pickle.dumps('\\n', 0).hex()
		{"\\\n", "565c75303035635c75303030610a70300a2e", `"\\\n"`}, // pickle.dumps('\\\n', 0).hex()
		// surrogate pairs (b/c python supports them)
		{"S", "565c75643830305c75646335330a70300a2e", `"S"`}, //pickle.dumps('\ud800\udc53', 0).hex() // a long way of producing "S" (53h)
	}

	for i, test := range tests {
		in, _ := hex.DecodeString(test[1])
		u := NewUnpickler(in)
		obj, err := u.Load()
		if err != nil {
			t.Errorf("test case %d FAIL: %v", i, err)
			continue
		}

		var raw string
		var str string
		switch s := obj.(type) {
		case *types.SimpleString:
			raw = string(*s)
			str = s.String()
		case *types.EscapedString:
			raw = string(*s)
			str = s.String()
		default:
			t.Errorf("test case %d FAIL: unexpected output type (%T) %v", i, obj, obj)
		}

		if str != test[0] {
			t.Errorf("test case %d FAIL: unexpected output value (%T) % x (%q) (%s); expected %q", i, obj, raw, str, raw, test[0])
		} else {
			var b strings.Builder
			obj.JSON(&b)
			if b.String() != test[2] {
				t.Errorf("test case %d FAIL: unexpected JSON value (%T) % x (%q) -> JSON %s (expected %s)", i, obj, obj, obj, b.String(), test[2])
			} else {
				t.Logf("test case %d PASS: (%T) % x (%q) -> JSON %s", i, obj, obj, obj, b.String())
			}
		}
	}
}

func TestPickledRawUnicodeEscapePreserveSP(t *testing.T) {
	// retest the surrogate pair case with conversion disabled
	in, _ := hex.DecodeString("565c75643830305c75646335330a70300a2e") // pickle.dumps('\ud800\udc53', 0).hex()
	u := NewUnpickler(in)
	u.ConvertSurrogatePairs(false)
	obj, err := u.Load()
	if err != nil {
		t.Errorf("surrogate pair test case FAIL: %v", err)
	} else {
		var b strings.Builder
		obj.JSON(&b)
		expect := `"\ud800\udc53"`
		if b.String() != expect {
			t.Errorf("surrogate pair test case FAIL: unexpected JSON value (%T) % x (%q) -> JSON %s (expected %s)", obj, obj, obj, b.String(), expect)
		} else {
			t.Logf("surrogate pair test case PASS: (%T) % x (%q) -> JSON %s", obj, obj, obj, b.String())
		}
	}
}

func TestInvalidPickledRawUnicodeEscape(t *testing.T) {
	// test cases that should fail
	tests := []struct {
		r []rune
		p string
	}{ // invalid UTF-8, pickle
		// 8, 16 and 32-bit unicode
		// backslashes
		// invalid surrogate pairs
		{[]rune{0xdc01, 0xd801}, "565c75646330315c75643830310a70300a2e"},        // pickle.dumps('\udc01\ud801', 0).hex() // wrong order
		{[]rune{0xd801}, "565c75643830310a70300a2e"},                            // pickle.dumps('\ud801', 0).hex() // an unpaired high surrogate
		{[]rune{0xd801, 0xd802}, "565c75643830315c75643830320a70300a2e"},        // pickle.dumps('\ud801\ud802', 0).hex() // two high surrogates
		{[]rune{0xdc01}, "565c75646330310a70300a2e"},                            // pickle.dumps('\udc01', 0).hex() // an unpaired low surrogate
		{[]rune{0xd801, '_', 0xdc02}, "565c75643830315f5c75646330320a70300a2e"}, // pickle.dumps('\ud801_\udc02', 0).hex() // surrogates split by a _
		// overlong utf-8 encodings
	}

	for i, test := range tests {
		in, _ := hex.DecodeString(test.p)
		u := NewUnpickler(in)
		obj, err := u.Load()
		if err == nil {
			t.Errorf("test case %d (%x) should have FAILED but returned %v instead", i, test.r, obj)
			continue
		}
		t.Logf("test case %d (%x) PASSED by returning error %v", i, test.r, err)
	}
}

func TestByteArray(t *testing.T) {
	in, _ := hex.DecodeString("8004950e00000000000000430a74657374206279746573942e") // pickle.dumps(b'test bytes').hex()
	u := NewUnpickler(in)
	obj, err := u.Load()
	if err != nil {
		t.Error(err)
		return
	}

	var b strings.Builder
	obj.JSON(&b)
	t.Logf("(%T) %v -> JSON %s", obj, obj, b.String())

	if b.String() != `"dGVzdCBieXRlcw=="` {
		t.Errorf("unexpected JSON value (%T) % x (%q) -> JSON %s", obj, obj, obj, b.String())
	}
}

func TestIntDictKeys(t *testing.T) {
	in, _ := hex.DecodeString("80049554000000000000007d94284b018c036f6e65944740000000000000008c0874776f2e7a65726f944740091eb851eb851f4740091eb851eb851f477ff00000000000008c03496e6694477ff8000000000000477ff8000000000000752e") // pickle.dumps({1:"one",2.0:"two.zero",3.14:3.14,math.inf:"Inf",math.nan:math.nan}).hex()

	u := NewUnpickler(in)
	obj, err := u.Load()
	if err != nil {
		t.Error(err)
		return
	}

	var b strings.Builder
	obj.JSON(&b)
	t.Logf("(%T) %v -> JSON %s", obj, obj, b.String())

	if b.String() != `{"1":"one","2":"two.zero","3.14":3.14,"+Inf":"Inf","NaN":"NaN"}` {
		t.Errorf("unexpected JSON value (%T) % x (%v) -> JSON %s", obj, obj, obj, b.String())
	}
}
