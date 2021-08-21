package encoding

import (
	"strconv"
	"strings"

	"github.com/renthraysk/encoding/ascii"
)

// Encoding represents a specific encoding/compression algorithm
type Encoding uint8

const (
	Identity Encoding = iota
	Deflate
	Compress
	Gzip
	Brotli
)

const (
	maxEncodingNameLen = len("identity") // NOTE: potentially needs to be updated if new encodings added.
)

func (e Encoding) String() string {
	const names string = "identity" + "deflate" + "compress" + "gzip" + "br"
	const index string = "\x00\x08\x0F\x17\x1B\x1D"

	if e <= Brotli {
		return names[index[e]:index[e+1]]
	}
	return "encoding(" + strconv.FormatUint(uint64(e), 10) + ")"
}

/* EncodingSet a set of Encodings */
type EncodingSet uint32

const (
	allSet EncodingSet = 1<<Gzip | 1<<Deflate | 1<<Brotli | 1<<Identity | 1<<Compress
)

// Contains returns true if EncodingSet es contains Encoding e, false otherwise.
func (es EncodingSet) Contains(e Encoding) bool {
	return e < 32 && (1<<e)&es != 0
}

func (es EncodingSet) String() string {
	var s strings.Builder

	for e := Identity; e <= Brotli; e++ {
		if !es.Contains(e) {
			continue
		}
		if s.Len() > 0 {
			s.WriteString(", ")
		}
		s.WriteString(e.String())
	}
	return s.String()
}

// encodingSet returns relevant EncodingSet value if it's recognised, 0 and false otherwise.
func encodingSet(name string) (EncodingSet, bool) {
	n := ascii.HorizontalSpace.TrimString(name)
	if len(n) > maxEncodingNameLen {
		// avoid the ToLowerString call and a possible pointless allocation
		// for a name that is too long to match any known encoding.
		return 0, false
	}
	switch ascii.ToLowerString(n) {
	case "identity":
		return 1 << Identity, true
	case "deflate":
		return 1 << Deflate, true
	case "compress":
		return 1 << Compress, true
	case "gzip":
		return 1 << Gzip, true
	case "br":
		return 1 << Brotli, true
	case "*":
		return allSet, true
	}
	return 0, false
}

// Parse parses an Accept-Encoding client request header and returns the set of supported encodings
func Parse(acceptEncoding string) EncodingSet {
	// Some common Accept-Encodings hardcoded to avoid actual parsing
	switch acceptEncoding {
	case "", "*":
		return allSet
	case "gzip, deflate":
		return 1<<Gzip | 1<<Deflate | 1<<Identity
	case "gzip, deflate, br":
		return 1<<Gzip | 1<<Deflate | 1<<Brotli | 1<<Identity
	}

	var supported, unsupported EncodingSet

	for len(acceptEncoding) > 0 {
		var part, name, pair string

		part, acceptEncoding = ascii.CutString(acceptEncoding, ',')
		name, part = ascii.CutString(part, ';')
		if encoding, ok := encodingSet(name); ok {
			for len(part) > 0 {
				pair, part = ascii.CutString(part, ';')
				pair = ascii.HorizontalSpace.TrimString(pair)
				param, value := ascii.CutString(pair, '=')
				if ascii.HorizontalSpace.TrimString(param) != "q" {
					continue
				}
				if len(value) > 0 {
					q, err := strconv.ParseFloat(value, 32)
					if err != nil || q <= 0 {
						unsupported |= encoding
						encoding = 0
					}
				}
				break
			}
			supported |= encoding
		}
	}
	return ((1 << Identity) &^ unsupported) | supported
}
