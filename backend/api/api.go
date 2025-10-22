package api

import (
	"strconv"
	"sistem-manajemen-armada/backend/db"

	"github.com/gofiber/fiber/v2"
)
// fungsi router
func SetupRouter() *fiber.App {
	app := fiber.New()
	app.Get("/vehicles/:vehicle_id/location", getLastLocation)
	app.Get("/vehicles/:vehicle_id/history", getHistory)
	return app
}
// fungsi metod get lokasi
func getLastLocation(c *fiber.Ctx) error {
	vehicleID := c.Params("vehicle_id")
	var loc db.VehicleLocation
	if result := db.DB.Where("vehicle_id = ?", vehicleID).Order("timestamp desc").First(&loc); result.Error != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}
	return c.JSON(loc)
}

// fungsi metod get histori
func getHistory(c *fiber.Ctx) error {
	vehicleID := c.Params("vehicle_id")
	start, _ := strconv.ParseInt(c.Query("start"), 10, 64)
	end, _ := strconv.ParseInt(c.Query("end"), 10, 64)
	var locs []db.VehicleLocation
	db.DB.Where("vehicle_id = ? AND timestamp BETWEEN ? AND ?", vehicleID, start, end).Find(&locs)
	return c.JSON(locs)
}
