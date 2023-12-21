package gorocks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	app := New()

	if app == nil {
		t.Fail()
	}
}

func TestHttpRequest(t *testing.T) {
	app := New()
	req, err := http.NewRequest(http.MethodGet, "testurl", nil)
	if err != nil {
		t.Fail()
	}

	app.ServeHTTP(nil, req)

	httpRequest := app.HttpRequest()
	if httpRequest.Method != http.MethodGet {
		fmt.Println("method")
		t.Fail()
	}

	if httpRequest.URL.Path != "testurl" {
		fmt.Println("path", httpRequest.URL.RawPath)
		t.Fail()
	}
}

func TestResponseWriter(t *testing.T) {
	app := New()
	req, err := http.NewRequest(http.MethodGet, "testurl", nil)
	if err != nil {
		t.Fail()
	}

	rr := httptest.NewRecorder()
	rr.Header().Set("key", "value")
	app.ServeHTTP(rr, req)

	respW := app.ResponseWriter()
	if respW.Header().Get("key") != "value" {
		fmt.Println("method")
		t.Fail()
	}
}

func TestPathFix(t *testing.T) {
	app := New()
	path := "/abc/"
	pathFixd := app.pathFix(path)

	if pathFixd != "/abc" {
		t.Fail()
	}

	path = ""
	pathFixd = app.pathFix(path)

	if path != pathFixd {
		t.Fail()
	}
}

func TestProcessPath(t *testing.T) {
	app := New()
	path, _ := app.processPath("/path/:id/:name")
	if path != "/path/*/*" {
		t.Fail()
	}

	path, _ = app.processPath("/")
	if path != "/" {
		t.Fail()
	}
}

func TestAddRegexToMap(t *testing.T) {
	app := New()
	app.addToRegexMap("/*")

	regex := app.routeRegexMap["/*"]
	if regex.String() != "^\\/[A-Za-z0-9]([A-Za-z0-9_-]*[A-Za-z0-9])?$" {
		t.Fail()
	}
}

func TestHeader(t *testing.T) {
	app := New()
	req, err := http.NewRequest(http.MethodGet, "testurl", nil)
	if err != nil {
		t.Fail()
	}

	req.Header.Add("Content-Type", "application/json")
	app.ServeHTTP(nil, req)
	value := app.Header("Content-Type")
	if value != "application/json" {
		t.Fail()
	}
}

func TestQuery(t *testing.T) {
	app := New()
	app.Post("/testquery", func(a *App) {
		a.JSON(200, Map{"q1": a.Query("q1"), "q2": a.Query("q2")})
	})

	req, err := http.NewRequest(http.MethodPost, "/testquery?q1=q1&q2=q2", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fail()
	}

	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	var m map[string]interface{}
	body, err := io.ReadAll(rr.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(body, &m)

	if err != nil {
		t.Fatal(err)
	}

	if m["q1"] != "q1" || m["q2"] != "q2" {
		t.Fail()
	}
}

func TestBindJSON(t *testing.T) {
	app := New()
	app.Post("/testquery", func(a *App) {
		a.JSON(200, Map{"q1": a.Query("q1"), "q2": a.Query("q2")})
	})

	req, err := http.NewRequest(http.MethodPost, "/testquery?q1=q1&q2=q2", bytes.NewBuffer([]byte(`{"key":"value"}`)))
	if err != nil {
		t.Fail()
	}

	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	m := struct {
		Key string `json:"key"`
	}{}
	err = app.BindJson(&m)
	if err != nil {
		t.Fatal(err)
	}

	if m.Key != "value" {
		t.Fail()
	}

	req, err = http.NewRequest(http.MethodPost, "/testquery?q1=q1&q2=q2", nil)
	if err != nil {
		t.Fail()
	}

	rr = httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	m = struct {
		Key string `json:"key"`
	}{}

	err = app.BindJson(&m)
	if err != ErrRequestBodyNil {
		t.Fatal(err)
	}
}
