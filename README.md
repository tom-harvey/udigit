# Golang Integer Conversions that work with all UNICODE digit sequences.

In UNICODE, there are 60 or more sequences, each containing ten code points,
where [`unicode.IsDigit()`](https://golang.org/pkg/unicode/#IsDigit) returns
true, but the [numeric conversions in package
strconv](https://golang.org/pkg/strconv/#hdr-Numeric_Conversions) only work on
digits in the ASCII range. The udigit package provides some analogous
routines, most notably `Atoi()`, that work with all UNICODE digits, and utilty
routines for transforming digit sequences embedded in UTF-8 strings. It also
provides a version of `unicode.IsDigit()` that will be faster on some mixes of
input, at the expense of some manual maintenance required when support for
newer versions of UNICODE are incorporated into the standard Go packages.

## Source

```sh
go get github.com/tom-harvey/udigit
```

## Caveats

Only ASCII '+' and '-' characters are supported in udigits.Atoi().

Floating point conversions are not supported because even I think that is
getting a little ridiculous.

This project was undertaken mainly for my educational purposes, with only a
cursory amount of research into the problem space. If there is another
facility for doing what these routines purport to do, or some reason why these
routines are utterly unneeded, I'd be interested in hearing.

The source files are full of magic numbers which should be generated from
UNICODE reference documents, but aren't. On the other hand, I believe the
provided test routines are exhaustive enough to demonstrate when the code is
completely consistent with the behavior of `unicode.IsDigit()`, and when it
needs some manual maintenance to regain compatibility.
