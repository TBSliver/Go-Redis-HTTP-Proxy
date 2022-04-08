package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (p Proxy) GetGetKey(ctx echo.Context, key string) error {
	value, err := p.rdb.Get(ctx.Request().Context(), key).Result()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Redis Error")
	}
	return ctx.JSON(http.StatusOK, value)
}

func (p Proxy) PostSetKey(ctx echo.Context, key string) error {
	var newString PostSetKeyJSONRequestBody
	err := ctx.Bind(&newString)
	if err != nil {
		ctx.Logger().Error(err)
		return sendError(ctx, http.StatusBadRequest, "No Value Sent")
	}
	err = p.rdb.Set(ctx.Request().Context(), key, newString, 0).Err()
	if err != nil {
		return sendError(ctx, http.StatusInternalServerError, "Redis Error")
	}
	return ctx.JSON(http.StatusOK, "ok")
}
