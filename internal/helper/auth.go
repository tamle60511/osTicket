package helper

import (
	"ecommerce/internal/domain"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

func SetupAuth(s string) Auth {
	return Auth{
		Secret: s,
	}
}

func (a Auth) CreateHashedPassword(p string) (string, error) {
	if len(p) < 6 {
		return "", errors.New("password length should be at least 6 characters long")
	}
	hashP, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", errors.New("password hash fail")
	}
	return string(hashP), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {
	if id == 0 || email == "" || role == "" {
		return "", errors.New("required inputs are missing to generate token")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", errors.New("unable to sign the token")
	}
	return tokenStr, nil
}

func (a Auth) VerifyPassword(pP string, hP string) error {
	if len(pP) < 6 {
		return errors.New("password length should be at least 6 characters long")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hP), []byte(pP))
	if err != nil {
		return errors.New("password does not match")
	}
	return nil
}
func IsValidRole(role string) bool {
	switch domain.Role(role) {
	case domain.Admin, domain.Staff, domain.Manager, domain.IT:
		return true
	}
	return false
}

func (a Auth) VerifyToken(t string) (domain.User, error) {
	tokenStr := strings.Split(t, " ")
	if len(tokenStr) != 2 || tokenStr[0] != "Bearer" {
		return domain.User{}, errors.New("invalid token format")
	}

	actualToken := tokenStr[1]
	token, err := jwt.Parse(actualToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown signing method %v", token.Header["alg"])
		}
		return []byte(a.Secret), nil
	})
	if err != nil {
		log.Printf("Failed to parse the token: %v", err)
		return domain.User{}, errors.New("failed to parse token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); !ok || float64(time.Now().Unix()) > exp {
			return domain.User{}, errors.New("token is expired")
		}

		user := domain.User{}
		if id, ok := claims["user_id"].(float64); ok {
			user.ID = uint(id)
		} else {
			return domain.User{}, errors.New("user ID not found in token")
		}
		if email, ok := claims["email"].(string); ok {
			user.Email = email
		} else {
			return domain.User{}, errors.New("email not found in token")
		}
		if role, ok := claims["role"].(string); ok {
			if !IsValidRole(role) {
				return domain.User{}, errors.New("invalid role found in token")
			}
			user.Role = domain.Role(role)
		} else {
			return domain.User{}, errors.New("role not found in token")
		}

		return user, nil
	}

	return domain.User{}, errors.New("token verification failed")
}

func (a Auth) Authorize(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  "missing authorization header",
		})
	}

	user, err := a.VerifyToken(authHeader)
	if err != nil || user.ID <= 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"reason":  "invalid or expired token",
		})
	}

	c.Locals("user", user)
	return c.Next()
}

func (a Auth) GetCurrentUser(c fiber.Ctx) (domain.User, error) {
	user := c.Locals("user")
	u, ok := user.(domain.User)
	if !ok {
		return domain.User{}, errors.New("user not found in context")
	}
	return u, nil
}
