package controllers

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func MergePDF(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to retrieve form data")
	}

	files := form.File["files"]
	if len(files) < 2 {
		return c.Status(fiber.StatusBadRequest).SendString("Please upload at least two PDF files")
	}

	var inputBuffers []*bytes.Buffer
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to open uploaded file")
		}
		defer file.Close()

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to read uploaded file")
		}
		inputBuffers = append(inputBuffers, buf)
	}

	var inputReaders []io.Reader
	for _, buf := range inputBuffers {
		inputReaders = append(inputReaders, bytes.NewReader(buf.Bytes()))
	}

	outputBuffer := new(bytes.Buffer)
	var inputFilePaths []string
	for i, buf := range inputBuffers {
		inputFilePath := fmt.Sprintf("input%d.pdf", i)
		err = os.WriteFile(inputFilePath, buf.Bytes(), 0644)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to write temporary PDF file")
		}
		inputFilePaths = append(inputFilePaths, inputFilePath)
		defer os.Remove(inputFilePath)
	}

	err = api.Merge(inputFilePaths[0], inputFilePaths[1:], outputBuffer, nil, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to merge PDF files")
	}

	c.Attachment("merged.pdf")
	return c.SendStream(bytes.NewReader(outputBuffer.Bytes()))
}

// MergePDFs menghandle penggabungan PDF
type MergeInput struct {
	Files []string `json:"files"`
}

func MergePDFs(c *fiber.Ctx) error {
	var input MergeInput

	// Parse input JSON
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request")
	}

	// Membaca file yang diterima
	var inputBuffers []bytes.Buffer
	for _, file := range input.Files {
		buf, err := os.ReadFile(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to read file")
		}
		inputBuffers = append(inputBuffers, *bytes.NewBuffer(buf))
	}

	// Menulis file sementara
	var inputFilePaths []string
	for i, buf := range inputBuffers {
		inputFilePath := fmt.Sprintf("input%d.pdf", i)
		err := os.WriteFile(inputFilePath, buf.Bytes(), 0644)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to write temporary PDF file")
		}
		inputFilePaths = append(inputFilePaths, inputFilePath)
		defer os.Remove(inputFilePath) // Menghapus file setelah digunakan
	}

	// Menggabungkan PDF
	outputBuffer := new(bytes.Buffer)
	err := api.Merge(inputFilePaths[0], inputFilePaths[1:], outputBuffer, nil, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to merge PDF files")
	}

	// Mengirimkan hasil PDF sebagai respons
	c.Response().Header.Set("Content-Disposition", "attachment; filename=merged.pdf")
	c.Response().Header.Set("Content-Type", "application/pdf")
	return c.SendStream(outputBuffer)
}
func SendFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	filePath := fmt.Sprintf("./uploads/%s", filename)

	return c.SendFile(filePath)
}
func UploadSinglePDF(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to get file")
	}

	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	err = c.SaveFile(file, filePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file " + err.Error())
	}

	return c.JSON(fiber.Map{"filePath": filePath})
}
