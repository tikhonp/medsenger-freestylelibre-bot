package handler

import (
	"net/http"

	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
	"github.com/TikhonP/medsenger-freestylelibre-bot/util"
	"github.com/TikhonP/medsenger-freestylelibre-bot/view"
	"github.com/labstack/echo/v4"
)

type SettingsHandler struct {
}

func (h *SettingsHandler) Get(c echo.Context) error {
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
	return util.TemplRender(c,
		view.Settings(
			util.Safe(contract.PatientName, ""),
			lc,
		),
	)
}

type formData struct {
	Email    string `form:"email" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func (h *SettingsHandler) Post(c echo.Context) error {
	contractId := util.QueryParamInt(c, "contract_id")
	if contractId == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "contract_id is required")
	}
	var userCredentials formData
	if err := c.Bind(&userCredentials); err != nil {
		return err
	}
	if err := c.Validate(userCredentials); err != nil {
		return err
	}
	_, err := db.NewLibreClient(userCredentials.Email, userCredentials.Password, *contractId)
	if err != nil {
		return err
	}
	return h.Get(c)
}
