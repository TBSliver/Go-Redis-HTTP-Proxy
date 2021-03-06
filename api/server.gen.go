// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /get/{key})
	GetGetKey(ctx echo.Context, key string) error

	// (POST /set/{key})
	PostSetKey(ctx echo.Context, key string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetGetKey converts echo context to params.
func (w *ServerInterfaceWrapper) GetGetKey(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "key" -------------
	var key string

	err = runtime.BindStyledParameterWithLocation("simple", false, "key", runtime.ParamLocationPath, ctx.Param("key"), &key)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter key: %s", err))
	}

	ctx.Set(ApiKeyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetGetKey(ctx, key)
	return err
}

// PostSetKey converts echo context to params.
func (w *ServerInterfaceWrapper) PostSetKey(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "key" -------------
	var key string

	err = runtime.BindStyledParameterWithLocation("simple", false, "key", runtime.ParamLocationPath, ctx.Param("key"), &key)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter key: %s", err))
	}

	ctx.Set(ApiKeyScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostSetKey(ctx, key)
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

	router.GET(baseURL+"/get/:key", wrapper.GetGetKey)
	router.POST(baseURL+"/set/:key", wrapper.PostSetKey)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xUTW/bMAz9KwK3wwYYdtbefNwHuqGHFnWwSxAMis0kmmNJFemiRuD/PlB2lqR1gO2w",
	"3WRafHx8j9QeStd4Z9EyQb6HgI8tEn90lcEYeMDKUMHB2I18ls4yWpaj9n5nSs3G2ewnOSsxKrfYaDm9",
	"DbiGHN5kR/xs+EvZKWbf9wlUSGUwXqAgh+9616JipwhZ8RZVjZ1il6pP2qoVKq0o5ioXlG2bFQZIInET",
	"sIKcQ4sCGpC8szS0cVcXbVki0cMY/atm8Fk3foeQw90tJMCdlzNd7OCuVofqyllFQ2mQi6MGE9Iei6x2",
	"iH6iTgKEZRsMd4WgDJ1pb26xk5OR0o8thg4SsLqRXO3Njxq7I9h4PZI2du0k8Zx8pKW+zuf36j6455hr",
	"OBKb+PWEgYa8WfohvYI+AefRam8gh+t0ll5DAl7zNpLNNsjZvsaul68N8uvyN6PpDYqzpIwVw5FTtRjK",
	"33yZq8+upOW7LbOnPMuCxFPjZNYabSuSMu8hMgnR02/VgHyDfBvV8DroBhkDQb4YpROWR+UG1c7HKjmZ",
	"iZfmLF9M3NVs9r8W5gG5DZaibE9xedYuHDYnzg3rjTQKAwbBUoIZnZrhHU24UYxuDLBOvKix++1F8Qde",
	"0JQX9464+MdmHF6y7pK8Z4/ducjTZk6DjPey12/MBelPFjk2fFjhxbJf9r8CAAD//5u1XKmSBQAA",
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
