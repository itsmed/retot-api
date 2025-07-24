package handler

import (
	"fmt"
	"log"
	"strconv"

	"app/database"
	"app/model"

	"github.com/gofiber/fiber/v2"
)

// func GetPostFromUser gets all posts from user with id
func GetPostFromUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	db := database.DB
	id, err := strconv.Atoi(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var posts []model.Post
	// Fetch posts for the user with the given ID
	if err := db.Where("user_id = ?", id).Find(&posts).Error; err != nil {
		log.Printf("Error fetching posts for user ID %d: %v", id, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error fetching posts",
			"data":    nil,
		})
	}

	// Check if posts were found
	if len(posts) == 0 {
		log.Printf("No posts found for user with ID %d", id)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No posts found for user with ID " + userId,
			"data":    nil,
		})
	}

	// Return the posts for the user
	return c.JSON(posts)
}

// CreatePost creates a new post
func CreatePost(c *fiber.Ctx) error {
	fmt.Println("create post called", c.Locals("user_id"))
	var post model.Post

	// Parse the request body into the Post struct
	if err := c.BodyParser(&post); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"data":    nil,
		})
	}

	db := database.DB

	// Save the post to the database
	if err := db.Create(&post).Error; err != nil {
		log.Printf("Error creating post: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Error creating post",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Post created successfully",
		"data":    post,
	})
}
