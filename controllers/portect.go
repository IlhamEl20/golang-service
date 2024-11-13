package controllers

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	model "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type Command struct {
	InFile   *string
	OutFile  *string
	Password *string
	Conf     *model.Configuration
}

// Encrypt inFile and write result to outFile with a password.
func Encrypt(cmd *Command) ([]string, error) {
	if cmd.Password != nil {
		cmd.Conf.UserPW = *cmd.Password
		cmd.Conf.OwnerPW = *cmd.Password // Ensure both user and owner passwords are set
	}
	return nil, api.EncryptFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
}

func UploadAndEncryptPDF(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Unable to parse form")
	}

	// Retrieve the file from form data
	fileHeader := form.File["pdf"][0]
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Unable to retrieve file")
	}
	defer file.Close()

	// Create a temporary file to store the uploaded PDF
	tempDir := os.TempDir()
	tempFile, err := os.Create(filepath.Join(tempDir, fileHeader.Filename))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Unable to create temporary file")
	}
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Unable to save file")
	}

	// Get the password from form data
	password := c.FormValue("password")
	if password == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Password is required")
	}

	// Encrypt the PDF
	inFile := tempFile.Name()
	outFile := filepath.Join(tempDir, "encrypted_"+fileHeader.Filename)
	cmd := &Command{
		InFile:   &inFile,
		OutFile:  &outFile,
		Password: &password,
		Conf:     model.NewDefaultConfiguration(),
	}
	_, err = Encrypt(cmd)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Unable to encrypt file")
	}

	// Serve the encrypted PDF with a download prompt
	c.Attachment(outFile)
	return c.SendFile(outFile)
}
