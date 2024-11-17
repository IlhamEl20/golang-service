package controllers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	model "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func AddAnnotation(c *fiber.Ctx) error {
	// Retrieve the file from the form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("File upload failed")
	}

	// Save the uploaded file temporarily
	tempFilePath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveFile(file, tempFilePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
	}

	defer os.Remove(tempFilePath) // Clean up

	// Open the PDF file
	f, err := os.Open(tempFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to open file")
	}
	defer f.Close()

	// Create an output file to save the result
	outFile := filepath.Join(os.TempDir(), "annotated_"+file.Filename)
	out, err := os.Create(outFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create output file")
	}
	defer out.Close()

	// Example annotation (change it based on your needs)
	// Create annotation
	annotation := &model.Annotation{
		SubType:  model.AnnText,
		Contents: "This is an annotation",
	}
	// Add annotation to the PDF
	err = api.AddAnnotations(f, out, []string{"1"}, annotation, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to add annotation")
	}

	// Return the annotated PDF file
	annotatedFile, err := ioutil.ReadFile(outFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read output file")
	}

	return c.Send(annotatedFile)
}

func RemoveAnnotationFromFile(c *fiber.Ctx) error {
	// Parse form data
	inFile := c.FormValue("inFile")
	outFile := c.FormValue("outFile")

	// Example: Page(s) from which to remove the annotation
	selectedPages := []string{"1"}

	// Example: Annotation ID to remove
	idsAndTypes := []string{"1"}

	// Example: Object number for the annotation (could be empty if not applicable)
	objNrs := []int{}

	// Create PDF configuration
	conf := model.NewDefaultConfiguration()

	// Call the pdfcpu API to remove annotations
	err := api.RemoveAnnotationsFile(inFile, outFile, selectedPages, idsAndTypes, objNrs, conf, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Annotation removed successfully",
	})
}
