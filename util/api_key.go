package util

import (
    "bytes"
    "encoding/json"
    "github.com/TikhonP/medsenger-freestylelibre-bot/config"
    "github.com/labstack/echo/v4"
    "io"
    "net/http"
)

type apiKeyModel struct {
    ApiKey string `json:"api_key" validate:"required"`
}

func (k *apiKeyModel) isValid(cfg *config.Server) bool {
    return cfg.MedsengerAgentKey == k.ApiKey
}

func ApiKeyJSON(cfg *config.Server) echo.MiddlewareFunc {
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

func ApiKeyGetParam(cfg *config.Server) echo.MiddlewareFunc {
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
