package handlers

import (
	"PrescriptionDispensingSystem/internals/db"
	"PrescriptionDispensingSystem/internals/models"
	"PrescriptionDispensingSystem/internals/utils"
	"context"
	"errors"
	"strings"

	// "strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

// Register godoc
// @Summary Register the first admin user
// @Description Only one admin can be created. Creates an admin account.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "Admin user details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "invalid request or missing fields"
// @Failure 403 {object} map[string]string "forbidden action"
// @Failure 409 {object} map[string]string "email already exists"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /register [post]

func Register(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid request body",
			"details": err.Error(),
		})
	}

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Password = strings.TrimSpace(user.Password)
	user.Name = strings.TrimSpace(user.Name)
	user.Role = strings.TrimSpace(user.Role)

	if user.Name == "" || user.Password == "" || user.Email == "" || user.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "please fill all required fields",
		})
	}

	if user.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "only admin can create users",
		})
	}

	// Check if admin exists
	var existingAdmin string
	err := db.DB.QueryRow(context.Background(),
		"SELECT email FROM users WHERE role = 'admin'",
	).Scan(&existingAdmin)

	if err == nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "an admin already exists",
		})
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "database error",
			"details": err.Error(),
		})
	}

	// Check if email already exists
	var existingEmail string
	err = db.DB.QueryRow(context.Background(),
		"SELECT email FROM users WHERE email = $1", user.Email,
	).Scan(&existingEmail)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "email already exists",
		})
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "database query error",
			"details": err.Error(),
		})
	}

	// Insert user (unnamed statement)
	_, err = db.DB.Exec(context.Background(),
		"INSERT INTO users (name, email, role, password, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())",
		user.Name, user.Email, user.Role, user.Password,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to register user",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "admin created successfully",
		"user": fiber.Map{
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}


// CreateUsers godoc
// @Summary Create a new user (Admin only)
// @Description Admin can create users with role 'user' or other roles.
// @Tags admin
// @Accept json
// @Produce json
// @Param user body models.User true "User details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "invalid request or missing fields"
// @Failure 409 {object} map[string]string "email already exists"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /admin/createUser [post]

func CreateUsers(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid request body",
			"details": err.Error(),
		})
	}

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	user.Password = strings.TrimSpace(user.Password)
	user.Name = strings.TrimSpace(user.Name)
	user.Role = strings.TrimSpace(user.Role)

	if user.Name == "" || user.Password == "" || user.Email == "" || user.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "please fill all required details",
		})
	}

	// Check if email exists
	var existingEmail string
	err := db.DB.QueryRow(context.Background(),
		"SELECT email FROM users WHERE email = $1", user.Email,
	).Scan(&existingEmail)

	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "email already exists",
		})
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "database query error",
			"details": err.Error(),
		})
	}

	// Insert user
	_, err = db.DB.Exec(context.Background(),
		"INSERT INTO users (name, email, role, password, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW())",
		user.Name, user.Email, user.Role, user.Password,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to register user",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}


// Login godoc
// @Summary Login user
// @Description Authenticates user and returns JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body struct{Email string; Password string} true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "invalid request"
// @Failure 401 {object} map[string]string "invalid credentials"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /login [post]

func Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	var user models.User
	err := db.DB.QueryRow(context.Background(),
		"SELECT id, name, email, role, password FROM users WHERE email = $1", input.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "database query error",
			"details": err.Error(),
		})
	}

	if user.Password != input.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

