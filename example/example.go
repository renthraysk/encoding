package main

import (
	"compress/gzip"
	_ "embed"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/renthraysk/encoding"
)

func ServeContent(contentType string, e encoding.Encoding, content string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Accept-Encoding")
		encodingSet := encoding.Parse(r.Header.Get("Accept-Encoding"))

		if encodingSet.Contains(e) {
			writeResponseString(w, r, contentType, e, content)
			return
		}

		// Check useragent can accept identity encoding
		// if so can decompress internal gzip representation on the fly
		if encodingSet.Contains(encoding.Identity) && e == encoding.Gzip {
			gz, err := gzip.NewReader(strings.NewReader(content))
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
				return
			}
			defer gz.Close()
			writeResponseReader(w, r, contentType, encoding.Identity, gz)
			return
		}

		http.Error(w, http.StatusText(http.StatusNotAcceptable), http.StatusNotAcceptable)
	})
}

func writeResponseString(w http.ResponseWriter, r *http.Request, contentType string, e encoding.Encoding, content string) {
	setHeaders(w.Header(), contentType, e)
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	if sw, ok := w.(io.StringWriter); ok {
		sw.WriteString(content)
		return
	}
	io.CopyN(w, strings.NewReader(content), int64(len(content)))
}

func writeResponseReader(w http.ResponseWriter, r *http.Request, contentType string, e encoding.Encoding, content io.Reader) {
	setHeaders(w.Header(), contentType, e)
	w.WriteHeader(http.StatusOK)
	if r.Method == http.MethodHead {
		return
	}
	io.Copy(w, content)
}

func setHeaders(h http.Header, contentType string, e encoding.Encoding) {
	h.Set("Content-Type", contentType)
	if e != encoding.Identity {
		h.Set("Content-Encoding", e.String())
	}
}

//go:embed data/index.html.gz
var indexHTMLGz string

func main() {
	err := http.ListenAndServe(":8080", ServeContent("text/html", encoding.Gzip, indexHTMLGz))
	if err != nil {
		log.Fatalf(": %v", err)
	}
}
