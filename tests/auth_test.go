package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NikolayOskin/go-trello-clone/app"
	"github.com/NikolayOskin/go-trello-clone/controller"
	"github.com/NikolayOskin/go-trello-clone/db"
	"github.com/NikolayOskin/go-trello-clone/db/seeder"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var router *chi.Mux

func TestMain(m *testing.M) {
	_ = os.Setenv("APP_ENV", "test")
	db.InitDB()
	db.FreshDb()
	seeder.Seed()
	a := app.New()
	a.InitRouting()
	router = a.Router
	os.Exit(m.Run())
}

func TestSignUpSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, _ := testRequest(
		t,
		ts,
		"POST",
		"/auth/sign-up",
		bytes.NewBuffer([]byte(`{"email":"test@gmail.com","password":"qwertyqwerty"}`)))

	if response.StatusCode != http.StatusCreated {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusCreated)
	}
}

func TestSignUpWithMalformedFields(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, _ := testRequest(
		t,
		ts,
		"POST",
		"/auth/sign-up",
		bytes.NewBuffer([]byte(`{"email":"test@gmail.com","password":"qwertyqwerty","field":"something"}`)))

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusBadRequest)
	}
}

func TestSignUpWithExistedUser(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, _ := testRequest(
		t,
		ts,
		"POST",
		"/auth/sign-up",
		bytes.NewBuffer([]byte(`{"email":"testuser@gmail.com","password":"qwertyqwerty"}`)))

	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusBadRequest)
	}
}

func TestSignInNonExistUser(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, body := testRequest(
		t,
		ts,
		"POST",
		"/auth/sign-in",
		bytes.NewBuffer([]byte(`{"email":"notexist@gmail.com","password":"qwertyqwerty"}`)))

	if response.StatusCode != 400 {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, 400)
	}
	if !strings.Contains(body, "invalid credentials") {
		t.Errorf("returned wrong response body")
	}
}

func TestSignInWithInvalidPassword(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, body := testRequest(
		t,
		ts,
		"POST",
		"/auth/sign-in",
		bytes.NewBuffer([]byte(`{"email":"testuser@gmail.com","password":"invalid"}`)))

	if response.StatusCode != 400 {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, 400)
	}
	if !strings.Contains(body, "invalid credentials") {
		t.Errorf("returned wrong response body")
	}
}

func TestSignInSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, bodyStr := testRequest(
		t,
		ts,
		"POST",
		"/auth/sign-in",
		bytes.NewBuffer([]byte(`{"email":"testuser@gmail.com","password":"qwerty"}`)))

	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
	if !strings.Contains(bodyStr, "access_token") {
		t.Errorf("returned wrong response body")
	}
}

func TestGetAuthUserSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "testuser@gmail.com", "qwerty")

	// test get authenticated user
	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"GET",
		"/users/me",
		bytes.NewBuffer([]byte("")),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}

func TestGetAuthUserWithInvalidToken(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	response, _ := testRequestWithJwt(
		"invalid_access_token",
		t,
		ts,
		"GET",
		"/users/me",
		bytes.NewBuffer([]byte("")),
	)
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusUnauthorized)
	}
}

func TestVerifyEmailWrongCode(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "testuser@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"PUT",
		"/auth/verify/00000000000000",
		bytes.NewBuffer([]byte(`{"email":"testuser@gmail.com"}`)),
	)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusBadRequest)
	}
}

func TestVerifiedMiddleware(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "testuser@gmail.com", "qwerty")

	// trying to fetch endpoint for verified users with unverified user token
	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"GET",
		"/users/me/boards",
		bytes.NewBuffer([]byte("")),
	)
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusUnauthorized)
	}
}

func TestJwtMiddlewareWithEmptyToken(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	// trying to fetch endpoint with empty jwt token
	response, _ := testRequestWithJwt(
		"",
		t,
		ts,
		"GET",
		"/users/me/boards",
		bytes.NewBuffer([]byte("")),
	)
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusUnauthorized)
	}
}

func TestVerifyEmailSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "testuser@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"PUT",
		"/auth/verify/12345",
		bytes.NewBuffer([]byte(`{"email":"testuser@gmail.com"}`)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}

func signIn(t *testing.T, ts *httptest.Server, email string, password string) string {
	var jwtResponse controller.JWTResponse
	_, body := testRequest(
		t,
		ts,
		"POST",
		"/auth/sign-in",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"email":"%v","password":"%v"}`, email, password))),
	)
	_ = json.Unmarshal([]byte(body), &jwtResponse)

	return jwtResponse.AccessToken
}
