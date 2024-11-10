package routes

import (
	"provinsi/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Define route to get all provinces
	app.Get("/api/provinsi", controllers.GetAllProvinsi)

	// Define route to get cities by province
	app.Get("/api/provinsi/:provinsiID/kota", controllers.GetKotasByProvinsi)

	// Define route to get kecamatan by city
	app.Get("/api/kota/:kotaID/kecamatan", controllers.GetKecamatansByKota)

	// Define route to get kelurahan by kecamatan
	app.Get("/api/kecamatan/:kecamatanID/kelurahan", controllers.GetKelurahansByKecamatan)

	app.Post("/check-pdf", controllers.UploadPDF)
	app.Get("/get-captcha", controllers.GetCaptchaHandler)
	app.Post("/verify-captcha", controllers.VerifyCaptchaHandler)
}
