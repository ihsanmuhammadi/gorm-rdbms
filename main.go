package main

import (
	"gorm-rdbms/database"
	"gorm-rdbms/routes"

	"github.com/gofiber/fiber/v2"
)

func main()  {
	// Connect to database & migrate
	database.Database()

	// Server
	app := fiber.New()

	// Routes
	routes.Routes(app)

	app.Listen(":9000")
}
