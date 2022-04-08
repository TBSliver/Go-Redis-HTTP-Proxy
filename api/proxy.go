package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=server.yaml ../api.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=types.yaml ../api.yaml

type Proxy struct {
	rdb *redis.Client
}

func NewProxy(url string) *Proxy {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)
	return &Proxy{
		rdb: rdb,
	}
}

type Error struct {
	Code    int32
	Message string
}

func sendError(ctx echo.Context, code int, message string) error {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, petErr)
	return err
}
