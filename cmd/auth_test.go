package cmd

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthValid(t *testing.T) {
	gdc := getDefaultClient
	defer func() {
		getDefaultClient = gdc
	}()

	getDefaultClient = func(req *http.Request) (*http.Response, error) {
		rr := httptest.NewRecorder()
		rr.WriteString("{\"token\":\"token\",\"refresh_token\":\"refresh_token\",\"id\":6,\"name\":\"Elton Minetto\"}")
		rr.WriteHeader(http.StatusOK)
		return rr.Result(), nil
	}

	err := doLogin("user", "password")
	assert.EqualValues(t, nil, err, "Error on valid auth")
}

func TestAuthInvalid(t *testing.T) {
	gdc := getDefaultClient
	defer func() {
		getDefaultClient = gdc
	}()

	getDefaultClient = func(req *http.Request) (*http.Response, error) {
		rr := httptest.NewRecorder()
		rr.WriteString("User or Password is invalid. Please verify your credentials.")
		resp := rr.Result()
		resp.StatusCode = http.StatusNotFound
		return resp, nil
	}

	err := doLogin("user", "password")
	assert.EqualValues(t, "Invalid credentials", err.Error(), "Error on invalid auth")
}
