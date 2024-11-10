package controllers

import (
	"bytes"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"provinsi/utils"
)

func UploadPDF(c *fiber.Ctx) error {
	// Ambil file dari form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("File not found")
	}

	// Cek apakah file tersebut berformat PDF
	if file.Header.Get("Content-Type") != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).SendString("Only PDF files are allowed")
	}

	// Baca file ke dalam buffer
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
	}
	defer fileContent.Close()

	// Buat buffer untuk menyimpan konten file
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(fileContent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read file")
	}
	// Validasi file PDF langsung dari buffer
	err = utils.ValidatePDF(buf.Bytes())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Invalid PDF: %v", err))
	}

	return c.SendString("PDF uploaded and validated successfully!")
}
