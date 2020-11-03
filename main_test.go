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
	jsonStr := []byte(`{"name":"xyz","value":"123", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}`)
	request := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(200), response.Code, "Should be succeed")
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	assert.Equal(t, true, strings.Contains(string(bodyString), "OK"))
	result, _ := ioutil.ReadFile("./result")
	assert.Equal(t, true, strings.Contains(string(result), "123"))

	jsonStr = []byte(`{"name":"xyz","value":"abc-456", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}`)
	request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler = http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(200), response.Code, "Should be succeed")
	bodyBytes, _ = ioutil.ReadAll(response.Body)
	bodyString = string(bodyBytes)
	assert.Equal(t, true, strings.Contains(string(bodyString), "OK"))
	result, _ = ioutil.ReadFile("./result")
	assert.Equal(t, true, strings.Contains(string(result), "abc-456"))

	jsonStr = []byte(`{"name":"xyz","value":"abc_456", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}`)
	request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler = http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(200), response.Code, "Should be succeed")
	bodyBytes, _ = ioutil.ReadAll(response.Body)
	bodyString = string(bodyBytes)
	assert.Equal(t, true, strings.Contains(string(bodyString), "OK"))
	result, _ = ioutil.ReadFile("./result")
	assert.Equal(t, true, strings.Contains(string(result), "abc_456"))
}

func TestTriggerArbitraryCodeExecution(t *testing.T) {
	jsonStr := []byte(`{"name":"xyz","value":"456 && echo 'sip' > /tmp/ok", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}`)
	request := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(400), response.Code, "Should be fail")

	jsonStr = []byte(`{"name":"xyz","value":"|", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}`)
	request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler = http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(400), response.Code, "Should be fail")

	jsonStr = []byte(`{"name":"xyz","value":"\", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}`)
	request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler = http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(400), response.Code, "Should be fail")

	jsonStr = []byte(`{"name":"xyz","value":">", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}`)
	request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler = http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(400), response.Code, "Should be fail")

	jsonStr = []byte(`{"name":"xyz","value":"a b", "token": "RYWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLp7A"}`)
	request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response = httptest.NewRecorder()
	handler = http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(400), response.Code, "Should be fail")
}

func TestTriggerWithInvalidToken(t *testing.T) {
	jsonStr := []byte(`{"name":"xyz","value":"456", "token": "XXWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLpXX"}`)
	request := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(401), response.Code, "Should be failed")
}

func TestTriggerWithInvalidDeploymentName(t *testing.T) {
	jsonStr := []byte(`{"name":"xxyzz","value":"456", "token": "XXWKkMSGK7tCb7jCSVZNmJzWneNDb2funq6kSLUPDVCgL8gAMPBfUWLyKtQdLpXX"}`)
	request := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler := http.HandlerFunc(handler)
	handler.ServeHTTP(response, request)
	assert.Equal(t, int(401), response.Code, "Should be failed")
}
