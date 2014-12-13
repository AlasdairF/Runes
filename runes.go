package runes

import "unicode"

// Fields splits the slice s around each instance of one or more consecutive white space characters, returning a slice of subslices of s or an empty list if s contains only white space.
func Fields(s []rune) [][]rune {
	return FieldsFunc(s, unicode.IsSpace)
}

// Fields splits the slice s around each instance of one or more consecutive white space characters, returning a slice of subslices of s or an empty list if s contains only white space.
func FieldsBytes(s []byte) [][]rune {
	return FieldsFuncBytes(s, unicode.IsSpace)
}

// FieldsFunc interprets s as a sequence of UTF-8-encoded Unicode code points.
// It splits the slice s at each run of code points c satisfying f(c) and
// returns a slice of subslices of s.  If all code points in s satisfy f(c), or
// len(s) == 0, an empty slice is returned.
// FieldsFunc makes no guarantees about the order in which it calls f(c).
// If f does not return consistent results for a given c, FieldsFunc may crash.
func FieldsFunc(s []rune, f func(rune) bool) [][]rune {
	var n int
	var inField, wasInField bool
	var r rune
	for _, r = range s {
		wasInField = inField
		inField = !f(r)
		if inField && !wasInField {
			n++
		}
	}

	a := make([][]rune, n)
	var i int
	var inf bool
	fieldStart := -1
	n = 0
	for i, r = range s {
		inf = f(r)
		if fieldStart < 0 && !inf {
			fieldStart = i
			continue
		}
		if fieldStart >= 0 && inf {
			a[n] = s[fieldStart:i]
			n++
			fieldStart = -1
		}
	}
	
	return a[0:n]
}

// FieldsFuncBytes is the same as FieldsFunc but the input is a slice of bytes and output is runes.
func FieldsFuncBytes(s []byte, f func(rune) bool) [][]rune {
	var n, size, i int
	var r rune
	var inField, wasInField bool
	l := len(s)
	for i = 0; i < l; i += size {
		r, size = utf8.DecodeRune(s[i:])
		wasInField = inField
		inField = !f(r)
		if inField && !wasInField {
			n++
		}
	}

	a := make([][]rune, n)
	buf := make([]rune, 0, 1)
	var na int
	for i = 0; i <= l && na < n; {
		r, size = utf8.DecodeRune(s[i:])
		if size == 0 {
			break
		}
		if f(r) {
			if len(buf) > 0 {
				a[na] = buf
				na++
				buf = make([]rune, 0, 1)
			}
		} else {
			buf = append(buf, r)
		}
		i += size
	}
	if len(buf) > 0 {
		a[na] = buf
		na++
	}
	return a[0:na]
}

// Contains reports whether subslice is within b.
func Contains(b, subslice []rune) bool {
	return Index(b, subslice) != -1
}

// Index returns the index of the first instance of sep in s, or -1 if sep is not present in s.
func Index(s, sep []rune) int {
	n := len(sep)
	if n == 0 {
		return 0
	}
	if n > len(s) {
		return -1
	}
	c := sep[0]
	if n == 1 {
		return IndexRune(s, c)
	}
	i := 0
	var o int
	t := s[:len(s)-n+1]
	for i < len(t) {
		if t[i] != c {
			o = IndexRune(t[i:], c)
			if o < 0 {
				break
			}
			i += o
		}
		if Equal(s[i:i+n], sep) {
			return i
		}
		i++
	}
		return -1
}

// LastIndex returns the index of the last instance of sep in s, or -1 if sep is not present in s.
func LastIndex(s, sep []rune) int {
	n := len(sep)
	if n == 0 {
		return len(s)
	}
	c := sep[0]
	for i := len(s) - n; i >= 0; i-- {
		if s[i] == c && (n == 1 || Equal(s[i:i+n], sep)) {
			return i
		}
	}
	return -1
}

// IndexRune returns the rune index of the first occurrence in s of the given rune c, or -1 if c is not present in s.
func IndexRune(s []rune, c rune) int {
	for i, b := range s {
		if b == c {
			return i
		}
	}
	return -1
}

// explode splits s into a slice of UTF-8 sequences, one per Unicode character (still slices of bytes), up to a maximum of n byte slices. Invalid UTF-8 sequences are chopped into individual bytes.
func explode(s []rune, n int) [][]rune {
	if n <= 0 {
		n = len(s)
	}
	a := make([][]rune, n)
	na := 0
	for len(s) > 0 {
		if na+1 >= n {
			a[na] = s
			na++
			break
		}
		a[na] = s[0:1]
		s = s[1:]
		na++
	}
	return a[0:na]
}

func Count(s, sep []rune) int {
	n := len(sep)
	if n == 0 {
		return len(s) + 1
	}
	if n > len(s) {
		return 0
	}
	count := 0
	c := sep[0]
	i := 0
	t := s[:len(s)-n+1]
	for i < len(t) {
		if t[i] != c {
 			o := IndexRune(t[i:], c)
			if o < 0 {
				break
			}
			i += o
		}
		if n == 1 || Equal(s[i:i+n], sep) {
			count++
			i += n
			continue
		}
		i++
	}
	return count
}

// Generic split: splits after each instance of sep, including sepSave bytes of sep in the subslices.
func genSplit(s, sep []rune, sepSave, n int) [][]rune {
	if n == 0 {
		return nil
	}
	if len(sep) == 0 {
		return explode(s, n)
	}
	if n < 0 {
		n = Count(s, sep) + 1
	}
 	c := sep[0]
	start := 0
	a := make([][]rune, n)
	na := 0
	for i := 0; i+len(sep) <= len(s) && na+1 < n; i++ {
		if s[i] == c && (len(sep) == 1 || Equal(s[i:i+len(sep)], sep)) {
			a[na] = s[start : i+sepSave]
			na++
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	a[na] = s[start:]
	return a[0 : na+1]
}

// SplitN slices s into subslices separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, SplitN splits after each rune.
// The count determines the number of subslices to return:
//   n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//   n == 0: the result is nil (zero subslices)
//   n < 0: all subslices
func SplitN(s, sep []rune, n int) [][]rune { return genSplit(s, sep, 0, n) }

// SplitAfterN slices s into subslices after each instance of sep and
// returns a slice of those subslices.
// If sep is empty, SplitAfterN splits after each rune.
// The count determines the number of subslices to return:
//   n > 0: at most n subslices; the last subslice will be the unsplit remainder.
//   n == 0: the result is nil (zero subslices)
//   n < 0: all subslices
func SplitAfterN(s, sep []rune, n int) [][]rune {
	return genSplit(s, sep, len(sep), n)
}

// Split slices s into all subslices separated by sep and returns a slice of
// the subslices between those separators.
// If sep is empty, Split splits after each rune.
// It is equivalent to SplitN with a count of -1.
func Split(s, sep []rune) [][]rune { return genSplit(s, sep, 0, -1) }

// SplitAfter slices s into all subslices after each instance of sep and
// returns a slice of those subslices.
// If sep is empty, SplitAfter splits after each rune.
// It is equivalent to SplitAfterN with a count of -1.
func SplitAfter(s, sep []rune) [][]rune {
	return genSplit(s, sep, len(sep), -1)
}

// Map returns a copy of the rune slice s with all its characters modified according to the mapping function.
// If mapping returns a negative value, the character is dropped from the string with no replacement.
func Map(mapping func(r rune) rune, s []rune) []rune {
	newslice := s
	var on int
	var c rune
	for _, c = range s {
		if c = mapping(c); c >= 0 {
			newslice[on] = c
			on++
		}
	}
	return newslice[0:on]
}

// ToUpper returns a copy of the rune slice s with all Unicode letters mapped to their upper case.
func ToUpper(s []rune) []rune { return Map(unicode.ToUpper, s) }

// ToLower returns a copy of the rune slice s with all Unicode letters mapped to their lower case.
func ToLower(s []rune) []rune { return Map(unicode.ToLower, s) }

// ToTitle returns a copy of the rune slice s with all Unicode letters mapped to their title case.
func ToTitle(s []rune) []rune { return Map(unicode.ToTitle, s) }

// ToUpperSpecial returns a copy of the rune slice s with all Unicode letters mapped to their upper case, giving priority to the special casing rules.
func ToUpperSpecial(_case unicode.SpecialCase, s []rune) []rune {
	return Map(func(r rune) rune { return _case.ToUpper(r) }, s)
}

// ToLowerSpecial returns a copy of the rune slice s with all Unicode letters mapped to their lower case, giving priority to the special casing rules.
func ToLowerSpecial(_case unicode.SpecialCase, s []rune) []rune {
	return Map(func(r rune) rune { return _case.ToLower(r) }, s)
}
   
// ToTitleSpecial returns a copy of the rune slice s with all Unicode letters mapped to their title case, giving priority to the special casing rules.
func ToTitleSpecial(_case unicode.SpecialCase, s []rune) []rune {
	return Map(func(r rune) rune { return _case.ToTitle(r) }, s)
}

func Equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range a {
		if c != b[i] {
			return false
		}
	}
	return true
}
