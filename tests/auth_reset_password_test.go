package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	resetPasswordUrl  string = "/auth/reset-password"
	setNewPasswordUrl string = "/auth/new-password"
)

func TestResetPasswordSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, _ := testRequest(
		t,
		ts,
		"POST",
		resetPasswordUrl,
		bytes.NewBuffer([]byte(`{"email":"testuser@gmail.com"}`)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}

func TestResetPasswordInvalidEmail(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, bodyStr := testRequest(
		t,
		ts,
		"POST",
		resetPasswordUrl,
		bytes.NewBuffer([]byte(`{"email":"nonexist@gmail.com"}`)),
	)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusBadRequest)
	}
	if !strings.Contains(bodyStr, "user with this email does not exist") {
		t.Errorf("returned wrong response body")
	}
}

func TestSetNewPasswordInvalidCode(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, _ := testRequest(
		t,
		ts,
		"POST",
		setNewPasswordUrl,
		bytes.NewBuffer([]byte(`{"email":"withresetpasswordcode@gmail.com","password":"newpass","code":"invalid"}`)),
	)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusBadRequest)
	}
}

func TestSetNewPasswordWithNonExistEmail(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, _ := testRequest(
		t,
		ts,
		"POST",
		setNewPasswordUrl,
		bytes.NewBuffer([]byte(`{"email":"notexist@gmail.com","password":"newpass","code":"123456"}`)),
	)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusBadRequest)
	}
}

func TestSetNewPasswordWithExpiredCode(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, bodyStr := testRequest(
		t,
		ts,
		"POST",
		setNewPasswordUrl,
		bytes.NewBuffer([]byte(`{"email":"withexpiredresetpasswordcode@gmail.com","password":"newpass","code":"123456"}`)),
	)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusBadRequest)
	}
	if !strings.Contains(bodyStr, "code is expired") {
		t.Errorf("returned wrong response body")
	}
}

func TestSetNewPasswordSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, _ := testRequest(
		t,
		ts,
		"POST",
		setNewPasswordUrl,
		bytes.NewBuffer([]byte(`{"email":"withresetpasswordcode@gmail.com","password":"newpass","code":"123456"}`)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}
