package main

import (
	"encoding/binary"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Router() *mux.Router {
	route := mux.NewRouter()
	route.HandleFunc("/", mainPage)
	route.HandleFunc("/login", login).Methods("POST")
	route.HandleFunc("/login/pass", changePass).Methods("POST")
	route.HandleFunc("/task", doWork).Methods("POST")

	return route
}

func TestMainPage(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestLoginOK(t *testing.T) {
	request := httptest.NewRequest("POST", "/login", nil)
	request.Form = make(url.Values)
	request.Form.Add("login", "admin")
	request.Form.Add("pass", "1000000")

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response for login user")
}

func TestLoginBadRequest(t *testing.T) {
	request := httptest.NewRequest("POST", "/login", nil)
	request.Form = make(url.Values)
	request.Form.Add("login", "admin")
	request.Form.Add("pass", "1000001")

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)

	assert.Equal(t, 400, response.Code, "Bad request for login user")
}

func TestChangePass(t *testing.T) {

	request := httptest.NewRequest("POST", "/login/pass", nil)
	request.Form = make(url.Values)
	request.Form.Add("login", "admin")
	request.Form.Add("pass", "1000000")
	request.Form.Add("newPass", "1000001")

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)

	assert.Equal(t, 200, response.Code, "OK response for change password")
}

func TestDoWork(t *testing.T) {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(18446744073709451615))
	result = append(result, "mnbvcxz"...)

	json_byte := []byte(`{"bigNumber": 100000, "text": "zxcvbnm"}`)
	request := httptest.NewRequest("POST", "/task", nil)
	request.Form = make(url.Values)
	request.Form.Add("login", "admin")
	request.Form.Add("value", string(json_byte))

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)

	assert.Equal(t, response.Body.Bytes(), result, "Wrong reverse result")
}
