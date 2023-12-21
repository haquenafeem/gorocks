package gorocks

import (
	"fmt"
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
