package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	api "github.com/teamdigitale/api-starter-kit-go/api"
)

func DoJson(handler http.Handler, method string, url string, body interface{}) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(b))
	handler.ServeHTTP(rr, req)

	return rr
}

func ErrorHandler(w http.ResponseWriter, err string, code int) {
	http.Error(w, err, code)
}

type Error struct {
	Message string `json:"message,omitempty"`
}

var app *api.MyApplication
var h http.Handler

func TestMain(m *testing.M) {
	store := api.CreateApplication()
	h = api.HandlerCustom(store)
	// run tests
	code := m.Run()
	os.Exit(code)
}

func Test404(t *testing.T) {
	var err error
	result := DoJson(h, "GET", "/missing", nil)
	assert.Equal(t, http.StatusNotFound, result.Code)
	return
	var error api.Problem
	var bytes []byte
	bytes = result.Body.Bytes()
	err = json.Unmarshal(bytes, &error)
	assert.NoError(t, err, "Cannot parse response", err)
	fmt.Println("Error response: ", error, bytes)
}

func TestEcho(t *testing.T) {
	// We should get a 404 on invalid ID
	result := DoJson(h, "GET", "/echo", nil)
	bytes := result.Body

	assert.Equal(t, http.StatusOK, result.Code)
	var ts api.Timestamps
	err := json.NewDecoder(result.Body).Decode(&ts)
	assert.NoError(t, err, "error parsing response", bytes)

}

func TestStatus(t *testing.T) {
	result := DoJson(h, "GET", "/status", nil)
	bytes := result.Body

	assert.Equal(t, http.StatusOK, result.Code)
	var ts api.Problem
	err := json.NewDecoder(result.Body).Decode(&ts)
	assert.NoError(t, err, "error parsing response", bytes)
}
