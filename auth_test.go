package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	var cookie *http.Cookie
	{
		//login Check
		const loginPath = "http://localhost:5000/auth/login"
		reqJson := `{"username" : "123456781240", "password" : "12345678901232345345"}`
		resp, err := http.Post(loginPath, "application/json", strings.NewReader(reqJson))
		require.NoError(t, err)
		respBody, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t,
			`{"location":"/","status":"success"}`,
			string(respBody))
		cookies := resp.Cookies()
		require.Len(t, cookies, 1)
		cookie = cookies[0]
	}

	if true {
		//logout Check
		const newBlogPath = "http://localhost:5000/auth/logout"
		req, err := http.NewRequest("DELETE", newBlogPath, nil)
		require.NoError(t, err)
		req.AddCookie(cookie)
		client := http.DefaultClient
		resp, err := client.Do(req)
		require.NoError(t, err)
		respBody, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t,
			``,
			string(respBody))
	}
	{
		//cookie Check
		const newBlogPath = "http://localhost:5000/blogs/new"
		req, err := http.NewRequest("POST", newBlogPath, nil)
		require.NoError(t, err)
		req.AddCookie(cookie)
		client := http.DefaultClient
		resp, err := client.Do(req)
		require.NoError(t, err)
		respBody, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, 200)
		require.Equal(t,
			``,
			string(respBody))
	}

}

func TestRegister(t *testing.T) {
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
