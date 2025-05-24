package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "app/docs"
	"app/handlers"
)

// @title API для студентов (Data provider)
// @version 1.0
// @description
// Здесь представлено описание всех endpoint сервиса для студентов
// @securityDefinitions.apikey x-api-key
// @in header
// @name x-api-key
// @BasePath /api/v1
// @host localhost:3000
func main() {
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	db.AutoMigrate(&handlers.User{})
	if err != nil {
		panic(err)
	}
	handlers.NewAuthHandler(db)
	app := fiber.New()
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	authHandler := handlers.NewAuthHandler(db)
	v1.Post("/auth", authHandler.Auth)

	v1.Group("/", authHandler.AuthMiddleware)
	v1.Get("/temperature", handlers.TemperatureHandler)
	v1.Get("/pressure", handlers.PressureHandler)
	v1.Get("/results", func(ctx *fiber.Ctx) error {
		var users []handlers.User
		if err := db.Find(&users).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch data",
			})
		}

		var response []map[string]interface{}
		for _, user := range users {
			response = append(response, map[string]interface{}{
				"id":        user.ID,
				"username":  user.Name,
				"surname":   user.Surname,
				"createdAt": user.CreatedAt.Format(time.RFC3339),
			})
		}

		return ctx.JSON(fiber.Map{"data": response})
	})

	app.Listen(":3000")
}
