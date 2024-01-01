// Package sdk provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.3 DO NOT EDIT.
package sdk

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	AuthenticationScopes = "Authentication.Scopes"
)

// ContactFormRequest Message and attributes from the contact form.
type ContactFormRequest struct {
	// Honeypot The message from the user.
	Honeypot string `json:"honeypot"`

	// Message The message from the user.
	Message string `json:"message"`
}

// HTTPError defines model for HTTPError.
type HTTPError struct {
	Code      string `json:"code"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Operation string `json:"operation"`
}

// HTTPResponse defines model for HTTPResponse.
type HTTPResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// SendContactFormJSONRequestBody defines body for SendContactForm for application/json ContentType.
type SendContactFormJSONRequestBody = ContactFormRequest

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetCredentials request
	GetCredentials(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// SendContactFormWithBody request with any body
	SendContactFormWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SendContactForm(ctx context.Context, body SendContactFormJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// Ping request
	Ping(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetCredentials(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCredentialsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SendContactFormWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSendContactFormRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) SendContactForm(ctx context.Context, body SendContactFormJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewSendContactFormRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) Ping(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPingRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetCredentialsRequest generates requests for GetCredentials
func NewGetCredentialsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/credentials/")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewSendContactFormRequest calls the generic SendContactForm builder with application/json body
func NewSendContactFormRequest(server string, body SendContactFormJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewSendContactFormRequestWithBody(server, "application/json", bodyReader)
}

// NewSendContactFormRequestWithBody generates requests for SendContactForm with any type of body
func NewSendContactFormRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/forms/contact/")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewPingRequest generates requests for Ping
func NewPingRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/ping/")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetCredentialsWithResponse request
	GetCredentialsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCredentialsResponse, error)

	// SendContactFormWithBodyWithResponse request with any body
	SendContactFormWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SendContactFormResponse, error)

	SendContactFormWithResponse(ctx context.Context, body SendContactFormJSONRequestBody, reqEditors ...RequestEditorFn) (*SendContactFormResponse, error)

	// PingWithResponse request
	PingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingResponse, error)
}

type GetCredentialsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetCredentialsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCredentialsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type SendContactFormResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *HTTPResponse
	JSONDefault  *HTTPError
}

// Status returns HTTPResponse.Status
func (r SendContactFormResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r SendContactFormResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PingResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSONDefault  *HTTPError
}

// Status returns HTTPResponse.Status
func (r PingResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PingResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetCredentialsWithResponse request returning *GetCredentialsResponse
func (c *ClientWithResponses) GetCredentialsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetCredentialsResponse, error) {
	rsp, err := c.GetCredentials(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCredentialsResponse(rsp)
}

// SendContactFormWithBodyWithResponse request with arbitrary body returning *SendContactFormResponse
func (c *ClientWithResponses) SendContactFormWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SendContactFormResponse, error) {
	rsp, err := c.SendContactFormWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSendContactFormResponse(rsp)
}

func (c *ClientWithResponses) SendContactFormWithResponse(ctx context.Context, body SendContactFormJSONRequestBody, reqEditors ...RequestEditorFn) (*SendContactFormResponse, error) {
	rsp, err := c.SendContactForm(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseSendContactFormResponse(rsp)
}

// PingWithResponse request returning *PingResponse
func (c *ClientWithResponses) PingWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PingResponse, error) {
	rsp, err := c.Ping(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePingResponse(rsp)
}

// ParseGetCredentialsResponse parses an HTTP response from a GetCredentialsWithResponse call
func ParseGetCredentialsResponse(rsp *http.Response) (*GetCredentialsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCredentialsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseSendContactFormResponse parses an HTTP response from a SendContactFormWithResponse call
func ParseSendContactFormResponse(rsp *http.Response) (*SendContactFormResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &SendContactFormResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest HTTPResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest HTTPError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParsePingResponse parses an HTTP response from a PingWithResponse call
func ParsePingResponse(rsp *http.Response) (*PingResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PingResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest HTTPError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /credentials/)
	GetCredentials(ctx echo.Context) error

	// (POST /forms/contact/)
	SendContactForm(ctx echo.Context) error

	// (GET /ping/)
	Ping(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetCredentials converts echo context to params.
func (w *ServerInterfaceWrapper) GetCredentials(ctx echo.Context) error {
	var err error

	ctx.Set(AuthenticationScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetCredentials(ctx)
	return err
}

// SendContactForm converts echo context to params.
func (w *ServerInterfaceWrapper) SendContactForm(ctx echo.Context) error {
	var err error

	ctx.Set(AuthenticationScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.SendContactForm(ctx)
	return err
}

// Ping converts echo context to params.
func (w *ServerInterfaceWrapper) Ping(ctx echo.Context) error {
	var err error

	ctx.Set(AuthenticationScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Ping(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/credentials/", wrapper.GetCredentials)
	router.POST(baseURL+"/forms/contact/", wrapper.SendContactForm)
	router.GET(baseURL+"/ping/", wrapper.Ping)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xW32/jNgz+VwRtj67d7d7ytF7XuxX3o0HTGwYUxcDIjK2rLOkkOtfg4P99oOzE+eGu",
	"2G0Pe2rDUB/Jj+THfJPKNd5ZtBTl7JuMqsYG0r+XzhIoeuNCc4tfWozE1hKjCtqTdlbO5AeMESoUYEsB",
	"REEvW8IoVsE1gmoUqscQKxeaXGbSB+cxkMYUoXYWN95N4N7VKJoBewfWRgwMQhuPciYjBW0r2WVy8Px3",
	"MF0mA35pdcBSzu53mNmY5cPujVt+RkUc+re7u/lVCC5w8MPqlCtTSvgEjTf8DNnxz2SfqAK3MEcPXqh4",
	"9E55ZKK1sDQoyAloqUZLWgFNRuRsoadqH+dNaxVb8w9ItSv/AVVDadu0R/znqLvF6J2NeMpeCQT898Ju",
	"fgfTskf3XOGfIoY0YkCE5WnHL+bXLxaR4o2pb+OcJt5lMqJqg6bNgpelT/dipHrgU/MA1gglMp6FhkH+",
	"OLuYX5+9w82YD3jNnzsG1nbl+tFJa5OqbECbhGSM+wW0jQY3eYnrEXTPKO4QGpnJNqQ3RD7OiiJ+harC",
	"kGvHvB+uyMX8Wiw8KuYvkbUE9Yi2FG4lDqMZrXBo1RD49eJX8ers0kAbUbwfvj4OXmmq22WuXFMMeMpA",
	"eCz2wIulccuigUgYivfXl1cfF1ecKWFo4s1qgWGtFe5h7r9NTgXzqckc09F3fo0h9tX+lJ/n58PoW/Ba",
	"zuSrZMqkB6pTKwsVsORegokFGyqcUKibJXEgsefcD9whabsVuC7lTL5Fuhz9JU9gP/8p8M/n59vmo00h",
	"CZ+o8Aa0HbU52Y9G+aSrN+/EmbhFChrXWIrYKoUxrlpjNjn7d5kseGFiMYxaKtS7KY2fu0hRgLD4NS2Z",
	"iO2y0ZEZHVbM7XTeQ4X5SdkLtOXeNZH95mGk167cHJUM3pthi4rP0R0V/mPAlZzJH4rxahXDySom7tUE",
	"M3wSkl5w1gGBMBNgjFhpNGUUEFBsZYFdPMQo1mB0mVLKZfdi076/ggNFfK6rC7R00NB+p1fQGvpPM+nP",
	"2kQanyw+eVQstbj14YHy2lbPL8xljeoxbn8VWFSk15o2wgVRIxiqWXEGtT4doTkP+nfvy3gq5jcf307c",
	"gmmqWRt1FLtMwPw/uN47QXJ2f3p87h+6B3YJLHzJ42ijgyvbdOLFlV3r4GzDVRwr977KslZ22THQgqDS",
	"tvpblNj75CdoD91fAQAA//8LLRQSfQoAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
