package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

const freestyleLibreScenarioObjectType = "freestylelibre_devices"

type ScenarioCapabilitiesHandler struct{}

func freestyleLibreScenarioItem() map[string]any {
	return map[string]any{
		"id":          1,
		"title":       "FreeStyle Libre",
		"description": "Мониторинг глюкозы через LibreLinkUp.",
	}
}

func (h ScenarioCapabilitiesHandler) Capabilities(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"version": 1,
		"agent": map[string]any{
			"code":  "freestylelibre",
			"title": "FreeStyle Libre",
		},
		"params": []any{},
		"object_types": []map[string]any{
			{
				"type":        freestyleLibreScenarioObjectType,
				"title":       "FreeStyle Libre",
				"description": "Подключение мониторинга глюкозы.",
				"icon":        "activity",
				"selection":   "many",
			},
		},
	})
}

func (h ScenarioCapabilitiesHandler) Objects(c *echo.Context) error {
	if c.Param("object_type") != freestyleLibreScenarioObjectType {
		return echo.NewHTTPError(http.StatusNotFound, "Object type not found.")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"type":  freestyleLibreScenarioObjectType,
		"items": []map[string]any{freestyleLibreScenarioItem()},
	})
}

func (h ScenarioCapabilitiesHandler) Object(c *echo.Context) error {
	if c.Param("object_type") != freestyleLibreScenarioObjectType || c.Param("object_id") != "1" {
		return echo.NewHTTPError(http.StatusNotFound, "Object not found.")
	}
	return c.JSON(http.StatusOK, map[string]any{
		"type":   freestyleLibreScenarioObjectType,
		"object": freestyleLibreScenarioItem(),
		"params": []any{},
	})
}
