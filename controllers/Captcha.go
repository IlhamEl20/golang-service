package controllers

import (
	"bytes"
	"encoding/base64"
	"provinsi/utils"

	"github.com/dchest/captcha"
	"github.com/gofiber/fiber/v2"
)

// GetCaptchaHandler handles the GET request to generate and serve a new captcha
// @Summary Generate a new captcha
// @Description Generate a new captcha and return the captcha ID and image
// @Tags Captcha
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /get-captcha [get]
func GetCaptchaHandler(c *fiber.Ctx) error {

	captchaID := captcha.New()
	// Set custom header for captcha ID
	c.Set("X-Captcha-ID", captchaID)

	// Expose custom headers to client-side JS
	c.Set("Access-Control-Expose-Headers", "X-Captcha-ID")

	// Log the captcha ID

	// Generate the captcha image as a base64 string
	var captchaImageBuffer bytes.Buffer
	err := captcha.WriteImage(&captchaImageBuffer, captchaID, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		return err
	}
	captchaImageBase64 := base64.StdEncoding.EncodeToString(captchaImageBuffer.Bytes())

	// Write the captcha image base64 string to the response
	return c.JSON(fiber.Map{
		"image": captchaImageBase64,
	})
}

// VerifyCaptchaHandler handles the POST request to verify a captcha
// @Summary Verify a captcha
// @Description Verify the captcha solution for the given captcha ID
// @Tags Captcha
// @Accept json
// @Produce json
// @Param captcha_id body string true "Captcha ID"
// @Param solution body string true "Captcha Solution"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /verify-captcha [post]
func VerifyCaptchaHandler(c *fiber.Ctx) error {
	// Parse JSON body
	var request struct {
		CaptchaID string `json:"captcha_id"`
		Solution  string `json:"solution"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Verify the captcha
	if utils.VerifyCaptcha(request.CaptchaID, request.Solution) {
		return c.JSON(fiber.Map{"message": "Captcha verified successfully", "code": 200, "status": "ok"})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid captcha solution or expired captcha"})
	}
}
