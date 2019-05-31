package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewNetClient_ProperlySetsThePassedInTimeout(t *testing.T) {
	expected := 20 * time.Second

	sut := NewNetClient(expected)
	assert.Equal(t, expected, sut.Timeout, "Client.Timeout doesn't match. Got %v; expected %v", sut.Timeout, expected)
	sutTransport := sut.Transport.(*http.Transport)
	assert.Equal(t, expected, sutTransport.TLSHandshakeTimeout, "Client.Transport.TLSHandshakeTimeout doesn't match. Got %v; expected %v", sutTransport.TLSHandshakeTimeout, expected)
	assert.Equal(t, expected, sutTransport.IdleConnTimeout, "Client.Transport.IdleConnTimeout doesn't match. Got %v; expected %v", sutTransport.IdleConnTimeout, expected)
	assert.Equal(t, expected, sutTransport.ResponseHeaderTimeout, "Client.Transport.ResponseHeaderTimeout doesn't match. Got %v; expected %v", sutTransport.ResponseHeaderTimeout, expected)
	assert.Equal(t, expected, sutTransport.ExpectContinueTimeout, "Client.Transport.ExpectContinueTimeout doesn't match. Got %v; expected %v", sutTransport.ExpectContinueTimeout, expected)
}
