package routes

import (
	"mongoapi/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProductsRoute(route fiber.Router) {
	route.Get("/", controllers.GetAllProducts)
	route.Get("/:id", controllers.GetProduct)
	route.Post("/", controllers.AddProduct)
	route.Put("/:id", controllers.UpdateProduct)
}
