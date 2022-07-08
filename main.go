package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abe27/api/v2/database"
	"github.com/abe27/api/v2/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Hello(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Hello, World",
		"data":    nil,
	})
}

func Welcome(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Welcome to the API server",
		"data":    nil,
	})
}

func setUpRouter(app *fiber.App) {
	// Welcome endpoint
	app.Get("/", Hello)
	route := app.Group("api/v2")
	route.Get("/", Welcome)
	// Test Todos
	route.Get("/todos", models.GetTodos)
	route.Post("/todos", models.CreateTodo)
	route.Get("/todos/:id", models.GetTodoById)
	route.Put("/todos/:id", models.UpdateTodo)
	route.Delete("/todos/:id", models.DeleteTodo)

	// User Interface
	route.Post("/register", models.Register)
	route.Post("/login", models.Login)
	route.Delete("/logout", models.Logout)
	route.Get("/Profile", models.Profile)
}

// Initialize connect DB
func init() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	// dns := "host=" + os.Getenv("DBHOST") + " user=" + os.Getenv("DBUSER") + " password=" + os.Getenv("DBPASSWD") + " dbname=" + os.Getenv("DBNAME") + " port=" + os.Getenv("DBPORT") + ""
	dns := os.Getenv("DBUSER") + ":@tcp(" + os.Getenv("DBHOST") + ":" + os.Getenv("DBPORT") + ")/" + os.Getenv("DBNAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	database.DBConn, err = gorm.Open(mysql.Open(dns), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tbt_", // table name prefix, table for `User` would be `t_users`
			SingularTable: false,  // use singular table name, table for `User` would be `user` with this option enabled
			NoLowerCase:   false,  // skip the snake_casing of names
			NameReplacer:  strings.NewReplacer("CID", "Cid"),
		},
	})

	if err != nil {
		panic("Failed to connect to database")
	}
	fmt.Println("Database connected!")
	database.DBConn.AutoMigrate(&models.Todo{})
	database.DBConn.AutoMigrate(&models.Whs{})
	database.DBConn.AutoMigrate(&models.User{})
	fmt.Println("Migrated DB")
}

func main() {
	config := fiber.Config{
		AppName:      "API Service",
		ServerHeader: "Taweechai Yuenyang API Server",
	}
	// Initialize Fiber Framework
	app := fiber.New(config)
	// Initialize set up router
	app.Use(logger.New(logger.Config{
		Format:     "${pid} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2022",
		TimeZone:   "Asia/Bangkok",
	}))
	// Initialize RequestID
	app.Use(requestid.New())
	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	setUpRouter(app)
	log.Fatal(app.Listen(":3000"))
}
