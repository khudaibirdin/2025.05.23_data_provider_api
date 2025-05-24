package handlers

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TemperatureResponse struct {
    Value int `json:"value" example:"42"`
}

// @Summary Получить значение температуры
// @Security x-api-key
// @Description Возвращает значение температуры в градусах
// @Tags Данные
// @Produce json
// @Success 200 {object} TemperatureResponse
// @Router /api/v1/temperature [get]
func TemperatureHandler(ctx *fiber.Ctx) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ctx.JSON(
		fiber.Map{
			"value": r.Intn(100),
		},
	)
}

type PressureResponse struct {
    Value int `json:"value" example:"125"`
}

// @Summary Получить значение давления
// @Security x-api-key
// @Description Возвращает значение давления в Паскалях
// @Tags Данные
// @Produce json
// @Success 200 {object} PressureResponse
// @Router /api/v1/pressure [get]
func PressureHandler(ctx *fiber.Ctx) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ctx.JSON(
		fiber.Map{
			"value": r.Intn(100) + 100,
		},
	)
}
