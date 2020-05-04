package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogicCycle(t *testing.T) {
	var IINs = [4]string{"123456781240", "123456781241", "123456781242", "123456781243"}
	var cookies [4]*http.Cookie
	for i, IIN := range IINs {
		cookies[i] = registerOrLoginWithIIN(t, IIN)
	}
	require.Equal(t, len(cookies), len(IINs))

	addInteraction(t, cookies[0], IINs[1])
	checkIsInfected(t, cookies[0], false)
	addInfected(t, IINs[1])
	checkIsInfected(t, cookies[0], true)
}

func addInfected(t *testing.T, IIN string) {
	listBefore := getInfectedListRemote(t)
	{
		//Add infected
		const newInteraction = "http://localhost:5000/infected/new"
		reqJson := fmt.Sprintf(`{"IIN" : "%s"}`, IIN)
		req, err := http.NewRequest("POST", newInteraction, strings.NewReader(reqJson))
		require.NoError(t, err)

		client := http.DefaultClient
		resp, err := client.Do(req)
		require.NoError(t, err)

		_, err = ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, 200)
	}
	listAfter := getInfectedListRemote(t)
	require.Equal(t, len(listBefore)+1, len(listAfter))
}

func getInfectedListRemote(t *testing.T) []string {
	var list ResponseList

	//Add infected
	const newInteraction = "http://localhost:5000/infected/list"
	req, err := http.NewRequest("GET", newInteraction, nil)
	require.NoError(t, err)

	client := http.DefaultClient
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
	str := string(respBody)
	require.NotEqual(t, "", str)
	require.NoError(t, json.Unmarshal(respBody, &list))

	return list.List
}

func TestUnmarshal(t *testing.T) {
	str := `["123456781241","123456781241","123456781241","123456781241","123456781241"]`
	res := make([]string, 1)
	require.NoError(t, json.Unmarshal([]byte(str), res))
}

func addInteraction(t *testing.T, cookie *http.Cookie, withIIN string) {
	//cookie Check
	const newInteraction = "http://localhost:5000/interactions/new"
	reqJson := fmt.Sprintf(`{"IIN" : "%s"}`, withIIN)
	req, err := http.NewRequest("POST", newInteraction, strings.NewReader(reqJson))
	require.NoError(t, err)

	req.AddCookie(cookie)
	client := http.DefaultClient
	resp, err := client.Do(req)
	require.NoError(t, err)

	_, err = ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200)
}

func checkIsInfected(t *testing.T, cookie *http.Cookie, shouldBeInfected bool) {
	{
		//check is infected
		const newInteraction = "http://localhost:5000/interactions/status"
		req, err := http.NewRequest("POST", newInteraction, nil)
		require.NoError(t, err)

		req.AddCookie(cookie)
		client := http.DefaultClient
		resp, err := client.Do(req)
		require.NoError(t, err)

		respBody, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, resp.StatusCode, 200)
		require.Equal(t, strconv.FormatBool(shouldBeInfected), string(respBody))
	}
}

func registerOrLoginWithIIN(t *testing.T, IIN string) *http.Cookie {
	{
		//Register User
		const registerPath = "http://localhost:5000/auth/register"
		reqJson := fmt.Sprintf(`{"IIN": "%s", "password": "12345678901232345345", "phone": "12345678897", "username" : "%s"}`, IIN, IIN)
		_, err := http.Post(registerPath, "application/json", strings.NewReader(reqJson))
		require.NoError(t, err)
	}

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
	return cookies[0]
}
