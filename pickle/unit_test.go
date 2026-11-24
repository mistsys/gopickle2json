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
		{"\u00ff", "\xff"},
		{"abc 123", "abc 123"},
		{"abc\n123", "abc\n123"},
		{"\n", `\u000a`},
		{"\n", `\U0000000a`},
		{"\n", `\U0000000A`},
		{"\\u000a", `\\u000a`},
		{"\\\n", `\\\u000a`},
		{"\\\n", `\\\u000A`},
		{"\\\\u000a", `\\\\u000a`},
		{"abc\U0001000a123\\", `abc\ud800\udc0a123\\`},
		{"\U00010000", `\ud800\udc00`},
		{"\U0010ffff", `\udbff\udfff`},
	}
	u := NewUnpickler(nil)
	for i, test := range tests {
		out, err := u.decode_raw_unicode_escaped([]byte(test[1]))
		if err != nil {
			t.Errorf("test case %d FAILED: decode_raw_unicode_escaped(%q) failed: %v", i, test[1], err)
		} else if string(out) != test[0] {
			t.Errorf("test case %d FAILED: decode_raw_unicode_escaped(%q) returned %q; expected %q", i, test[1], out, test[0])
		} else {
			t.Logf("test case %d passed: decode_raw_unicode_escaped(%q) returned %q", i, test[1], out)
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
	}
	u := NewUnpickler(nil)
	for i, test := range tests {
		out, err := u.decode_raw_unicode_escaped([]byte(test))
		if err == nil {
			t.Errorf("test case %d FAILED: decode_raw_unicode_escaped(%q) succeeded and returned %q", i, test, out)
		} else {
			t.Logf("test case %d passed: decode_raw_unicode_escaped(%q) failed and returned %v", i, test, err)
		}
	}
}

func TestBitEqualDecodeRawUnicodeEscape(t *testing.T) {
	tests := []string{ // raw-unicode-escaped inputs that should pass if BitEqualStrings is enabled
		`\ud800\ud900`,
		`\ud800ud900`,
		`\ud800\U0000dc00`,
		`\ud800\ue000`,
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
	u.BitEqualStrings(true)
	for i, test := range tests {
		out, err := u.decode_raw_unicode_escaped([]byte(test))
		if err != nil {
			t.Errorf("test case %d FAILED: decode_raw_unicode_escaped(%q) failed and returned %v", i, test, err)
		} else {
			t.Logf("test case %d passed: decode_raw_unicode_escaped(%q) succeeded and returned %q", i, test, out)
		}
	}

}

func TestPickledRawUnicodeEscape(t *testing.T) {
	tests := [][3]string{ // UTF-8, pickle, JSON
		// 8, 16 and 32-bit unicode
		{"\"café'\t🐍…\\\n", "5622636166e927095c5530303031663430645c75323032365c75303035635c75303030610a70300a2e", `"\"caf\u00e9\'\t🐍…\\\n"`}, // pickle.dumps('"café\'\t🐍…\\\n', 0).hex()
		// backslashes
		{"\\", "565c75303035630a70300a2e", `"\\"`},                 // pickle.dumps('\\', 0).hex()
		{"\n", "565c75303030610a70300a2e", `"\n"`},                 // pickle.dumps('\n', 0).hex()
		{"\\n", "565c75303035636e0a70300a2e", `"\\n"`},             // pickle.dumps('\\n', 0).hex()
		{"\\\n", "565c75303035635c75303030610a70300a2e", `"\\\n"`}, // pickle.dumps('\\\n', 0).hex()
		// surrogate pairs (b/c python supports them)
		{"\U00010053", "565c75643830305c75646335330a70300a2e", "\"\U00010053\""}, // pickle.dumps('\ud800\udc53', 0).hex() // a long way of producing "S"+0x10000 (10053h)
		// fffd (the substitution char)
		{"\ufffd", "565c75666666640a70300a2e", `"\ufffd"`}, // pickle.dumps('\ufffd', protocol=0).hex() // a real fffd rune
		{"\ufffd\ufffd", "565c75666666645c75666666640a70300a2e", `"\ufffd\ufffd"`},
		{"\ufffd\ufffd\ufffd", "565c75666666645c75666666645c75666666640a70300a2e", `"\ufffd\ufffd\ufffd"`},
		{"\ufffd\ufffd\ufffd\ufffd", "565c75666666645c75666666645c75666666645c75666666640a70300a2e", `"\ufffd\ufffd\ufffd\ufffd"`},
		{"\ufffd\ufffd\ufffd\ufffd", "80049510000000000000008c0cefbfbdefbfbdefbfbdefbfbd942e", `"\ufffd\ufffd\ufffd\ufffd"`}, // protocol=5 (BINUNICODE) encoding
		{"\ufffd\ufffd\ufffd", "8004950d000000000000008c09efbfbdefbfbdefbfbd942e", `"\ufffd\ufffd\ufffd"`},
		{"\ufffd\ufffd", "8004950a000000000000008c06efbfbdefbfbd942e", `"\ufffd\ufffd"`},
		{"\ufffd", "80049507000000000000008c03efbfbd942e", `"\ufffd"`},
		{"\\udcd", "565c7530303563756463640a70300a2e", `"\\udcd"`},         // pickle.dumps('\\udcd', 0).hex() // not a \u code at all
		{"\\udcd", "80059509000000000000008c055c75646364942e", `"\\udcd"`}, // pickle.dumps('\\udcd', 5).hex()
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
				t.Logf("test case %d pass: (%T) % x (%q) -> JSON %s", i, obj, obj, obj, b.String())
			}
		}
	}
}

func TestPickledRawUnicodeEscapeWithBitEqualEnabled(t *testing.T) {
	// retest with surrogate pair conversion disabled
	tests := [][3]string{ // UTF-8, pickle, JSON
		// 8, 16 and 32-bit unicode
		{"\"café'\t🐍…\\\n", "5622636166e927095c5530303031663430645c75323032365c75303035635c75303030610a70300a2e", `"\"caf\u00e9\'\t🐍…\\\n"`}, // pickle.dumps('"café\'\t🐍…\\\n', 0).hex(); json.dumps('"café\'\t🐍…\\\n') returns '"\\"caf\\u00e9\'\\t\\ud83d\\udc0d\\u2026\\\\\\n"'

		// backslashes
		{"\\", "565c75303035630a70300a2e", `"\\"`},                 // pickle.dumps('\\', 0).hex()
		{"\n", "565c75303030610a70300a2e", `"\n"`},                 // pickle.dumps('\n', 0).hex()
		{"\\n", "565c75303035636e0a70300a2e", `"\\n"`},             // pickle.dumps('\\n', 0).hex()
		{"\\\n", "565c75303035635c75303030610a70300a2e", `"\\\n"`}, // pickle.dumps('\\\n', 0).hex()
		// surrogate pairs (b/c python supports them)
		{"\U00010053", "565c75643830305c75646335330a70300a2e", `"\ud800\udc53"`}, // pickle.dumps('\ud800\udc53', 0).hex()
		{"\ufffd", "565c75643830300a70300a2e", `"\ud800"`},                       // pickle.dumps('\ud800', protocol=0).hex() // unpaired surrogate
		{"\ufffd", "80049507000000000000008c03eda1b3942e", `"\ud873"`},           // pickle.dumps('\ud873', 5).hex() // same, in binary form
		{"\ufffd", "565c75646364650a70300a2e", `"\udcde"`},                       // pickle.dumps('\udcde', 0).hex()
		{"\ufffd", "80059507000000000000008c03edb39e942e", `"\udcde"`},           // pickle.dumps('\udcde', 5).hex()

		// weird chars
		{"\ufffd", "565c75666666640a70300a2e", `"\ufffd"`},                                           // pickle.dumps('\ufffd', protocol=0).hex() // a real fffd rune
		{"\ufffd", "80059507000000000000008c03efbfbd942e", `"\ufffd"`},                               // pickle.dumps('\ufffd', 5).hex() // a real fffd rune
		{"\ufffdabc\ufffd", "565c75666666646162635c75666666640a70300a2e", `"\ufffdabc\ufffd"`},       // pickle.dumps('\ufffdabc\ufffd', protocol=0).hex()
		{"\ufffdabc\ufffd", "8005950d000000000000008c09efbfbd616263efbfbd942e", `"\ufffdabc\ufffd"`}, // pickle.dumps('\ufffdabc\ufffd', 5).hex()

		// 0xff varies by protocol
		{"\u00ff", "56ff0a70300a2e", `"\u00ff"`},                     // pickle.dumps('\xff', 0).hex()
		{"\u00ff", "5802000000c3bf71002e", `"\u00ff"`},               // pickle.dumps('\xff', 1).hex()
		{"\u00ff", "80025802000000c3bf71002e", `"\u00ff"`},           // pickle.dumps('\xff', 2).hex()
		{"\u00ff", "80059506000000000000008c02c3bf942e", `"\u00ff"`}, // pickle.dumps('\xff', 5).hex()
	}

	for i, test := range tests {
		in, _ := hex.DecodeString(test[1])
		u := NewUnpickler(in)
		u.BitEqualStrings(true)
		obj, err := u.Load()
		if err != nil {
			t.Errorf("test case %d FAIL: %v", i, err)
		} else {
			var b strings.Builder
			obj.JSON(&b)
			expect := test[2]
			if b.String() != expect {
				t.Errorf("test case %d FAIL: unexpected JSON value (%T) % x (%q) -> JSON %s (expected %s)", i, obj, obj, obj, b.String(), expect)
			} else {
				t.Logf("test case %d pass: (%T) % x (%q) -> JSON %s", i, obj, obj, obj, b.String())
			}
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
		t.Logf("test case %d (%x) passed by returning error %v", i, test.r, err)
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

	if b.String() != `{"1":"one","2":"two.zero","3.14":3.14,"Infinity":"Inf","NaN":"NaN"}` {
		t.Errorf("unexpected JSON value (%T) % x (%v) -> JSON %s", obj, obj, obj, b.String())
	}
}

func TestBoolAndNoneDictKeys(t *testing.T) {
	in, _ := hex.DecodeString("8004950e000000000000007d9428884b01894b024e4b03752e") // pickle.dumps({True:1,False:2,None:3}).hex()

	u := NewUnpickler(in)
	obj, err := u.Load()
	if err != nil {
		t.Error(err)
		return
	}

	var b strings.Builder
	obj.JSON(&b)
	t.Logf("(%T) %v -> JSON %s", obj, obj, b.String())

	if b.String() != `{"true":1,"false":2,"null":3}` {
		t.Errorf("unexpected JSON value (%T) (%v) -> JSON %s", obj, obj, b.String())
	}
}

func TestEmptydict(t *testing.T) {
	in, _ := hex.DecodeString("80047d942e") // pickle.dumps({}).hex()

	u := NewUnpickler(in)
	obj, err := u.Load()
	if err != nil {
		t.Error(err)
		return
	}

	var b strings.Builder
	obj.JSON(&b)
	t.Logf("(%T) %v -> JSON %s", obj, obj, b.String())

	if b.String() != `{}` {
		t.Errorf("unexpected JSON value (%T) (%v) -> JSON %s", obj, obj, b.String())
	}
}

func TestTupleDictKey(t *testing.T) {
	// tuples can be dict keys in python, but not in JSON
	in, _ := hex.DecodeString("8004950b000000000000007d94884b0186944b02732e") // pickle.dumps({(True,1):2}).hex()

	u := NewUnpickler(in)
	obj, err := u.Load()
	if err != nil {
		t.Logf("tuple dict key correctly failed and returned: %v", err)
		return
	}

	var b strings.Builder
	obj.JSON(&b)
	t.Logf("(%T) %v -> JSON %s", obj, obj, b.String())
	t.Errorf("unexpected JSON value (%T) (%v) -> JSON %s", obj, obj, b.String())
}
