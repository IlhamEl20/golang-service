package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// WatermarkRequest represents the request payload for watermarking.
type WatermarkRequest struct {
	Text   string `json:"text" validate:"required"`
	OnTop  bool   `json:"onTop"`
	Update bool   `json:"update"`
	Pages  string `json:"pages"`
	Desc   string `json:"desc"`
	Image  string `json:"image"`
}

// CreateWatermark handles watermark creation for images or PDFs.
func CreateWatermark(c *fiber.Ctx) error {
	req := new(WatermarkRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file upload is required"})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer src.Close()

	// Detect file type from MIME
	mimeType := file.Header.Get("Content-Type")
	fmt.Println("Detected MIME type:", mimeType)

	// Convert file to byte buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, src); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	var outputBuf bytes.Buffer
	switch {
	case mimeType == "application/pdf":
		err = applyPDFWatermark(&buf, &outputBuf, req.Text, req.OnTop, req.Update, req.Pages, req.Desc, req.Image)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "unsupported file type"})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// Ensure the watermark directory exists
	watermarkDir := "watermark"
	if err := os.MkdirAll(watermarkDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Save the watermarked file in the watermark directory
	outputFile := watermarkDir + "/output_" + file.Filename
	if err := os.WriteFile(outputFile, outputBuf.Bytes(), 0644); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// // Send the watermarked file directly
	// c.Set("Content-Type", "application/pdf")
	// c.Set("Content-Disposition", "attachment; filename=watermarked_" + file.Filename)
	// return c.SendStream(bytes.NewReader(outputBuf.Bytes()))
	return c.Download(outputFile)
}

func applyPDFWatermark(inputBuf, outputBuf *bytes.Buffer, text string, onTop, update bool, pages, desc, image string) error {
	if text == "" && image == "" {
		text = "Confidential"
	}

	// Set default description if not provided
	if desc == "" {
		desc = "points:24, scale:0.8, color:.8 .8 .4, op:.6, rotation:45"
	}
	fmt.Println("PDF Watermark Configuration:", desc)

	// Create a temporary input file
	inputFile := "temp_input.pdf"
	if err := os.WriteFile(inputFile, inputBuf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to create temporary input file: %w", err)
	}
	defer os.Remove(inputFile)

	// Create a temporary output file
	outputFile := "temp_output.pdf"
	defer os.Remove(outputFile)

	// Apply the watermark to the specified pages
	selectedPages := []string{pages}
	if pages == "" {
		selectedPages = []string{"1-"} // Apply to all pages if no specific pages are provided
	}

	var err error
	if image != "" {
		// Decode base64 image
		imageData, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			return fmt.Errorf("failed to decode base64 image: %w", err)
		}

		// Save the image to a temporary file
		imageFile := "temp_image_watermark.png"
		if err := os.WriteFile(imageFile, imageData, 0644); err != nil {
			return fmt.Errorf("failed to create temporary image file: %w", err)
		}
		defer os.Remove(imageFile)
		if update {
			// Update image watermark
			if err = api.UpdateImageWatermarksFile(inputFile, outputFile, selectedPages, onTop, imageFile, desc, model.NewDefaultConfiguration()); err != nil {
				return fmt.Errorf("failed to update image watermark: %w", err)
			}
		} else {
			if err = api.AddImageWatermarksFile(inputFile, outputFile, selectedPages, onTop, imageFile, desc, model.NewDefaultConfiguration()); err != nil {
				return fmt.Errorf("failed to apply image watermark: %w", err)
			}
		}

	} else {
		if update {
			// Update text watermark
			err = api.UpdateTextWatermarksFile(inputFile, outputFile, selectedPages, onTop, text, desc, model.NewDefaultConfiguration())
		} else {
			// Apply text watermark
			err = api.AddTextWatermarksFile(inputFile, outputFile, selectedPages, onTop, text, desc, model.NewDefaultConfiguration())
		}
	}

	if err != nil {
		return fmt.Errorf("failed to apply watermark: %w", err)
	}

	// Read the output file into the output buffer
	outputData, err := os.ReadFile(outputFile)
	if err != nil {
		return fmt.Errorf("failed to read temporary output file: %w", err)
	}
	outputBuf.Write(outputData)

	return nil
}
