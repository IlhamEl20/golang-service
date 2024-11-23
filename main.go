// package main

// import (
// 	"log"
// 	"provinsi/routes"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// )

// // @title Fiber CAPTCHA API
// // @version 1.0
// // @description This is a simple API to generate and verify CAPTCHA using GoFiber and dchest/captcha.
// // @contact.name Muhamad Ilham
// // @contact.email example@example.com
// // @host localhost:5000
// // @BasePath /
// // @schemes http
// func main() {
// 	// database.ConnectDB()
// 	// err := database.DB.AutoMigrate(
// 	// // 	&models.Product{})
// 	// if err != nil {
// 	// 	log.Fatalf("Error during migration: %v", err)
// 	// } else {
// 	// 	fmt.Println("Migration successfully completed.")
// 	// }

// 	app := fiber.New()
// 	app.Use(cors.New(cors.Config{
// 		AllowOrigins: "*", // Ganti dengan domain frontend yang diizinkan
// 		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
// 		AllowHeaders: "Origin, Content-Type, Accept",
// 		// AllowCredentials: true, // If you need to send cookies or other credentials
// 	}))
// 	routes.Setup(app)
// 	// app.Get("/swagger/*", swagger.HandlerDefault) // default swagger
// 	log.Fatal(app.Listen(":5000"))
// }
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/IlhamEl20/golang-service/routes"
)

// Handler - fungsi yang diekspor untuk Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	routes.Setup(app)

	if err := app.Handler()(w, r); err != nil {
		log.Printf("Error handling request: %v", err)
	}
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	routes.Setup(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

