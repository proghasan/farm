package middleware

import (
	"farm/config"
	"farm/database"
	"farm/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Login    string `json:"login"` // email, phone, or username
	Password string `json:"password"`
}

type AuthClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func Login(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		var user models.User
		result := database.DB.Where("email = ? OR phone = ? OR username = ?", req.Login, req.Login, req.Login).First(&user)
		if result.Error != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		if user.Status != "Active" {
			return c.Status(403).JSON(fiber.Map{"error": "Account is not active"})
		}

		now := time.Now()
		claims := AuthClaims{
			UserID: user.ID,
			Role:   user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(72 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		database.DB.Model(&user).Update("last_login_at", now)

		return c.JSON(fiber.Map{
			"token": signedToken,
			"user": fiber.Map{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"phone": user.Phone,
				"role":  user.Role,
			},
		})
	}
}

func GetUserID(c *fiber.Ctx) uint {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(AuthClaims)
	return claims.UserID
}

func GetUserRole(c *fiber.Ctx) string {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(AuthClaims)
	return claims.Role
}
