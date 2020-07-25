package tests

import (
	"io"
	"io/ioutil"

	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUpValidation(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	cases := []struct {
		Json           string
		ExpectedStatus int
	}{
		{`{"email":"","password":""}`, 422},
		{`{"email":"some@gmail.com","password":""}`, 422},
		{`{"email":"","password":"qwertyqwerty"}`, 422},
		{`{"email":"some@gmail.com","password":"123"}`, 422},
		{`{"email":"notemail","password":"qwertyqwerty"}`, 422},
	}

	for _, testcase := range cases {
		response, _ := testRequest(t, ts, "POST", "/auth/sign-up", bytes.NewBuffer([]byte(testcase.Json)))

		if response.StatusCode != testcase.ExpectedStatus {
			t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, testcase.ExpectedStatus)
		}
	}
}

func TestSignInValidation(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	cases := []struct {
		Json           string
		ExpectedStatus int
	}{
		{`{"email":"","password":""}`, 422},
		{`{"email":"some@gmail.com","password":""}`, 422},
		{`{"email":"","password":"qwertyqwerty"}`, 422},
	}

	for _, testcase := range cases {
		response, _ := testRequest(t, ts, "POST", "/auth/sign-in", bytes.NewBuffer([]byte(testcase.Json)))

		if response.StatusCode != testcase.ExpectedStatus {
			t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, testcase.ExpectedStatus)
		}
	}
}

func TestResetPasswordValidation(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	// empty email
	if res, _ := testRequest(
		t,
		ts,
		"POST",
		"/auth/reset-password",
		bytes.NewBuffer([]byte(`{"email":""}`))); res.StatusCode != 422 {
		t.Errorf("handler returned wrong status code: got %v want %v", res.StatusCode, 422)
	}

	// string which is not email
	if res, _ := testRequest(
		t,
		ts,
		"POST",
		"/auth/reset-password",
		bytes.NewBuffer([]byte(`{"email":"string"}`))); res.StatusCode != 422 {
		t.Errorf("handler returned wrong status code: got %v want %v", res.StatusCode, 422)
	}
}

func TestSetNewPasswordValidation(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	cases := []struct {
		Json           string
		ExpectedStatus int
	}{
		{`{"email":"somebody@gmail.com","password":"","code":""}`, 422},
		{`{"email":"somebody@gmail.com","password":"qwerty","code":""}`, 422},
		{`{"email":"somebody@gmail.com","password":"short","code":"qwerty"}`, 422},
		{`{"email":"","password":"qwertyqwerty","code":"qwerty"}`, 422},
		{`{"email":"notemail","password":"qwertyqwerty","code":"qwerty"}`, 422},
	}

	for _, testcase := range cases {
		response, _ := testRequest(t, ts, "POST", "/auth/new-password", bytes.NewBuffer([]byte(testcase.Json)))

		if response.StatusCode != testcase.ExpectedStatus {
			t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, testcase.ExpectedStatus)
		}
	}
}

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	method string,
	path string,
	body io.Reader,
) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}

func testRequestWithJwt(
	accessToken string,
	t *testing.T,
	ts *httptest.Server,
	method string,
	path string,
	body io.Reader,
) (*http.Response, string) {

	req, err := http.NewRequest(method, ts.URL+path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", accessToken)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
