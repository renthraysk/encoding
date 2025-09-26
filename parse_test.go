package encoding

import (
	"testing"
)

func TestEncodingLookup(t *testing.T) {

	if e, ok := encodingSet(Identity.String()); !ok || e != 1<<Identity {
		t.Errorf("expected %q to return constant %s(%d)", Identity.String(), e, e)
	}
	if e, ok := encodingSet(Deflate.String()); !ok || e != 1<<Deflate {
		t.Errorf("expected %q to return constant %s(%d)", Deflate.String(), e, e)
	}
	if e, ok := encodingSet(Compress.String()); !ok || e != 1<<Compress {
		t.Errorf("expected %q to return constant %s(%d)", Compress.String(), e, e)
	}
	if e, ok := encodingSet(Gzip.String()); !ok || e != 1<<Gzip {
		t.Errorf("expected %q to return constant %s(%d)", Gzip.String(), e, e)
	}
	if e, ok := encodingSet(Brotli.String()); !ok || e != 1<<Brotli {
		t.Errorf("expected %q to return constant %s(%d)", Brotli.String(), e, e)
	}
	if e, ok := encodingSet(Zstd.String()); !ok || e != 1<<Zstd {
		t.Errorf("expected %q to return constant %s(%d)", Zstd.String(), e, e)
	}
	if e, ok := encodingSet(DCZ.String()); !ok || e != 1<<DCZ {
		t.Errorf("expected %q to return constant %s(%d)", DCZ.String(), e, e)
	}
	if e, ok := encodingSet(DCB.String()); !ok || e != 1<<DCB {
		t.Errorf("expected %q to return constant %s(%d)", DCB.String(), e, e)
	}
	// Zstd is the last supported encoding, make sure test is update if another is added.
	if _, ok := encodingSet((DCZ + 1).String()); ok {
		t.Errorf("TestEncodingLookup needs updating with new encoding")
	}
}

var parseTests = map[string]EncodingSet{
	"":                                   allSet,
	"*":                                  allSet,
	"identity":                           1 << Identity,
	"br":                                 1<<Identity | 1<<Brotli,
	"bR":                                 1<<Identity | 1<<Brotli,
	"gzip":                               1<<Identity | 1<<Gzip,
	"GzIp":                               1<<Identity | 1<<Gzip,
	"deflate":                            1<<Identity | 1<<Deflate,
	"compress":                           1<<Identity | 1<<Compress,
	"GzIp, Br":                           1<<Identity | 1<<Gzip | 1<<Brotli,
	" gzip, br ":                         1<<Identity | 1<<Brotli | 1<<Gzip,
	"gzip, deflate, br":                  1<<Identity | 1<<Gzip | 1<<Deflate | 1<<Brotli,
	"identity;q=0":                       0,
	" gzip, br;q=0 ":                     1<<Identity | 1<<Gzip,
	"gzip, br;p=x;q=0":                   1<<Identity | 1<<Gzip,
	"*;q=0, br":                          1 << Brotli,
	"identity;q=0, gzip":                 1 << Gzip,
	"compress, gzip":                     1<<Compress | 1<<Gzip | 1<<Identity,
	"compress;q=0.5, gzip;q=1.0":         1<<Compress | 1<<Gzip | 1<<Identity,
	"gzip;q=1.0, identity; q=0.5, *;q=0": 1<<Gzip | 1<<Identity,
	"ABC, gzip":                          1<<Gzip | 1<<Identity,
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ, gzip":   1<<Gzip | 1<<Identity,
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ012345":   1 << Identity,
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456":  1 << Identity,
	"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz, gzip": 1<<Gzip | 1<<Identity,

	"zstd":     1<<Zstd | 1<<Identity,
	"zstd;q=0": 1 << Identity,

	"dcz": 1<<DCZ | 1<<Identity,
	"dcb": 1<<DCB | 1<<Identity,

	"dcz;q=0": 1 << Identity,
	"dcb;q=0": 1 << Identity,
}

func TestParse(t *testing.T) {

	n := testing.AllocsPerRun(1, func() {
		for acceptEncoding, expected := range parseTests {
			got := Parse(acceptEncoding)
			if expected != got {
				t.Errorf("Parse(%q): expected %q, got %q", acceptEncoding, expected, got)
			}
		}
	})
	if testing.CoverMode() == "" && n != 0 {
		t.Errorf("expected Parse() not to allocate, allocated %v times", n)
	}
}
