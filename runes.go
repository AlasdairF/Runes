package runes

import "unicode"

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
		if equalPortable(s[i:i+n], sep) {
			return i
		}
		i++
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

// Map returns a copy of the rune slice s with all its characters modified according to the mapping function.
// If mapping returns a negative value, the character is dropped from the string with no replacement.
func Map(mapping func(r rune) rune, s []rune) []rune {
	newslice := s
	var on int
	var r rune
	for i, c := range s {
		r = mapping(r)
		if r >= 0 {
			newslice[on] = r
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

func equalPortable(a, b []rune) bool {
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
