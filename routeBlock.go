package gorocks

import (
	"strings"
)

type routeBlock struct {
	method        string
	fn            HandlerFunc
	urlParamOrder map[int]string
	urlParams     map[string]string
}

func (r *routeBlock) setUrlParams(path string) {
	segments := strings.Split(path, "/")[1:]

	for i, value := range r.urlParamOrder {
		r.urlParams[value] = segments[i-1]
	}
}
