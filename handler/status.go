package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
)

type statusResponseModel struct {
	IsTrackingData     bool     `json:"is_tracking_data"`
	SupportedScenarios []string `json:"supported_scenarios"`
	TrackedContracts   []int    `json:"tracked_contracts"`
}

type StatusHandler struct{}

func (h StatusHandler) Handle(c echo.Context) error {
	trackedContracts, err := db.GetActiveContractIds()
	if err != nil {
		return err
	}
	response := statusResponseModel{
		IsTrackingData:     true,
		SupportedScenarios: []string{},
		TrackedContracts:   trackedContracts,
	}
	return c.JSON(http.StatusOK, response)
}
