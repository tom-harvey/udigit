// Package udigit provides integer conversions for all UNICODE digits.
// The standard Go package routines, notably strconv.Atoi(), only work
// on digits in the ASCII range.
package udigit

import (
	"strconv"
	"unicode"
)

const (
	mathLo = 0x1D7CE // Mathematical digit range is mathLo to mathHi
	mathHi = 0x1D7FF // and holds 5 consecutive 0-9 ranges

	loDigit = 0x0660 // lowest UNICODE digit above '9'
)

// FoldDigits returns a string with higher UNICODE digits replaced by '0'-'9'
func FoldDigits(s string) string {
	return MapDigits(s, '0')
}

// MapDigits returns a string with UNICODE digits replaced by digits in the
// ten digit range that contains base.  If no replacements are needed, or if
// base is not a UNICODE digit, the argument is returned unchanged.
func MapDigits(s string, base rune) string {
	// faster path for base '0'
	if base == '0' || IsDigit(base) {
		var rs []rune
		var rLo, rHi rune // range needing no conversion
		if base == '0' {
			rLo = 0
			rHi = loDigit - 1
		} else {
			base -= rune(Udtoi(base)) // lower to the zero of the range
			rLo = base
			rHi = base + 9
		}
		ri := 0
		for _, r := range s {
			if (r < rLo || r > rHi) && IsDigit(r) {
				if rs == nil {
					rs = []rune(s) // deferred conversion
				}
				rs[ri] = rune(Udtoi(r)) + base
			}
			ri++
		}
		if rs != nil {
			s = string(rs)
		}
	}
	return s
}

// Atoi works as strconv.Atoi() for sequences of unicode.IsDigit()
func Atoi(s string) (int, error) {
	return strconv.Atoi(FoldDigits(s))
}

// ParseInt works as strconv.ParseInt() for all UNICODE digits.
func ParseInt(s string, base int, bitSize int) (i int64, err error) {
	return strconv.ParseInt(FoldDigits(s), base, bitSize)
}

// ParseUint works as strconv.ParseUint() for all UNICODE digits.
func ParseUint(s string, base int, bitSize int) (i uint64, err error) {
	return strconv.ParseUint(FoldDigits(s), base, bitSize)
}

// Udtoi returns 0 through 9 for all runes where unicode.IsDigit() is true.
// It returns -1 for all runes where unicode.IsDigit() is false.
func Udtoi(r rune) int {
	if r <= unicode.MaxRune {
		switch r >> 4 {
		// the high 16 bits of UNICODE digit ranges that end in 0x0
		case 0x10D3, 0x11DA:
			//fallthrough // these two cases are digits in UNICODE 11
		case 0x0003, 0x0066, 0x006F, 0x007C, 0x00E5, 0x00ED, 0x00F2, 0x0104,
			0x0109, 0x017E, 0x0181, 0x019D, 0x01A8, 0x01A9, 0x01B5, 0x01BB,
			0x01C4, 0x01C5, 0x0A62, 0x0A8D, 0x0A90, 0x0A9D, 0x0A9F, 0x0AA5,
			0x0ABF, 0x0FF1, 0x104A, 0x110F, 0x111D, 0x112F, 0x1145,
			0x114D, 0x1165, 0x116C, 0x1173, 0x118E, 0x11C5, 0x11D5,
			0x16A6, 0x16B5, 0x1E95:
			if r&0xf <= 9 {
				return int(r & 0xf)
			}
		}
		r6 := r - 6
		switch r6 >> 4 {
		// the high 16 bits of UNICODE digit ranges that end in 0x6
		case 0x0096, 0x009E, 0x00A6, 0x00AE, 0x00B6, 0x00BE, 0x00C6, 0x00CE,
			0x00D6, 0x00DE, 0x0194, 0x1106, 0x1113, 0x1D7F:
			if r6&0xf <= 9 {
				return int(r6 & 0xf)
			}
		}
		if mathLo <= r && r <= mathHi {
			return int((r - mathLo) % 10)
		}
	}
	return -1
}

// IsDigit acts as unicode.IsDigit() but is faster for many inputs.
func IsDigit(r rune) bool {
	return Udtoi(r) != -1
}
