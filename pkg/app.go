package pkg

import (
	"project-go/db"
	"project-go/models"
	"strings"

	"gorm.io/gorm"

	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const secretKey = "your-secret-key"

// LoginHandler handles the login process
func LoginHandler(c *fiber.Ctx) error {
	db.Init()
	var loginCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&loginCredentials); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate the username and password
	var user models.Users
	result := db.DB.Where("username = ?", loginCredentials.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// User not found
			// return c.SendStatus(fiber.StatusUnauthorized)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User tidak terdaftar",
				"code":  fiber.StatusNotFound,
			})
		}

		// Handle other database errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
			"code":  fiber.StatusInternalServerError,
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginCredentials.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
			"code":  fiber.StatusUnauthorized,
		})
	}
	// // Generate a Bearer Token for the new user
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"username": user.Username,
	// 	"role":     user.Role,
	// })

	// Create the Claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating JWT token",
		})
	}

	// Set the user's role in the session or token (replace this with your actual session or token logic)
	// For simplicity, setting it as a cookie in this example
	c.Cookie(&fiber.Cookie{
		Name:  "role",
		Value: user.Role,
	})

	// Redirect to a route that checks the user's role
	// return c.Redirect("/dashboard")
	// Respond with a JSON object
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
		"token":   tokenString,
	})
}

// RegisterHandler handles user registration
func RegisterHandler(c *fiber.Ctx) error {
	db.Init()
	// Parse the request JSON to get registration details
	var registrationDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
		// Add any other registration fields as needed
	}
	if err := c.BodyParser(&registrationDetails); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error hashing the password",
		})
	}

	// Create a new user record in the database
	newUser := models.Users{
		Username: registrationDetails.Username,
		Password: string(hashedPassword),
		Role:     registrationDetails.Role, // Set a default role or adjust based on your requirements
	}

	result := db.DB.Create(&newUser)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating user",
		})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"username": registrationDetails.Username,
		"role":     registrationDetails.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating JWT token",
		})
	}

	// Respond with the Bearer Token
	return c.JSON(fiber.Map{
		"status": "success",
		"token":  tokenString,
	})
}

func UserHandler(c *fiber.Ctx) error {

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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"code":    fiber.StatusUnauthorized,
			"message": "Invalid or expired token",
		})
	}

	// Decode the token and get user details
	user, err := DecodeToken(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Respond with user details
	return c.JSON(fiber.Map{
		"status": "success",
		"user":   user,
	})
}

// DecodeToken extracts user information from the token
func DecodeToken(tokenString string) (*models.Users, error) {
	// Your JWT secret key used for signing tokens
	jwtSecret := []byte("your_secret_key")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	// Retrieve user details from claims
	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}

	// Optionally, you can retrieve other user details from claims as needed

	// Create a User instance with retrieved details
	user := &models.Users{
		Username: username,
	}

	return user, nil
}
