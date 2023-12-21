package gorocks

import (
	"net/http"
	"testing"
)

func TestRouteBlock(t *testing.T) {
	rb := &routeBlock{
		method: http.MethodGet,
		urlParamOrder: map[int]string{
			1: "id",
			3: "name",
			5: "key",
		},
		urlParams: map[string]string{},
	}

	rb.setUrlParams("path/id/path2/name/path3/key")
	if rb.urlParams["id"] != "id" ||
		rb.urlParams["name"] != "name" ||
		rb.urlParams["key"] != "key" {
		t.Fail()
	}
}
