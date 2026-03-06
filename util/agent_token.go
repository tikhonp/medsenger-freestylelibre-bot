// Package util provides utility functions for entire application.
package util

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/maigo"
)

const contractIDKey = "contract_id"

type agentTokenModel struct {
	AgentToken string `json:"agent_token" validate:"required"`
}

func GetContractID(c echo.Context) int {
	return c.Get(contractIDKey).(int)
}

func processAgentToken(agentToken string, c echo.Context, client *maigo.Client, roles []maigo.RequestRole) error {
	data, err := client.DecodeAgentJWT(agentToken)
	if err != nil {
		log.Printf("Failed to decode agent JWT: %s", err.Error())
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid jwt key.")
	}
	for _, role := range roles {
		if !slices.Contains(data.Roles, role) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid jwt key role.")
		}
	}
	c.Set(contractIDKey, *data.ContractID)
	return nil
}

func AgentTokenJSON(client *maigo.Client, roles ...maigo.RequestRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Workaround to read request body twice
			req := c.Request()
			bodyBytes, _ := io.ReadAll(req.Body)
			if err := req.Body.Close(); err != nil {
				return err
			}
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			c.SetRequest(req)

			data := new(agentTokenModel)
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON.")
			}
			if err := c.Validate(data); err != nil {
				return err
			}
			if err := processAgentToken(data.AgentToken, c, client, roles); err != nil {
				return err
			}
			return next(c)
		}
	}
}

func AgentTokenGetParam(client *maigo.Client, roles ...maigo.RequestRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			agentToken := c.QueryParam("agent_token")
			if err := processAgentToken(agentToken, c, client, roles); err != nil {
				return err
			}
			return next(c)
		}
	}
}
