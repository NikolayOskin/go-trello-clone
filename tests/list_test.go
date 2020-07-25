package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateList(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	tests := []struct {
		name           string
		expectedStatus int
		json           string
	}{
		{
			"with all valid fields",
			http.StatusCreated,
			`{"title":"New list", "board_id": "507f191e810c19729de860ea", "pos": 2}`,
		},
		{
			"with title more than 20 chars",
			422,
			`{"title":"Very very very very very long long long title", "board_id": "507f191e810c19729de860ea", "pos": 2}`,
		},
		{
			"with empty title",
			422,
			`{"title":"", "board_id": "507f191e810c19729de860ea", "pos": 2}`,
		},
		{
			"with empty board_id",
			422,
			`{"title":"Some title", "board_id": "", "pos": 2}`,
		},
		{
			"without pos field",
			422,
			`{"title":"Some title", "board_id": "507f191e810c19729de860ea"}`,
		},
		{
			"with not existed board_id",
			http.StatusBadRequest,
			`{"title":"Some title", "board_id": "507f191e810c19729de860eb", "pos": 2}`,
		},
		{
			"with not expected fields",
			http.StatusBadRequest,
			`{"title":"Some title", "board_id": "507f191e810c19729de860ea", "pos": 2, "unknownfield":"something"}`,
		},
		{
			"with invalid json",
			http.StatusBadRequest,
			`{"title":"Some title" "board_id": "507f191e810c19729de860ea", "pos": 2, "unknownfield":"something"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			response, _ := testRequestWithJwt(
				jwtToken,
				t,
				ts,
				"POST",
				"/lists",
				bytes.NewBuffer([]byte(test.json)),
			)
			if response.StatusCode != test.expectedStatus {
				t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, test.expectedStatus)
			}
		})
	}
}

func TestUpdateList(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	tests := []struct {
		name           string
		expectedStatus int
		path           string
		json           string
	}{
		{
			"with invalid json",
			http.StatusBadRequest,
			"/lists/507f191e810c19729de860ea",
			`{"titl}`,
		},
		{
			"with title more than 20 chars",
			422,
			"/lists/507f191e810c19729de860ea",
			`{"title":"Very very very very very long long long title"}`,
		},
		{
			"with empty title",
			422,
			"/lists/507f191e810c19729de860ea",
			`{"title":""}`,
		},
		{
			"success update",
			http.StatusOK,
			"/lists/507f191e810c19729de860ea",
			`{"title":"New updated title"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			response, _ := testRequestWithJwt(
				jwtToken,
				t,
				ts,
				"PUT",
				test.path,
				bytes.NewBuffer([]byte(test.json)),
			)
			if response.StatusCode != test.expectedStatus {
				t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, test.expectedStatus)
			}
		})
	}
}

func TestDeleteListSuccess(t *testing.T) {
	ts := httptest.NewServer(router)
	defer ts.Close()

	jwtToken := signIn(t, ts, "verified@gmail.com", "qwerty")

	response, _ := testRequestWithJwt(
		jwtToken,
		t,
		ts,
		"DELETE",
		"/lists/507f191e810c19729de860ea",
		bytes.NewBuffer([]byte(``)),
	)
	if response.StatusCode != http.StatusOK {
		t.Errorf("returned wrong status code: got %v want %v", response.StatusCode, http.StatusOK)
	}
}
