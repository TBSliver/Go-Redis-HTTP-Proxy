package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"strings"
	"tbsliver.me/armco/redis-proxy/api"
)

//go:embed openapi.html
var openapiHtml string

func main() {

	var port = flag.Int("port", 3000, "Port to listen on")
	var host = flag.String("host", "0.0.0.0", "Host to listen on")
	var redis = flag.String("redis", "redis://localhost:6379/0", "Redis server to connect to")
	var authKey = flag.String("auth", "", "Auth Key to use")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		// Ignore errors here as it's too late already
		_, _ = fmt.Fprintf(os.Stderr, "Error loading swager spec\n: %s", err)
		os.Exit(1)
	}

	redisProxy := api.NewProxy(*redis)

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Skipper: skipper(),
		Options: openapi3filter.Options{
			AuthenticationFunc: authenticator(*authKey),
		},
	}))

	e.File("/docs", "openapi.html")
	e.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, openapiHtml)
	})
	e.GET("/docs.json", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, swagger)
	})

	api.RegisterHandlers(e, redisProxy)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", *host, *port)))
}

func skipper() echoMiddleware.Skipper {
	return func(c echo.Context) bool {
		return strings.HasPrefix(c.Request().URL.Path, "/docs")
	}
}

func authenticator(authKey string) openapi3filter.AuthenticationFunc {
	return func(c context.Context, input *openapi3filter.AuthenticationInput) error {
		fmt.Print(input.SecuritySchemeName)
		if input.SecuritySchemeName != "apiKey" {
			return fmt.Errorf("security scheme %s != 'apiKey'", input.SecuritySchemeName)
		}

		q := input.RequestValidationInput.GetQueryParams()
		keys := q[input.SecurityScheme.Name]
		if len(keys) != 1 {
			return fmt.Errorf("missing key")
		}

		key := keys[0]
		if key != authKey {
			return fmt.Errorf("invalid key")
		}

		return nil
	}
}
