package util

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// QueryParamInt returns the value of the query parameter with the specified name as an int.
func QueryParamInt(c echo.Context, name string) *int {
	param := c.QueryParam(name)
	result, err := strconv.Atoi(param)
	if err != nil {
		return nil
	}
	return &result
}
