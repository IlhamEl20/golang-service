package controllers

import (
	"provinsi/database"
	"provinsi/models"

	"github.com/gofiber/fiber/v2"
)

// Get all provinces
func GetAllProvinsi(c *fiber.Ctx) error {
	var provinsis []models.Provinsi
	database.DB.Find(&provinsis)
	return c.JSON(provinsis)
}

// Get cities based on selected province
func GetKotasByProvinsi(c *fiber.Ctx) error {
	provinsiID := c.Params("provinsiID")
	var kotas []models.Kota
	database.DB.Where("provinsi_id = ?", provinsiID).Find(&kotas)
	return c.JSON(kotas)
}

// Get kecamatan based on selected city
func GetKecamatansByKota(c *fiber.Ctx) error {
	kotaID := c.Params("kotaID")
	var kecamatans []models.Kecamatan
	database.DB.Where("kota_id = ?", kotaID).Find(&kecamatans)
	return c.JSON(kecamatans)
}

// Get kelurahan based on selected kecamatan
func GetKelurahansByKecamatan(c *fiber.Ctx) error {
	kecamatanID := c.Params("kecamatanID")
	var kelurahans []models.Kelurahan
	database.DB.Where("kecamatan_id = ?", kecamatanID).Find(&kelurahans)
	return c.JSON(kelurahans)
}
