package main

import (
	"net"
	"net/http"
	"time"
)

// NewNetClient initializes a new `http.Client` with the specified `timeout`
// as the timeout for the `Dial`, `TLSHandshake`, `IdleConnTimeout`, `ResponseHeaderTimeout` & `ExpectContinueTimeout.
func NewNetClient(timeout time.Duration) (netClient NetworkClient) {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: timeout,
		}).Dial,
		TLSHandshakeTimeout:   timeout,
		IdleConnTimeout:       timeout,
		ResponseHeaderTimeout: timeout,
		ExpectContinueTimeout: timeout,
	}
	netClient = &http.Client{
		Timeout:   timeout,
		Transport: netTransport,
	}
	return
}
