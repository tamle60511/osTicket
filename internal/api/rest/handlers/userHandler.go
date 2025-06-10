package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/dto"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	svc service.UserService
}

func SetupUserRoutes(rh *rest.RestHandler) {
	app := rh.App
	svc := service.UserService{
		Repo: repository.NewRepository(rh.DB),
		Auth: rh.Auth,
	}
	handler := UserHandler{
		svc: svc,
	}
	pubRoutes := app.Group("/users")
	pubRoutes.Post("/register", handler.Register)
	pubRoutes.Post("/login", handler.Login)
	pvtroutes := pubRoutes.Group("/", rh.Auth.Authorize)
	pvtroutes.Get("/profile", handler.GetProfile)
	pvtroutes.Post("/profile", handler.CreateProfile)

}

func (h *UserHandler) Register(c fiber.Ctx) error {
	user := dto.UserSignup{}

	if err := c.Bind().Body(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	if user.Email == "" || user.Password == "" || user.UserName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Email, password and username are required",
		})
	}

	token, err := h.svc.Signup(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Registration successful",
		"token":   token,
	})
}

func (h *UserHandler) Login(c fiber.Ctx) error {
	loginInput := dto.UserLogin{}

	if err := c.Bind().Body(&loginInput); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}

	if loginInput.Email == "" || loginInput.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Email and password are required",
		})
	}

	token, err := h.svc.Login(loginInput.Email, loginInput.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid email or password",
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User logged in successfully",
		"token":   token,
	})
}

func (h *UserHandler) CreateProfile(c fiber.Ctx) error {
	return c.Status(http.StatusNotImplemented).JSON(fiber.Map{
		"message": "Create profile logic not implemented",
	})
}

func (h *UserHandler) GetProfile(c fiber.Ctx) error {
	user, err := h.svc.Auth.GetCurrentUser(c)
	if err != nil {
		log.Println("Error getting current user:", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization failed",
			"reason":  err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Profile retrieved successfully",
		"user": fiber.Map{
			"id":        user.ID,
			"email":     user.Email,
			"user_name": user.UserName,
			"role":      user.Role,
		},
	})
}
