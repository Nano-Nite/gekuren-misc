package v1

import (
	"gakuren-system.com/pkg/helper"
	"gakuren-system.com/pkg/model"
	"github.com/gofiber/fiber/v3"
)

func SetupMiscRoute(app *fiber.App, apiVersion string) {
	miscBaseURL := apiVersion + "/misc"

	app.Post(miscBaseURL+"/gender", func(c fiber.Ctx) error {
		payload := new(model.SearchPayload)
		tenantUUID, _, err := helper.ValidateRequest(c)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusUnauthorized, "Missing or invalid authentication data", nil, err)
		}
		if err = c.Bind().Body(payload); err != nil {
			return helper.ReturnResponse(c, fiber.StatusBadRequest, "Invalid request body format", nil, err)
		}
		data, stats, err := helper.SearchGender(tenantUUID.String(), *payload)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusBadRequest, "Failed to search users", nil, err)
		}
		return helper.ReturnResponse(c, fiber.StatusOK, "success", map[string]any{"data_statistic": stats, "result": data}, nil)
	})

	app.Post(miscBaseURL+"/title", func(c fiber.Ctx) error {
		payload := new(model.SearchPayload)
		tenantUUID, _, err := helper.ValidateRequest(c)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusUnauthorized, "Missing or invalid authentication data", nil, err)
		}
		if err = c.Bind().Body(payload); err != nil {
			return helper.ReturnResponse(c, fiber.StatusBadRequest, "Invalid request body format", nil, err)
		}
		data, stats, err := helper.SearchTitle(tenantUUID.String(), *payload)
		if err != nil {
			return helper.ReturnResponse(c, fiber.StatusBadRequest, "Failed to search users", nil, err)
		}
		return helper.ReturnResponse(c, fiber.StatusOK, "success", map[string]any{"data_statistic": stats, "result": data}, nil)
	})
}
