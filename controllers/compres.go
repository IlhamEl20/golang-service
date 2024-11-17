package controllers

import (
	"bytes"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	model "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func CompressPDF(c *fiber.Ctx) error {
	// Terima file PDF dari request
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to receive file: " + err.Error())
	}

	in, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to open uploaded file: " + err.Error())
	}
	defer in.Close()

	// Baca file PDF ke dalam buffer
	var inputBuffer bytes.Buffer
	if _, err := io.Copy(&inputBuffer, in); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read uploaded file: " + err.Error())
	}

	// Kompres file PDF
	var outputBuffer bytes.Buffer
	config := model.NewDefaultConfiguration()
	err = api.Optimize(bytes.NewReader(inputBuffer.Bytes()), &outputBuffer, config)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to compress PDF: " + err.Error())
	}

	// Kirim file PDF yang sudah dikompresi sebagai response
	return c.SendStream(bytes.NewReader(outputBuffer.Bytes()), outputBuffer.Len())
}
