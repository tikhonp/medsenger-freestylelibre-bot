package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
	"github.com/tikhonp/medsenger-freestylelibre-bot/util"
	"github.com/tikhonp/medsenger-freestylelibre-bot/view"
)

type SettingsHandler struct{}

func (h SettingsHandler) renderPage(c echo.Context, showAddAccount bool) error {
	contractID, err := util.GetContractID(c)
	if err != nil {
		return err
	}
	contract, err := db.GetContractByID(contractID)
	if err != nil {
		return err
	}
	var lc *db.LibreClient
	lc, _ = contract.LibreClient()
	return util.TemplRender(
		c,
		view.Settings(util.Safe(contract.PatientName, ""), lc, showAddAccount),
	)
}

func (h SettingsHandler) Get(c echo.Context) error {
	return h.renderPage(c, true)
}

type userCredentials struct {
	Email    string `form:"email" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func (h SettingsHandler) Post(c echo.Context) error {
	var uc userCredentials
	if err := c.Bind(&uc); err != nil {
		return err
	}
	if err := c.Validate(uc); err != nil {
		return err
	}
	contractID, err := util.GetContractID(c)
	if err != nil {
		return err
	}
	if _, err := db.NewLibreClient(uc.Email, uc.Password, contractID); err != nil {
		return err
	}
	return h.renderPage(c, false)
}
