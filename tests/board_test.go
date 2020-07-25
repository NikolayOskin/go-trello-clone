package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUserBoards(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"GET",
		"/users/me/boards",
		bytes.NewBuffer([]byte(``)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}

func TestCreateBoardSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "testuser@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"POST",
		"/boards",
		bytes.NewBuffer([]byte(`{"title":"My new board"}`)),
	)
	if response.StatusCode != http.StatusCreated {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusCreated)
	}
}

func TestUpdateBoardSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"PUT",
		"/boards/507f191e810c19729de860ea", // getting id from seeded board
		bytes.NewBuffer([]byte(`{"title":"Updated title"}`)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}

func TestGetFullBoard(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"GET",
		"/boards/507f191e810c19729de860ea", // getting id from seeded board
		bytes.NewBuffer([]byte(``)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}

func TestGetFullBoardNotExisted(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"GET",
		"/boards/507f191e810c19729de860eb", // getting id from seeded board
		bytes.NewBuffer([]byte(``)),
	)
	if response.StatusCode != http.StatusNotFound {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusNotFound)
	}
}
