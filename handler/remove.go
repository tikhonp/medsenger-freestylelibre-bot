package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
)

type RemoveHandler struct{}

func (h RemoveHandler) Handle(c echo.Context) error {
	contractID := c.Get("contract_id").(int)
	if err := db.MarkInactiveContractWithID(contractID); err != nil {
		return err
	}
	return c.String(http.StatusCreated, "ok")
}
