package handler

import (
	"net/http"

	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
	"github.com/TikhonP/medsenger-freestylelibre-bot/util"
	"github.com/TikhonP/medsenger-freestylelibre-bot/view"
	"github.com/labstack/echo/v4"
)

type SettingsHandler struct{}

func (h SettingsHandler) renderPage(c echo.Context, showAddAccount bool) error {
	contractId := util.QueryParamInt(c, "contract_id")
	if contractId == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "contract_id is required")
	}
	contract, err := db.GetContractById(*contractId)
	if err != nil {
		return err
	}
	var lc *db.LibreClient = nil
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
	Email      string `form:"email" validate:"required"`
	Password   string `form:"password" validate:"required"`
	ContractId int    `form:"contract_id" validate:"required"`
}

func (h SettingsHandler) Post(c echo.Context) error {
	var uc userCredentials
	if err := c.Bind(&uc); err != nil {
		return err
	}
	if err := c.Validate(uc); err != nil {
		return err
	}
	if _, err := db.NewLibreClient(uc.Email, uc.Password, uc.ContractId); err != nil {
		return err
	}
    return h.renderPage(c, false)
}
