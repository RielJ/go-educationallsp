package rpc_test

import (
	"testing"

	"github.com/go-educationallsp/rpc"
)

type EncodingExample struct {
	Testing bool
}

func TestEncoding(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})
	if expected != actual {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestDecoding(t *testing.T) {
	msg := []byte("Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}")
	method, content, err := rpc.DecodeMessage(msg)
	contentLength := len(content)
	if err != nil {
		t.Errorf("Error decoding message: %s", err)
	}
	if contentLength != 15 {
		t.Errorf("Expected %d but got %d", 15, contentLength)
	}
	if method != "hi" {
		t.Errorf("Expected hi but got %s", method)
	}
}
