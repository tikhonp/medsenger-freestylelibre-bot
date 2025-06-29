// Package handler provides handlers for HTTP endpoints.
package handler

import (
	"net/http"

	"github.com/TikhonP/maigo"
	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-freestylelibre-bot/db"
)

type initModel struct {
	ContractID        int    `json:"contract_id" validate:"required"`
	ClinicID          int    `json:"clinic_id" validate:"required"`
	AgentToken        string `json:"agent_token" validate:"required"`
	PatientAgentToken string `json:"patient_agent_token" validate:"required"`
	DoctorAgentToken  string `json:"doctor_agent_token" validate:"required"`
	AgentID           int    `json:"agent_id" validate:"required"`
	AgentName         string `json:"agent_name" validate:"required"`
	Locale            string `json:"locale" validate:"required"`
}

type InitHandler struct {
	MaigoClient *maigo.Client
}

func (h InitHandler) fetchContractDataOnInit(c db.Contract, ctx echo.Context) {
	ci, err := h.MaigoClient.GetContractInfo(c.ID)
	if err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	c.PatientName = &ci.PatientName
	c.PatientEmail = &ci.PatientEmail
	if err := c.Save(); err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	_, err = h.MaigoClient.SendMessage(
		c.ID,
		"Подключен агент для интеграции глюкометров freestyle libre! Пожалуйста настройте аккаунт Libre Link Up.",
		maigo.WithAction("Настроить", "/setup", maigo.Action),
		maigo.OnlyPatient(),
	)
	if err != nil {
		sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	ctx.Logger().Info("Successfully fetched contract data")
}

func (h InitHandler) Handle(c echo.Context) error {
	m := new(initModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	contract := db.Contract{
		ID:         m.ContractID,
		IsActive:   true,
		AgentToken: &m.AgentToken,
		Locale:     &m.Locale,
	}
	if err := contract.Save(); err != nil {
		return err
	}
	go h.fetchContractDataOnInit(contract, c)
	return c.String(http.StatusCreated, "ok")
}
