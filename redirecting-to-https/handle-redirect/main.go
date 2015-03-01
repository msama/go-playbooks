package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Factory method for helloWorldHandler
func NewDummyHandler() http.Handler {
	return &dummyHandler{}
}

// An http.Handler which says hello world
type dummyHandler struct{}

func (h *dummyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "You are being served with https (hopefully).\n")
}

// Factory method for redirectHandler
func NewRedirectHandler(scheme, host, port string) http.Handler {
	return &redirectHandler{
		scheme: scheme,
		host:   host,
		port:   port,
	}
}

// An http.Handler which redirects a request by overriding certain url params
type redirectHandler struct {
	scheme string
	host   string
	port   string
}

func (h *redirectHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	newUrl := *req.URL
	if len(h.scheme) > 0 {
		newUrl.Scheme = h.scheme
	}

	newHostAndPort := make([]string, 2)
	// Sadly the host is in the request and not in the url
	hostAndPort := strings.Split(req.Host, ":")
	copy(newHostAndPort, hostAndPort)

	if len(h.host) > 0 {
		newHostAndPort[0] = h.host
	}

	if len(h.port) > 0 {
		newHostAndPort[1] = h.port
	}

	if len(newHostAndPort[1]) > 0 {
		newUrl.Host = newHostAndPort[0] + ":" + newHostAndPort[1]
	} else {
		newUrl.Host = newHostAndPort[0]
	}
	http.Redirect(rw, req, newUrl.String(), http.StatusMovedPermanently)
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}

	port_tls := os.Getenv("PORT_TLS")
	if len(port_tls) == 0 {
		port_tls = "443"
	}

	errors := make(chan error, 1)

	go func() {
		// Redirect at runtime all the traffic to https.
		if err := http.ListenAndServe(":"+port, NewRedirectHandler("https", "", port_tls)); err != nil {
			errors <- err
		}
	}()

	go func() {
		// The certificate was generated with:
		// openssl genrsa -out test_private_key 2048
		// openssl req -new -x509 -key test_private_key -out test_cert -days 365
		if err := http.ListenAndServeTLS(":"+port_tls, "test_cert.pem", "test_private_key.pem",
			NewDummyHandler()); err != nil {
			errors <- err
		}
	}()

	log.Printf("Service up")
	log.Printf("HTTP listening on %v\n", port)
	log.Printf("HTTPS listening on %v\n", port_tls)

	select {
	case err := <-errors:
		// Handles the error from http.ListenAndServe
		log.Printf("Error: %v\n", err)
		break
	}
}
