package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	// http.Handle("/", &rootHandler{})
	if err := http.ListenAndServe(":8080", &rootHandler{}); err != nil {
		log.Fatal(err)
	}
}

type rootHandler struct{}

func (h *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"Method":           r.Method,           // string
		"URL":              r.URL,              // *url.URL
		"Proto":            r.Proto,            // string
		"ProtoMajor":       r.ProtoMajor,       // int
		"ProtoMinor":       r.ProtoMinor,       // int
		"Header":           r.Header,           // Header
		"Body":             r.Body,             // io.ReadCloser
		"GetBody":          r.GetBody,          // func() (io.ReadCloser, error)
		"ContentLength":    r.ContentLength,    // int64
		"TransferEncoding": r.TransferEncoding, // []string
		"Close":            r.Close,            // bool
		"Host":             r.Host,             // string
		"Form":             r.Form,             // url.Values
		"PostForm":         r.PostForm,         // url.Values
		"MultipartForm":    r.MultipartForm,    // *multipart.Form
		"Trailer":          r.Trailer,          // Header
		"RemoteAddr":       r.RemoteAddr,       // string
		"RequestURI":       r.RequestURI,       // string
		"TLS":              r.TLS,              // *tls.ConnectionState
		"Response":         r.Response,         // *Response
	}).Debug()
}
