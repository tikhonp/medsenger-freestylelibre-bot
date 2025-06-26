package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
)

type contractIDModel struct {
	ContractID int `json:"contract_id" validate:"required"`
}

type RemoveHandler struct{}

func (h RemoveHandler) Handle(c echo.Context) error {
	m := new(contractIDModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	if err := db.MarkInactiveContractWithID(m.ContractID); err != nil {
		return err
	}
	return c.String(http.StatusCreated, "ok")
}
