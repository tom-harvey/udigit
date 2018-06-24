// Package udigits provides conversions that work with all UNICODE digits.
// It extends the routines in strconv that work only with digits '0' to '9'.
package udigit

import (
	"strconv"
	"testing"
	"unicode"
)

const (
	// a sample of 1, 2, 3, and 4 byte UTF-8 digit runes
	udigits = "0123456789٠١٢٣٤٥٦٧٨٩௦௧௨௩௪௫௬௭௮௯𑁦𑁧𑁨𑁩𑁪𑁫𑁬𑁭𑁮𑁯０１２３４５６７８９"
)

var Sdummy string
var Idummy int
var Bdummy bool

/* TODO FoldDigits, MapDigits tests */
var testFoldDigits = []struct {
	description string
	input       string
	ok          string
}{
	{
		"empty input",
		"",
		"",
	},
	{
		"udigits sample",
		udigits,
		"01234567890123456789012345678901234567890123456789",
	},
	{
		"no substitutions needed",
		"It was a dark and stormy night...",
		"It was a dark and stormy night...",
	},
	{
		"embedded substitution",
		"ABC٠١٢٣٤٥٦٧٨٩XYZ",
		"ABC0123456789XYZ",
	},
	{
		"embedded in non-latin substitution",
		"世界٠١٢٣٤٥٦٧٨٩世界",
		"世界0123456789世界",
	},
}

var testMapDigits = []struct {
	description string
	input       string
	base        rune
	ok          string
}{
	{
		"empty input",
		"",
		'٤',
		"",
	},
	{
		"udigits sample, mapped to 𑁯",
		udigits,
		'𑁯',
		"𑁦𑁧𑁨𑁩𑁪𑁫𑁬𑁭𑁮𑁯𑁦𑁧𑁨𑁩𑁪𑁫𑁬𑁭𑁮𑁯𑁦𑁧𑁨𑁩𑁪𑁫𑁬𑁭𑁮𑁯𑁦𑁧𑁨𑁩𑁪𑁫𑁬𑁭𑁮𑁯𑁦𑁧𑁨𑁩𑁪𑁫𑁬𑁭𑁮𑁯",
	},
	{
		"no substitutions needed",
		"It was a dark and stormy night...",
		'௪',
		"It was a dark and stormy night...",
	},
	{
		"embedded substitution",
		"ABC0123456789XYZ",
		'٨',
		"ABC٠١٢٣٤٥٦٧٨٩XYZ",
	},
	{
		"embedded in non-latin substitution",
		"世界０１２３４５６７８９世界",
		'௫',
		"世界௦௧௨௩௪௫௬௭௮௯世界",
	},
}

// Atoi non-error cases
var testAtoi = []struct {
	description string
	input       string
	ok          int
}{
	{
		"empty input",
		"",
		0,
	},
	{
		"ascii digits",
		"123456789",
		123456789,
	},
	{
		"+digits",
		"+٠١٢٣٤٥٦٧٨٩",
		123456789,
	},
	{
		"-digits",
		"-௦௧௨௩௪௫௬௭௮௯",
		-123456789,
	},
}

func TestFoldDigits(t *testing.T) {
	for _, test := range testFoldDigits {
		if ok := FoldDigits(test.input); ok != test.ok {
			t.Fatalf("FoldDigits(\"%v\") %v:\n\twant: %v\n\t got: %v\n",
				test.input, test.description, test.ok, ok)
		}
	}
}
func TestMapDigits(t *testing.T) {
	for _, test := range testMapDigits {
		if ok := MapDigits(test.input, test.base); ok != test.ok {
			t.Fatalf("MapDigits(\"%v\") %v:\n\twant: %v\n\t got: %v\n",
				test.input, test.description, test.ok, ok)
		}
	}
}
func TestAtoi(t *testing.T) {
	for _, test := range testAtoi {
		ok, err := Atoi(test.input)
		if err != nil {
		}
		if ok != test.ok {
			t.Fatalf("Atoi(%v) %v: want:%v got:%v",
				test.input, test.description, test.ok, ok)
		}
	}
}

/*
if FoldDigits("") != "" || MapDigits("", '𑁩') != "" {
	t.Fatalf("FoldDigits()/MapDigits() fail on empty string")
}
*/

// TestUdigit runs tests that cover all UTF-8 rune values
func TestUdigit(t *testing.T) {
	last := -1
	dFound := 0
	for r := rune(0); r <= unicode.MaxRune; r++ {
		s := string(r)
		Sdummy = FoldDigits(s)
		Sdummy = MapDigits(s, '𑁩')
		if unicode.IsDigit(r) {
			if '9' < r && r < loDigit {
				t.Fatalf("loDigit is %v, should be %v", loDigit, r)
			}
			i := Udtoi(r)
			if i == -1 {
				t.Fatalf("Udtoi(%v), want -1, got %v", r, i)
			}
			if !IsDigit(r) {
				t.Fatalf("IsDigit(%v), want true, got false", r)
			}
			if i != last+1 {
				t.Fatalf("rune %v, want:%v got:%v", r, last+1, i)
			}
			if i == 9 {
				last = -1
			} else {
				last++
			}
			if i != dFound%10 {
				t.Fatalf("Unicode digit %v: want:%v got:%v", r, dFound, i)
			}
			iS, errS := strconv.Atoi(FoldDigits(s))
			if errS != nil {
				t.Fatalf("test error strconv.Atoi(%v): %v", s, errS)
			}
			iU, errU := Atoi(s)
			if errU != nil {
				t.Fatalf("test error Atoi(%v): %v", s, errU)
			}
			if iS != iU {
				t.Fatalf("Atoi(%v): want:%v got:%v", s, iS, iU)
			}
			dFound++
		} else { // non-digit
			if last != -1 {
				t.Fatalf("rune %v, bad sequence end %v", r, last)
			}
			if Udtoi(r) != -1 {
				t.Fatalf("Udtoi(0x%08X) want:-1 got:%v", r, Udtoi(r))
			}
			if IsDigit(r) {
				t.Fatalf("IsDigit(%v), want false, got true", r)
			}
		}
	}
}

func BenchmarkFoldDigits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sdummy = FoldDigits(udigits)
	}
}
func BenchmarkMapDigits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sdummy = MapDigits(udigits, '௫')
	}
}
func BenchmarkUdtoi(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, r := range udigits {
			Idummy = Udtoi(r)
		}
	}
}
func BenchmarkIsDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for r := rune(0); r < unicode.MaxRune; r++ {
			Bdummy = IsDigit(r)
		}
	}
}
func BenchmarkUnicodeIsDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for r := rune(0); r < unicode.MaxRune; r++ {
			Bdummy = unicode.IsDigit(r)
		}
	}
}
