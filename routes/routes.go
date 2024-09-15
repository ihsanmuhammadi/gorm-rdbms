package routes

import (
	"gorm-rdbms/controller"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App)  {
	// Get all posts
	app.Get("/post", controller.GetAllPosts)

	// Get post by id
	app.Get("/post/:id", controller.GetPostById)

	// Create post
	app.Post("/post", controller.CreatePost)

	// Update post
	app.Put("/post/:id", controller.UpdatePost)

	// Delete post
	app.Delete("/post/:id", controller.DeletePost)

	// Post News to DB
	app.Post("/news", controller.CreateNews)
}
