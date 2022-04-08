package main

import (
	"flag"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"os"
	"tbsliver.me/armco/redis-proxy/api"
)

func main() {

	var port = flag.Int("port", 3000, "Port to listen on (default 3000)")
	var host = flag.String("host", "0.0.0.0", "Host to listen on (default 0.0.0.0)")
	var redis = flag.String("redis", "redis://localhost:6379/0", "Redis server to connect to (default localhost:6379)")
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

	docs := e.Group("/docs")

	docs.File("", "public/openapi.html")
	docs.File(".yaml", "api.yaml")

	e.Group("/", middleware.OapiRequestValidator(swagger))

	api.RegisterHandlers(e, redisProxy)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", *host, *port)))
}
