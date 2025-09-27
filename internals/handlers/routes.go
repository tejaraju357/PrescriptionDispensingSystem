package handlers

import (
	"PrescriptionDispensingSystem/internals/cache"
	"PrescriptionDispensingSystem/internals/db"
	"PrescriptionDispensingSystem/internals/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetAllMedicnes godoc
// @Summary      Get all medicines
// @Description  Fetches the complete list of medicines from the database
// @Tags         medicines
// @Produce      json
// @Success      200 {array} models.Medicine
// @Failure      500 {object} map[string]string
// @Router       /medicines [get]

func GetAllMedicnes(c *fiber.Ctx) error {
	query := `SELECT name, dosage_form, stock_quantity FROM medicines`
	rows, err := db.DB.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var allMeds []models.Medicine
	for rows.Next() {
		var med models.Medicine
		if err := rows.Scan(&med.Name, &med.DosageForm, &med.StockQuantity); err != nil {
			return err
		}
		allMeds = append(allMeds, med)
	}

	return c.JSON(allMeds)
}


// DispenseStock godoc
// @Summary      Dispense medicine from stock
// @Description  Deducts stock of a specific medicine (handles concurrency with Redis locks)
// @Tags         medicines
// @Accept       json
// @Produce      json
// @Param        medicine body models.Medicine true "Medicine dispense request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /dispenseStock [post]
func DispenseStock(c *fiber.Ctx) error {
	var req models.Medicine
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "medicine name is required",
		})
	}

	ctx := context.Background()
	lockKey := fmt.Sprintf("lock:medicine:%s", req.Name)

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
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "lock acquisition failed",
				})
			}
			if acquired {
				defer cache.ReleaseLock(lockKey)
				goto DISPENSE
			}
		}
	}

DISPENSE:
	var stock int
	var medID int
	query := "SELECT id, stock_quantity FROM medicines WHERE LOWER(name) = LOWER($1)"
	err := db.DB.QueryRow(ctx, query, req.Name).Scan(&medID, &stock)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "medicine not found",
		})
	}

	if stock < req.StockQuantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("not enough stock. Available: %d", stock),
		})
	}

	updateQuery := "UPDATE medicines SET stock_quantity = stock_quantity - $1, updated_at=$2 WHERE id=$3"
	_, err = db.DB.Exec(ctx, updateQuery, req.StockQuantity, time.Now(), medID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update stock",
		})
	}

	var updatedStock int
	selectQuery := "SELECT name, dosage_form, stock_quantity FROM medicines WHERE id=$1"
	err = db.DB.QueryRow(ctx, selectQuery, medID).Scan(&req.Name, &req.DosageForm, &updatedStock)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch updated medicine info",
		})
	}
	req.StockQuantity = updatedStock

	return c.JSON(fiber.Map{
		"message":  "medicine dispensed successfully",
		"medicine": req,
	})
}

// CreatePriscription godoc
// @Summary      Create a prescription
// @Description  Stores a new prescription in the database
// @Tags         prescriptions
// @Accept       json
// @Produce      json
// @Param        prescription body models.Prescription true "Prescription data"
// @Success      200 {object} models.Prescription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /presc [post]
func CreatePriscription(c *fiber.Ctx) error {
	var presc models.Prescription
	if err := c.BodyParser(&presc); err != nil {
		return err
	}

	query := `INSERT INTO prescriptions 
	            (patient_name, medicine_name, dosage, quantity) 
	            VALUES ($1, $2, $3, $4) 
	            RETURNING id`

	var newPrescriptionID int
	if err := db.DB.QueryRow(context.Background(), query, presc.PatientName, presc.MedicineName, presc.Dosage, presc.Quantity).Scan(&newPrescriptionID); err != nil {
		return err
	}

	presc.ID = uint(newPrescriptionID)
	return c.JSON(presc)

}
