package handlers

import (
	"PrescriptionDispensingSystem/internals/cache"
	"PrescriptionDispensingSystem/internals/db"
	"PrescriptionDispensingSystem/internals/models"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AddMedicine godoc
// @Summary Add a new medicine
// @Description Admin endpoint: adds or updates medicine stock
// @Tags admin
// @Accept json
// @Produce json
// @Param medicine body models.Medicine true "Medicine details"
// @Success 200 {object} models.Medicine
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/addMedicine [post]

func AddMedicine(c *fiber.Ctx) error {
	var med models.Medicine
	if err := c.BodyParser(&med); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:medicine:%s", med.Name)

	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "system busy, try again later",
			})
		case <-ticker.C:
			acquired, err := cache.AcquireLock(lockKey, 5*time.Second)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "lock acquisition failed"})
			}
			if acquired {
				defer cache.ReleaseLock(lockKey)
				goto ATOMIC
			}
		}
	}

ATOMIC:
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM medicines WHERE name=$1)`
	err := db.DB.QueryRow(ctx, query, med.Name).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if exists {
		updateQuery := `UPDATE medicines 
		                SET stock_quantity = stock_quantity + $1, updated_at = NOW() 
		                WHERE name=$2 
		                RETURNING stock_quantity`
		err := db.DB.QueryRow(ctx, updateQuery, med.StockQuantity, med.Name).Scan(&med.StockQuantity)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		insertQuery := `INSERT INTO medicines (name, dosage_form, stock_quantity, created_at, updated_at) 
		                VALUES ($1,$2,$3,NOW(),NOW()) 
		                RETURNING name,dosage_form,stock_quantity`
		err := db.DB.QueryRow(ctx, insertQuery, med.Name, med.DosageForm, med.StockQuantity).
			Scan(&med.Name, &med.DosageForm, &med.StockQuantity)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(med)
}

// DeleteMedicine godoc
// @Summary Delete a medicine by name
// @Description Admin endpoint: deletes a medicine from the inventory
// @Tags admin
// @Accept json
// @Produce json
// @Param body body struct{name=string} true "Medicine name"
// @Success 200 {string} string "Medicine Deleted"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "internal server error"
// @Router /admin/deleteMedicine [delete]

func DeleteMedicine(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON("Invalid input")
	}

	query := `DELETE FROM medicines WHERE LOWER(name) = LOWER($1)`
	res, err := db.DB.Exec(context.Background(), query, body.Name)
	if err != nil {
		return c.JSON(err.Error())
	}

	if rows := res.RowsAffected(); rows == 0 {
		return c.JSON("No medicine found with that name")
	}

	return c.JSON("Medicine Deleted")
}
