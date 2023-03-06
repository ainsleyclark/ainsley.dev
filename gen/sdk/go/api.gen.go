// Package sdk provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
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
	ApiKeyAuthScopes = "ApiKeyAuth.Scopes"
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
	// SendContactForm request with any body
	SendContactFormWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	SendContactForm(ctx context.Context, body SendContactFormJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
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

	operationPath := fmt.Sprintf("/forms/contact")
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
	// SendContactForm request with any body
	SendContactFormWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*SendContactFormResponse, error)

	SendContactFormWithResponse(ctx context.Context, body SendContactFormJSONRequestBody, reqEditors ...RequestEditorFn) (*SendContactFormResponse, error)
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

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /forms/contact)
	SendContactForm(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// SendContactForm converts echo context to params.
func (w *ServerInterfaceWrapper) SendContactForm(ctx echo.Context) error {
	var err error

	ctx.Set(ApiKeyAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SendContactForm(ctx)
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

	router.POST(baseURL+"/forms/contact", wrapper.SendContactForm)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xUTY/bNhD9KwTbo1baNjed6mw3qLFJY8SbosBiUYypkcWEIhnOyIkQ6L8XpGRb/mhz",
	"SE6WyZk3bzhv3lepXOudRcsky6+SVIMtpM87ZxkUv3KhfYefOiSOpxWSCtqzdlaW8g0SwRYF2EoAc9Cb",
	"jpFEHVwruEGhRgxRu9DmMpM+OI+BNaYKjbPYe3cF97FB0U7YB7COMEQQ7j3KUhIHbbdyyOQU+X0wQyYD",
	"fup0wEqWTwfM7Mjy+ZDjNh9QcSz9x+Pj6j4EF2Lx0+6UqxIl/AKtNzENY+A/6fxKF7iHOUu4FhvrwNjk",
	"PP5VZ1U8zd8gN676ZpMTlX2ZI+p/tfoOyTtLeNltBQzxd2H7v8B0MWI4Hc2R5nvCkCQBzFhdTmixWn6T",
	"eqp3pL6vc0l8yCSh6oLmfh3FPdJdeP2A/aLjJv7TUSwNQoURy0IbAf6+WayWNw/YH7lAypJDBNW2duOY",
	"k8RThy1ok5CMcb+BtmSwzyvcHUFnh+IRoZWZ7ELKYfZUFgV9hu0WQ65dfPNTOS9WS7H2qOLbpYfagPqI",
	"thKuFqfVjFY4jWkq/HL9u3hxc2egIxSvp+vz4lvNTbfJlWuLCU8ZCB+LGXixMW5TtECMoXi9vLv/c30f",
	"mTKGlt7Waww7rXCGOc9NQUV8T83m/DnGqe8w0NjtL/ltfjuJ3YLXspQv0lEmPXCTxlhEFVExm4F313xq",
	"5YhJgLD4OQlPULdpNcVKk+zcwas8bDGfb8OykqVco61mjihHNSLxS1f1ex2gTbXBe6NVyi0+0Liko7HG",
	"r58D1rKUPxVH5y0m2y2ueG4S26WtpR2KrAMCYybAGFFrNBUJCCj2qxJDPBCJHRhdJUq5TKs0LnJ6xV9v",
	"b39YBycucYX72wdxI9ZoWVCnFBLVnTH9qPUaOsM/lMlozVdovLf4xaOK9oP7mLlVyPLp1CSenofneB2i",
	"QNPtmcKCq7pkvuLe7nRwto0dnG/YfBuipofsHGjNsNV2+78oNMbkF2jPw78BAAD//y6+TmrRBwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
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
	var res = make(map[string]func() ([]byte, error))
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
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
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
