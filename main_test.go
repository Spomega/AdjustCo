package main

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeRequestToServer(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	response, _ := http.Get(server.URL)

	if response.StatusCode != 200 {
		t.Errorf("HTTP Status expected: 200, got: %d", response.StatusCode)
	}

}

func TestMakeRequestWithURL(t *testing.T) {

	response, _ := http.Get("http://adjust.com")

	if response.StatusCode != 200 {
		t.Errorf("HTTP Status expected: 200, got: %d", response.StatusCode)
	}

}

func TestMD5Hash(t *testing.T) {

	byteHash := []byte("adjust.com")
	hash := md5.Sum(byteHash)
	expected := hex.EncodeToString(hash[:])

	actual := HashResponse(byteHash)

	if actual != expected {
		t.Errorf("Expected the MD5Hash to be %s but instead got %s", expected, actual)
	}

}
