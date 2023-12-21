package gorocks

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApplicationJSON(t *testing.T) {
	app := New()
	app.Get("/", func(a *App) {
		a.JSON(200, Map{"err": ""})
	})

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.Use(ApplicationJSON)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	value := app.ResponseWriter().Header().Get("Content-Type")
	if value != "application/json" {
		t.Fail()
	}
}

func TestResponseWithHeaders(t *testing.T) {
	app := New()
	app.Get("/", func(a *App) {
		a.JSON(200, Map{"err": ""})
	})

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app.Use(ResponseWithHeaders(map[string]string{"key": "value"}))
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	value := app.ResponseWriter().Header().Get("key")
	if value != "value" {
		t.Fail()
	}
}

var BasicAuthTestCases = []struct {
	Title        string
	User         string
	Password     string
	SetBasicAuth bool
	ErrMessage   string
	StatusCode   int
}{
	{
		Title:        `Given basic auth not set it should return not provided with 401 status`,
		SetBasicAuth: false,
		ErrMessage:   "not provided",
		StatusCode:   401,
	},
	{
		Title:        `Given basic auth provided with wrong credentials it should return wrong username/password with 401 status`,
		SetBasicAuth: true,
		ErrMessage:   "wrong username/password",
		StatusCode:   401,
	},
	{
		Title:        `Given basic auth provided with right credentials it should return empty with 200 status`,
		SetBasicAuth: true,
		ErrMessage:   "",
		User:         "un",
		Password:     "pw",
		StatusCode:   200,
	},
}

func TestBasicAuth(t *testing.T) {
	app := New()
	app.Get("/", func(a *App) {
		a.JSON(200, Map{"err": ""})
	})

	app.Use(BasicAuth("un", "pw"))

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, testCase := range BasicAuthTestCases {
		if testCase.SetBasicAuth {
			req.SetBasicAuth(testCase.User, testCase.Password)
		}

		rr := httptest.NewRecorder()
		app.ServeHTTP(rr, req)

		value := rr.Result().StatusCode
		if value != testCase.StatusCode {
			t.Fail()
		}

		var m map[string]interface{}
		body, err := io.ReadAll(rr.Result().Body)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(body, &m)

		if err != nil {
			t.Fatal(err)
		}

		if m["err"] != testCase.ErrMessage {
			t.Fail()
		}

		rr.Result().Body.Close()
	}
}
