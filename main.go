package main

import (
	"log"
	"os"

	"mongoapi/config"
	"mongoapi/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	})

	api := app.Group("/api")

	routes.ProductsRoute(api.Group("/products"))
}

func main() {
	path, _ := os.Getwd()
	fileName := filepath.Join(path, ".env")
	fmt.Println(fileName)
	err := godotenv.Load(filepath.Join(path, ".env"))
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	config.ConnectDB()

	setupRoutes(app)

	port := os.Getenv("PORT")
	erro := app.Listen(":" + port)

	if erro != nil {
		log.Fatal("Error app failed to start")
		log.Fatal(erro)
	}
}
