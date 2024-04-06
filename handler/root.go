package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type RootHandler struct {
}

func (h *RootHandler) Handle(c echo.Context) error {
	return c.String(http.StatusOK, "Купил мужик шляпу, а она ему как раз!")
}
