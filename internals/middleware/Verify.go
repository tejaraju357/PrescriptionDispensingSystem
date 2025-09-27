package middleware

import (
	"PrescriptionDispensingSystem/internals/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Verification(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, utils.ExtractSecertKey)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userId, ok := claims["userID"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}

	role, ok := claims["role"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}

	validRoles := map[string]bool{
		"admin": true,
		"doctor":   true,
		"pharmacist": true,
	}

	roleStr, ok := role.(string)
	if !ok || !validRoles[roleStr] {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role"})
	}

	c.Locals("user_id", userId)
	c.Locals("role", role)
	return c.Next()
}


func AdminVerification(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, utils.ExtractSecertKey)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userId, ok := claims["userID"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}

	role, ok := claims["role"]
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token data"})
	}

	validRoles := map[string]bool{
		"admin": true,
		
	}

	roleStr, ok := role.(string)
	if !ok || !validRoles[roleStr] {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role"})
	}

	c.Locals("user_id", userId)
	c.Locals("role", role)
	return c.Next()
}

