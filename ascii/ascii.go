package ascii

// CutString splits string s into two, and returns the prefix before the first occurrence of character c and the suffix after c
// or s and an empty string if c is not present.
func CutString(s string, c byte) (string, string) {
	i := 0
	for i < len(s) && s[i] != c {
		i++
	}
	if i < len(s) {
		return s[:i], s[i+1:]
	}
	return s, ""
}
