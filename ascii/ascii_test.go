package ascii

import "testing"

func TestChangeCase(t *testing.T) {
	tests := []*struct{ in, lower, upper string }{
		{in: "@ABCDEFGHIJKLMNOPQRSTUVWXYZ[", lower: "@abcdefghijklmnopqrstuvwxyz["},
		{in: "`abcdefghijklmnopqrstuvwxyz{", lower: "`abcdefghijklmnopqrstuvwxyz{"},
	}

	n0 := testing.AllocsPerRun(1, func() {
		for _, s := range tests {
			if got := ToLowerString(s.in); s.lower != got {
				t.Fatalf("ToLower %q failed expected %q, got %q", s.in, s.lower, string([]byte(got))) // prevent got escaping
			}
		}
	})
	if testing.CoverMode() == "" && n0 != 0 {
		t.Fatalf("expected toLower & toUpper not to allocate, allocated %v times", n0)
	}

	// Having to modify a string longer than translate() stack buffer will cause append to allocate
	n1 := testing.AllocsPerRun(1, func() {
		s := ToLowerString("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
		if s != "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz" {
			t.Fatalf("toLower failed to lower")
		}
	})
	if testing.CoverMode() == "" && 0 <= n1 && n1 <= 1 {
		t.Fatalf("expected ToLowerString() to allocate atmost once, allocated %v times", n1)
	}
}

func TestCut(t *testing.T) {
	tests := []struct {
		in, prefix, suffix string
		split              byte
	}{
		{in: "abc=def", split: '=', prefix: "abc", suffix: "def"},
		{in: "abcdef", split: '=', prefix: "abcdef", suffix: ""},
	}

	for _, s := range tests {
		prefix, suffix := CutString(s.in, s.split)
		if prefix != s.prefix || suffix != s.suffix {
			t.Fatalf("CutString(%q, %v) failed, expected prefix %q suffix %q, got prefix %q suffix %q", s.in, s.split, s.prefix, s.suffix, prefix, suffix)
		}
	}
}

func TestSet64TrimString(t *testing.T) {
	tests := []struct {
		name         string
		set          set64
		in, expected string
	}{
		{"HorizontalSpace", HorizontalSpace, "abc", "abc"},
		{"HorizontalSpace", HorizontalSpace, " \tabc", "abc"},
		{"HorizontalSpace", HorizontalSpace, "abc \t", "abc"},
		{"HorizontalSpace", HorizontalSpace, " \tabc \t", "abc"},
	}

	for _, s := range tests {
		if got := s.set.TrimString(s.in); s.expected != got {
			t.Fatalf("%s.TrimString(%q) failed: expected %q, got %q", s.name, s.in, s.expected, got)
		}
	}
}
