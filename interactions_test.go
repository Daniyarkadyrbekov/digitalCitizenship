package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//mux.Get("/infected/list", getInfectedList(database))
//mux.Get("/infected/new", newInfetcted(database))
//mux.Get("/interactions/status", interactedWithInfected(database))
//mux.Get("/interactions/new", newInteraction(database))

func TestInfected(t *testing.T) {
	const registerPath = "http://localhost:5000/infected/new"
	for i := 0; i < 100; i++ {
		_, err := http.Get(registerPath)
		require.NoError(t, err)
	}
}

func TestStub(t *testing.T) {
	const registerPath = "http://localhost:5000/auth/register"
	reqJson := `{"IIN": "123456781240", "password": "12345678901232345345", "phone": "12345678897", "username" : "123456781240"}`
	resp, err := http.Post(registerPath, "application/json", strings.NewReader(reqJson))
	require.NoError(t, err)
	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t,
		`{"location":"/","message":"Account successfully created, you are now logged in","status":"success"}`,
		string(respBody))
}
