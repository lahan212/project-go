package pkg

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware is a middleware to authenticate users
var jwtSecret = []byte("your_secret_key")

// AuthMiddleware is a middleware to authenticate users using JWT
func AuthMiddleware(expectedRole string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Implement your authentication logic here

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"code":    fiber.StatusUnauthorized,
				"message": "Missing Authorization header",
			})
		}
		// Expecting Bearer token, so check the format
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Bearer token format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// Verify the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// Check the role from JWT claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get claims"})
		}

		// Ensure the user has the expected role
		if role, ok := claims["role"].(string); !ok || role != expectedRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient privileges"})
		}
		// Extract user data from the claims
		username := claims["username"].(string)
		role := claims["role"].(string)

		// Set user data in the context
		c.Locals("userData", fiber.Map{"username": username, "role": role})
		return c.Next()
	}
}
