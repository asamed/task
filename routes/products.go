package routes

import (
	"mongoapi/controllers"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func ProductsRoute(route fiber.Router) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		route.Get("/", controllers.GetAllProducts)
		wg.Done()
	}()
	route.Get("/:id", controllers.GetProduct)
	go func() {
		route.Post("/", controllers.AddProduct)
		wg.Done()
	}()
	route.Put("/:id", controllers.UpdateProduct)
	route.Delete("/:id", controllers.DeleteProduct)
	wg.Wait()
}
