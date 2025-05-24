package handlers

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// AuthRequest represents authentication request
// @Description Authentication request payload
type AuthRequest struct {
	Name             string `json:"name" validate:"required" example:"Иван"`
	Surname          string `json:"surname" validate:"required" example:"Иванов"`
	UniversityNumber int    `json:"university_number" validate:"required" example:"182551"`
}

type AuthResponce struct {
	Token string `json:"token" example:"a3eb9104-bc88-45b4-b250-87e2ba0b5e65"`
}

// @Summary Авторизация/регистрация пользователя
// @Description Авторизация/регистрация пользователя и получение токена
// @Tags Авторизация
// @Accept json
// @Produce json
// @Success 200 {object} AuthResponce
// @Param request body AuthRequest true "Auth credentials"
// @Router /api/v1/auth [post]
func (h *AuthHandler) Auth(ctx *fiber.Ctx) error {
	var requestData AuthRequest
	err := ctx.BodyParser(&requestData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "Неверный запрос",
			},
		)
	}
	validate := validator.New()
	err = validate.Struct(requestData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": err.Error(),
			},
		)
	}
	userCondition := User{
		UniversityNumber: requestData.UniversityNumber,
	}
	var user User
	err = h.db.Where(userCondition).
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = userCondition
		user.Name = requestData.Name
		user.Surname = requestData.Surname
		user.Token = uuid.NewString()
		err := h.db.Create(&user).Error
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"error": err.Error(),
				},
			)
		}
	} else if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"error": "Ошибка при работе с БД",
			},
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"token": user.Token,
		},
	)
}

// проверка доступа пользователю
func (h *AuthHandler) AuthMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("x-api-key")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Отсутствует токен доступа в headers.Authorization",
		})
	}
	var user User
	err := h.db.Where(&User{Token: authHeader}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Вы не авторизованы, токен не валидный",
		})
	}
	return ctx.Next()
}
