package gorocks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type HandlerFunc func(*App)

type App struct {
	req           *http.Request
	resp          http.ResponseWriter
	routes        map[string]routeBlock
	middlewares   []Middleware
	routeRegexMap map[string]*regexp.Regexp
	curr          *routeBlock
}

func New() *App {
	return &App{
		routes:        make(map[string]routeBlock),
		middlewares:   make([]Middleware, 0),
		routeRegexMap: make(map[string]*regexp.Regexp),
	}
}

func (app *App) set404() {
	if app.resp == nil {
		return
	}

	app.resp.WriteHeader(http.StatusNotFound)
}

func (app *App) processPath(path string) (string, map[int]string) {
	urlParams := make(map[int]string)
	if path == emptyString || path == slash {
		return path, urlParams
	}

	segments := strings.Split(path, "/")
	if len(segments) == 0 {
		return path, urlParams
	}

	generalPath := emptyString
	count := 0

	for _, segment := range segments {
		if segment == emptyString {
			count++
			continue
		}

		if segment[0] == colonRune {
			generalPath += slashAsterisk
			urlParams[count] = segment[1:]
			count++
			continue
		}

		count++
		generalPath += slash + segment
	}

	return generalPath, urlParams
}

func (app *App) addToRegexMap(processedPath string) {
	if !strings.Contains(processedPath, asterisk) {
		return
	}

	regexString := caret + strings.Replace(processedPath, asterisk, alphanumericRegex, negativeOne)

	regexString = strings.Replace(regexString, slash, backwardForwardRegexSlash, negativeOne)
	regexString = regexString + dollar

	compiledRegex := regexp.MustCompile(regexString)
	app.routeRegexMap[processedPath] = compiledRegex
}

func (app *App) prefixFix(path string) string {
	if path == emptyString {
		return path
	}

	if path[0] != slashRune {
		path = slash + path
	}

	return path
}

func (app *App) suffixFix(path string) string {
	if path == emptyString {
		return path
	}

	if path[len(path)-1] == slashRune {
		path = path[:len(path)-1]
	}

	return path
}

func (app *App) pathFix(path string) string {
	if path == emptyString || path == slash {
		return path
	}

	path = app.suffixFix(path)
	path = app.prefixFix(path)

	if path == emptyString {
		return slash
	}

	return path
}

func (app *App) getMethodWithPath(path, method string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

func (app *App) addPath(path, method string, fn HandlerFunc) {
	path = app.pathFix(path)
	processedPath, urlParamOrderMap := app.processPath(path)
	processedPath = app.getMethodWithPath(processedPath, method)
	app.addToRegexMap(processedPath)
	app.routes[processedPath] = routeBlock{
		method:        method,
		fn:            fn,
		urlParamOrder: urlParamOrderMap,
		urlParams:     map[string]string{},
	}
}

func (app *App) getRouteBlock(path string) (*routeBlock, bool) {
	rb, okay := app.routes[path]
	if okay {
		return &rb, false
	}

	for urlPath, regexp := range app.routeRegexMap {
		if regexp.MatchString(path) {
			rb := app.routes[urlPath]
			return &rb, true
		}
	}

	return nil, false
}

func (app *App) HttpRequest() *http.Request {
	return app.req
}

func (app *App) ResponseWriter() http.ResponseWriter {
	return app.resp
}

func (app *App) SetStatusCode(code int) {
	if app.resp == nil {
		return
	}

	app.resp.WriteHeader(code)
}

func (app *App) SetHeader(key, value string) {
	if app.resp == nil {
		return
	}

	app.resp.Header().Set(key, value)
}

func (app *App) PrintRoutes() {
	for i := range app.routes {
		fmt.Println(i)
	}
}

func (app *App) Use(fns ...Middleware) {
	app.middlewares = append(app.middlewares, fns...)
}

func (app *App) Bind(fn func(reader io.Reader) error) error {
	if app.req == nil {
		return ErrRequestNil
	}

	if app.req.Body == nil {
		return ErrRequestBodyNil
	}

	return fn(app.req.Body)
}

func (app *App) BindJson(obj interface{}) error {
	fn := func(reader io.Reader) error {
		body, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		return json.Unmarshal(body, obj)
	}

	return app.Bind(fn)
}

func (app *App) JSON(statusCode int, obj interface{}) error {
	app.SetStatusCode(statusCode)

	responseBytes, err := json.MarshalIndent(obj, emptyString, "\t")
	if err != nil {
		return err
	}

	_, err = app.resp.Write(responseBytes)

	return err
}

func (app *App) Param(key string) string {
	return app.curr.urlParams[key]
}

func (app *App) Header(key string) string {
	if app.req == nil {
		return emptyString
	}

	if app.req.Header == nil {
		return emptyString
	}

	return app.req.Header.Get(key)
}

func (app *App) Query(key string) string {
	if app.req == nil {
		return emptyString
	}

	if app.req.URL == nil {
		return emptyString
	}

	return app.req.URL.Query().Get(key)
}

func (app *App) Get(path string, fn HandlerFunc) {
	app.addPath(path, http.MethodGet, fn)
}

func (app *App) Post(path string, fn HandlerFunc) {
	app.addPath(path, http.MethodPost, fn)
}

func (app *App) Delete(path string, fn HandlerFunc) {
	app.addPath(path, http.MethodDelete, fn)
}

func (app *App) Put(path string, fn HandlerFunc) {
	app.addPath(path, http.MethodPut, fn)
}

func (app *App) Patch(path string, fn HandlerFunc) {
	app.addPath(path, http.MethodPatch, fn)
}

func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.req = req
	app.resp = res
	path := app.pathFix(req.URL.Path)
	path = app.getMethodWithPath(path, req.Method)
	rb, isDynamicPath := app.getRouteBlock(path)
	if rb == nil {
		app.set404()
		return
	}

	if rb.method != req.Method {
		app.set404()
		return
	}

	app.curr = rb
	if isDynamicPath {
		rb.setUrlParams(path)
	}

	handler := rb.fn
	for _, ml := range app.middlewares {
		handler = ml(handler)
	}

	handler(app)
}

func (app *App) Run(portString string) error {
	err := http.ListenAndServe(portString, app)
	if err != nil {
		return err
	}

	return nil
}
