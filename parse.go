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
	Zstd

	DCB
	DCZ
)

func (e Encoding) String() string {
	const names string = "identity" + "deflate" + "compress" + "gzip" + "br" + "zstd" + "dcb" + "dcz"
	const index string = "\x00\x08\x0F\x17\x1B\x1D\x21\x24\x27"

	if Identity <= e && e <= DCZ {
		return names[index[e]:index[e+1]]
	}
	return "encoding(" + strconv.FormatUint(uint64(e), 10) + ")"
}

/* EncodingSet a set of Encodings */
type EncodingSet uint32

const (
	allSet EncodingSet = 1<<Gzip | 1<<Deflate | 1<<Brotli | 1<<Identity | 1<<Compress | 1<<Zstd
)

// Contains returns true if EncodingSet es contains Encoding e, false otherwise.
func (es EncodingSet) Contains(e Encoding) bool {
	return e < 32 && (1<<e)&es != 0
}

func (es EncodingSet) String() string {
	var s strings.Builder

	for e := Identity; e <= DCZ; e++ {
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

const maxNameLen = max(
	len("identity"),
	len("deflate"),
	len("compress"),
	len("gzip"),
	len("br"),
	len("zstd"),
	len("dcb"),
	len("dcz"))

func ParseEncoding(name string) (Encoding, bool) {
	n := ascii.HorizontalSpace.TrimString(name)
	if len(n) > maxNameLen {
		return 0, false
	}
	switch ascii.ToLowerString(n) {
	case "identity":
		return Identity, true
	case "deflate":
		return Deflate, true
	case "compress":
		return Compress, true
	case "gzip":
		return Gzip, true
	case "br":
		return Brotli, true
	case "zstd":
		return Zstd, true
	case "dcb":
		return DCB, true
	case "dcz":
		return DCZ, true
	}
	return 0, false
}

func parseEncodingSet(name string) (EncodingSet, bool) {
	n := ascii.HorizontalSpace.TrimString(name)
	if len(n) > maxNameLen {
		return 0, false
	}
	switch ascii.ToLowerString(n) {
	case "*":
		return allSet, true
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
	case "zstd":
		return 1 << Zstd, true
	case "dcb":
		return 1 << DCB, true
	case "dcz":
		return 1 << DCZ, true
	}
	return 0, false
}

// ParseAcceptEncoding parses an Accept-Encoding client request header and returns the set of supported encodings
func ParseAcceptEncoding(acceptEncoding string) EncodingSet {
	// Some common Accept-Encodings hardcoded to avoid actual parsing
	switch acceptEncoding {
	case "", "*":
		return allSet
	case "gzip, deflate":
		return 1<<Gzip | 1<<Deflate | 1<<Identity
	case "gzip, deflate, br":
		return 1<<Gzip | 1<<Deflate | 1<<Brotli | 1<<Identity
	case "gzip, deflate, br, zstd":
		return 1<<Gzip | 1<<Deflate | 1<<Brotli | 1<<Zstd | 1<<Identity
	case "gzip, deflate, br, zstd, dcb, dcz":
		return 1<<Gzip | 1<<Deflate | 1<<Brotli | 1<<Zstd | 1<<DCB | 1<<DCZ | 1<<Identity
	}

	var supported, unsupported EncodingSet

	for len(acceptEncoding) > 0 {
		var part, name, pair string

		part, acceptEncoding = ascii.CutString(acceptEncoding, ',')
		name, part = ascii.CutString(part, ';')
		if encoding, ok := parseEncodingSet(name); ok {
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
