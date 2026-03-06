package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
)

type RemoveHandler struct{}

func (h RemoveHandler) Handle(c echo.Context) error {
	contractID, err := util.GetContractID(c)
	if err != nil {
		return err
	}
	if err := db.MarkInactiveContractWithID(contractID); err != nil {
		return err
	}
	return c.String(http.StatusCreated, "ok")
}
