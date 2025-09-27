package ascii

// ToLowerString returns s with all upper case characters replaced with their lower case counterparts.
func ToLowerString(s string) string {
	i := 0
	for i < len(s) && !('A' <= s[i] && s[i] <= 'Z') {
		i++
	}
	if i >= len(s) {
		return s
	}
	b := append(make([]byte, 0, 32), s...)
	for ; i < len(b); i++ {
		if 'A' <= b[i] && b[i] <= 'Z' {
			b[i] ^= 0x20
		}
	}
	return string(b)
}
