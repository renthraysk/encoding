package ascii

// translateString translates any byte in s between (lower, upper) to (newLower, newLower + upper - lower)
func translateString(s string, lower, upper, newLower byte) string {
	// This func is expected to inline (go1.16) with
	// enough cost headroom so immediate callers,
	// will also inline.
	i := 0
	// unsigned comparison logic to reduce cost for inlining purposes
	for i < len(s) && s[i]-lower > upper-lower {
		i++
	}
	if i >= len(s) {
		// no bytes in (lower, upper) are present
		return s
	}
	// Fixed capacity make() for possible stack allocation if result doesn't escape
	// Related issue: https://github.com/golang/go/issues/47524
	b := append(make([]byte, 0, 32), s...)
	for ; i < len(b); i++ {
		if b[i]-lower <= upper-lower {
			b[i] += newLower - lower
		}
	}
	return string(b)
}

// ToLower returns s with all upper case characters replaced with their lower case counterparts.
func ToLowerString(s string) string {
	// This func is expected to inline (go1.16)
	return translateString(s, 'A', 'Z', 'a')
}
