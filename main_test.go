package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestTrigger(t *testing.T) {
	jsonStr := []byte(`{"name":"xyz","value":"123", "token": "abc"}`)

	request := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(200), response.Code, "Should be succeed")
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	assert.Equal(t, true, strings.Contains(string(bodyString), "123"))

	result, _ := ioutil.ReadFile("./result")
	assert.Equal(t, true, strings.Contains(string(result), "123"))
}
