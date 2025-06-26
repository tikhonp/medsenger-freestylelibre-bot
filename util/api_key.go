// Package util provides utility functions for entire application.
package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-freestylelibre-bot/config"
)

type apiKeyModel struct {
	APIKey string `json:"api_key" validate:"required"`
}

func (k *apiKeyModel) isValid(cfg *config.Server) bool {
	return cfg.MedsengerAgentKey == k.APIKey
}

func APIKeyJSON(cfg *config.Server) echo.MiddlewareFunc {
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

			data := new(apiKeyModel)
			if err := json.Unmarshal(bodyBytes, &data); err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON.")
			}
			if err := c.Validate(data); err != nil {
				return err
			}
			if !data.isValid(cfg) {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key.")
			}
			return next(c)
		}
	}
}

func APIKeyGetParam(cfg *config.Server) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.QueryParam("api_key")
			if apiKey != cfg.MedsengerAgentKey {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid API key.")
			}
			return next(c)
		}
	}
}
