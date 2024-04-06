package handler

import (
	"github.com/TikhonP/medsenger-freestylelibre-bot/db"
	"github.com/TikhonP/maigo"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type initModel struct {
	ContractId        int    `json:"contract_id" validate:"required"`
	ClinicId          int    `json:"clinic_id" validate:"required"`
	AgentToken        string `json:"agent_token" validate:"required"`
	PatientAgentToken string `json:"patient_agent_token" validate:"required"`
	DoctorAgentToken  string `json:"doctor_agent_token" validate:"required"`
	AgentId           int    `json:"agent_id" validate:"required"`
	AgentName         string `json:"agent_name" validate:"required"`
	Locale            string `json:"locale" validate:"required"`
}

type InitHandler struct {
	MaigoClient *maigo.Client
}

func (h *InitHandler) Handle(c echo.Context) error {
	m := new(initModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	contract := &db.Contract{
		Id:         m.ContractId,
		IsActive:   true,
		AgentToken: &m.AgentToken,
		Locale:     &m.Locale,
	}
	if err := contract.Save(); err != nil {
		return err
	}
	go func(c db.Contract) {
		ci, err := h.MaigoClient.GetContractInfo(m.ContractId)
		if err != nil {
			log.Println(err)
			return
		}
		contract.PatientName = &ci.PatientName
		contract.PatientEmail = &ci.PatientEmail
		if err := contract.Save(); err != nil {
			log.Println(err)
			return
		}
	}(*contract)
	return c.String(http.StatusCreated, "ok")
}
