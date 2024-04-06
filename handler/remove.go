package handler

import (
	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

type contractIdModel struct {
	ContractId int `json:"contract_id" validate:"required"`
}

type RemoveHandler struct {
}

func (h *RemoveHandler) Handle(c echo.Context) error {
	m := new(contractIdModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	if err := db.MarkInactiveContractWithId(m.ContractId); err != nil {
		return err
	}
	return c.String(http.StatusCreated, "ok")
}
