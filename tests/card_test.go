package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateCardSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"POST",
		"/cards",
		bytes.NewBuffer([]byte(`{"text":"New card", "board_id": "507f191e810c19729de860ea", "list_id": "507f191e810c19729de860eb", "pos": 2}`)),
	)
	if response.StatusCode != http.StatusCreated {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusCreated)
	}
}

func TestUpdateCardSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"PUT",
		"/cards/507f191e810c19729de860eb",
		bytes.NewBuffer([]byte(`{"text":"Updated card2", "board_id": "507f191e810c19729de860ea", "list_id": "507f191e810c19729de860eb", "pos": 2}`)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}

func TestDeleteCardSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"DELETE",
		"/cards/507f191e810c19729de860eb",
		bytes.NewBuffer([]byte(``)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}
