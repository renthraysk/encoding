package ascii

// set64 represents a subset of ASCII characters ranging from NUL (0) to ? (0x3F)
type set64 uint64

const (
	// NewLine set representing newline '\n' and carriage return '\r'.
	NewLine set64 = 1<<'\n' | 1<<'\r'
	// HorizontalSpace set representing space ' ' and tab '\t'.
	HorizontalSpace set64 = 1<<' ' | 1<<'\t'
	// VerticalSpace set representing form feed '\f' and vertical tab '\v'.
	VerticalSpace set64 = 1<<'\f' | 1<<'\v'
	// WhiteSpace set representing all whitespace characters, newline, horizontal space, and vertical space.
	WhiteSpace set64 = NewLine | HorizontalSpace | VerticalSpace
)

// Contains returns true if byte c is a member of set, false otherwise.
func (set set64) Contains(c byte) bool {
	return c < 64 && (1<<c)&set != 0
}

// TrimString removes characters in the set from the beginning and end of s.
func (set set64) TrimString(s string) string {
	for len(s) > 0 && set.Contains(s[0]) {
		s = s[1:]
	}
	for len(s) > 0 && set.Contains(s[len(s)-1]) {
		s = s[:len(s)-1]
	}
	return s
}
