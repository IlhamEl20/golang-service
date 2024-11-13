package controllers

// import (
// 	"strconv"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/pdfcpu/pdfcpu/pkg/api"
// 	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
// )

// func AddWatermark(c *fiber.Ctx) error {
// 	text := c.Query("text")
// 	angle := c.Query("angle")

// 	// Convert angle to float64
// 	angleFloat, err := strconv.ParseFloat(angle, 64)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid angle"})
// 	}

// 	// Create a configuration for the watermark
// 	wmConf, err := pdfcpu.ParsePDFWatermarkDetails(text, "font:Helvetica, points:24, rot:"+strconv.FormatFloat(angleFloat, 'f', 2, 64), true, pdfcpu.POINTS)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create watermark configuration"})
// 	}

// 	// Apply the watermark to the PDF
// 	err = api.AddWatermarksFile("input.pdf", "", []string{"1"}, wmConf, nil)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add watermark"})
// 	}

// 	return c.SendFile("input.pdf")
// }
