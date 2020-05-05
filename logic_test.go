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

const (
	baseUrl = "http://localhost:5000"
	//baseUrl = "https://stormy-hamlet-15402.herokuapp.com"
)

func TestLogicCycle(t *testing.T) {
	var IINs = [4]string{"223456781240", "223456781241", "223456781242", "223456781243"}
	var macs = [4]string{"12345678897", "12345678898", "12345678899", "12345678890"}
	var cookies [4]*http.Cookie
	for i, IIN := range IINs {
		cookies[i] = registerOrLoginWithIIN(t, IIN, macs[i])
	}
	require.Equal(t, len(cookies), len(IINs))

	addInteraction(t, cookies[0], macs[1])
	checkIsInfected(t, cookies[0], false)
	addInfected(t, IINs[1])
	checkIsInfected(t, cookies[0], true)
}

func addInfected(t *testing.T, IIN string) {
	listBefore := getInfectedListRemote(t)
	{
		//Add infected
		const newInteraction = baseUrl + "/infected/new"
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
	const newInteraction = baseUrl + "/infected/list"
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
	const newInteraction = baseUrl + "/interactions/new"
	reqJson := fmt.Sprintf(`{"mac" : "%s"}`, withIIN)
	req, err := http.NewRequest("POST", newInteraction, strings.NewReader(reqJson))
	require.NoError(t, err)

	req.AddCookie(cookie)
	client := http.DefaultClient
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, resp.StatusCode, 200, string(respBody))
}

func checkIsInfected(t *testing.T, cookie *http.Cookie, shouldBeInfected bool) {
	{
		//check is infected
		const newInteraction = baseUrl + "/interactions/status"
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

func registerOrLoginWithIIN(t *testing.T, IIN, mac string) *http.Cookie {
	{
		//Register User
		const registerPath = baseUrl + "/auth/register"
		reqJson := fmt.Sprintf(`{"IIN": "%s", "password": "12345678901232345345", "mac": "%s", "username" : "%s"}`, IIN, mac, IIN)
		_, err := http.Post(registerPath, "application/json", strings.NewReader(reqJson))
		require.NoError(t, err)
		//require.Equal(t, resp.StatusCode, 200)
		//require.Equal(t, resp.StatusCode, 200)
		//respBody, err := ioutil.ReadAll(resp.Body)
		//require.NoError(t, err)
		//require.Equal(t,
		//	`{"location":"/","status":"success"}`,
		//	string(respBody))
	}

	//login Check
	const loginPath = baseUrl + "/auth/login"
	reqJson := fmt.Sprintf(`{"username" : "%s", "password" : "12345678901232345345"}`, IIN)
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
