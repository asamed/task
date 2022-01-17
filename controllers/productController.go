package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"mongoapi/config"
	"mongoapi/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllProducts(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var products []model.Product

	filter := bson.M{}
	findOptions := options.Find()

	total, _ := productCollection.CountDocuments(ctx, filter)

	cursor, err := productCollection.Find(ctx, filter, findOptions)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Products Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var prod model.Product
		cursor.Decode(&prod)
		products = append(products, prod)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  products,
		"total": total,
	})
}

func GetProduct(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var prod model.Product
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		log.Fatal(err)
	}
	findResult := productCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&prod)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    prod,
		"success": true,
	})
}

func AddProduct(c *fiber.Ctx) error {
	fmt.Println(c)
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	prod := new(model.Product)
	if err := c.BodyParser(prod); err != nil {
        log.Println(err)
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Failed to parse body",
            "error":   err,
        })
    }
	fmt.Println(prod)
	result, err := productCollection.InsertOne(ctx, prod)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Product failed to insert",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Product inserted successfully",
	})

}

func UpdateProduct(c *fiber.Ctx) error {
	fmt.Println(c)
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	prod := new(model.Product)
	prod.ProductName = c.Query("prodName")
	str := c.FormValue("prodPrice")
	var err error
	prod.ProductPrice, err = strconv.Atoi(str)
	if err != nil {
		fmt.Println(err)
	}
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err,
		})
	}

	update := bson.M{
		"$set": prod,
	}
	_, err = productCollection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Product failed to update",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Product updated successfully",
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	productCollection := config.MI.DB.Collection("products")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err,
		})
	}
	_, err = productCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Product failed to delete",
			"error":   err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Product deleted successfully",
	})
}
