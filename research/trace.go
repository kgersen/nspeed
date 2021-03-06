package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
)

// transport is an http.RoundTripper that keeps track of the in-flight
// request and implements hooks to report HTTP tracing events.
type transport struct {
	current *http.Request
}

// RoundTrip wraps http.DefaultTransport.RoundTrip to keep track
// of the current request.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	return http.DefaultTransport.RoundTrip(req)
}

// GotConn prints whether the connection has been used previously
// for the current request.
func (t *transport) GotConn(info httptrace.GotConnInfo) {
	fmt.Printf("Connection reused for %v? %v\n", t.current.URL, info.Reused)
}

func (t *transport) DNSDone(dnsInfo httptrace.DNSDoneInfo) {
	fmt.Printf("DNS Info: %+v\n", dnsInfo)
}

func main() {
	t := &transport{}

	req, _ := http.NewRequest("GET", os.Args[1], nil)
	trace := &httptrace.ClientTrace{
		GotConn: t.GotConn,
		DNSDone: t.DNSDone,
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{Transport: t}
	if resp, err := client.Do(req); err != nil {
		log.Fatal(err)
	}
}
