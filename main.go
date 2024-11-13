package main

import (
	"log"
	"provinsi/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title Fiber CAPTCHA API
// @version 1.0
// @description This is a simple API to generate and verify CAPTCHA using GoFiber and dchest/captcha.
// @contact.name Muhamad Ilham
// @contact.email example@example.com
// @host localhost:5000
// @BasePath /
// @schemes http
func main() {
	// database.ConnectDB()
	// err := database.DB.AutoMigrate(

	// // 	&models.Product{})
	// if err != nil {
	// 	log.Fatalf("Error during migration: %v", err)
	// } else {
	// 	fmt.Println("Migration successfully completed.")
	// }

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Ganti dengan domain yang diizinkan
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	routes.Setup(app)
	// app.Get("/swagger/*", swagger.HandlerDefault) // default swagger
	log.Fatal(app.Listen(":5000"))
}
