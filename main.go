package main

import (
	"fmt"
	"os"

	handlers "github.com/abe27/api/v2/controllers"
	"github.com/abe27/api/v2/database"
	"github.com/abe27/api/v2/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDatabase() {
	var err error
	dsn := "host=" + os.Getenv("DBHOST") + " user=" + os.Getenv("DBUSER") + " password=" + os.Getenv("DBPASSWD") + " dbname=" + os.Getenv("DBNAME") + " port=" + os.Getenv("DBPORT") + ""
	database.DBConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connected!")
	database.DBConn.AutoMigrate(&models.Todo{})
	fmt.Println("Migrated DB")
}

func main() {
	app := fiber.New()
	initDatabase()
	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(handlers.Handler)
	})
	app.Listen(":3000")
}
