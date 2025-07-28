package handler

import (
	"fmt"
	"log"
	"strconv"

	"app/database"
	"app/middleware"
	"app/model"

	"github.com/gofiber/fiber/v2"
)

// GetAllItems gets all items from all categories
func GetAllItems(c *fiber.Ctx) error {
	db := database.DB
	var items []model.Item

	// Fetch all items from the database
	if err := db.Find(&items).Error; err != nil {
		log.Printf("Error fetching items: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching items",
			"data":    nil,
		})
	}

	// Check if items were found
	if len(items) == 0 {
		log.Println("No items found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No items found",
			"data":    nil,
		})
	}

	// Return the list of items
	return c.JSON(items)
}

// GetItemFromCategory gets all items from category with id
func GetItemFromCategory(c *fiber.Ctx) error {
	categoryId := c.Params("id")
	db := database.DB
	id, err := strconv.Atoi(categoryId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	var items []model.Item
	// Fetch items for the category with the given ID
	if err := db.Where("category_id = ?", id).Find(&items).Error; err != nil {
		log.Printf("Error fetching items for category ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching items",
			"data":    nil,
		})
	}

	// Check if items were found
	if len(items) == 0 {
		log.Printf("No items found for category with ID %d", id)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No items found for category with ID " + categoryId,
			"data":    nil,
		})
	}

	// Return the items for the category
	return c.JSON(items)
}

// CreateItem creates a new item
func CreateItem(c *fiber.Ctx) error {
	var item model.Item
	fmt.Println("create item called", c.Locals("user_id"))
	// Parse the request body into the Item struct
	if err := c.BodyParser(&item); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    nil,
		})
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		log.Printf("Error getting user ID: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
			"data":    nil,
		})
	}
	item.User.ID = userID

	// Save the item to the database
	db := database.DB
	if err := db.Create(&item).Error; err != nil {
		log.Printf("Error creating item: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error creating item",
			"data":    nil,
		})
	}

	return c.JSON(item)
}

// GetItemFromUser gets all items from user with id
func GetItemFromUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	db := database.DB
	id, err := strconv.Atoi(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var items []model.Item
	// Fetch items for the user with the given ID
	if err := db.Where("user_id = ?", id).Find(&items).Error; err != nil {
		log.Printf("Error fetching items for user ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching items",
			"data":    nil,
		})
	}

	// Check if items were found
	if len(items) == 0 {
		log.Printf("No items found for user with ID %d", id)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No items found for user with ID " + userId,
			"data":    nil,
		})
	}

	// Return the items for the user
	return c.JSON(items)
}

// GetItemFromId gets item with id
func GetItemFromId(c *fiber.Ctx) error {
	itemId := c.Params("id")
	db := database.DB
	id, err := strconv.Atoi(itemId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	var item model.Item
	// Fetch the item with the given ID
	if err := db.First(&item, id).Error; err != nil {
		log.Printf("Error fetching item with ID %d: %v", id, err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Item not found",
			"data":    nil,
		})
	}

	return c.JSON(item)
}

// UpdateItem updates an item with id
func UpdateItem(c *fiber.Ctx) error {
	db := database.DB
	itemId := c.Params("id")
	id, err := strconv.Atoi(itemId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var item model.Item
	if err := db.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	if item.UserID != uint(userID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not own this item"})
	}

	var input model.Item
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	item.Name = input.Name
	item.Description = input.Description
	item.Price = input.Price

	if err := db.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update item"})
	}

	return c.JSON(item)
}

// DeleteItem deletes an item with id
func DeleteItem(c *fiber.Ctx) error {
	itemId := c.Params("id")
	id, err := strconv.Atoi(itemId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	db := database.DB
	var item model.Item
	if err := db.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
	}

	if item.UserID != uint(userID) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You do not own this item"})
	}

	if err := db.Delete(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete item"})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Item deleted successfully"})
}
