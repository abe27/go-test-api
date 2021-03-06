package models

import (
	"time"

	"github.com/abe27/api/v2/database"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Todo struct {
	ID          string `gorm:"primarykey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func GetTodos(c *fiber.Ctx) error {
	db := database.DBConn
	var todos []Todo
	db.Find(&todos)
	return c.JSON(&todos)
}

func CreateTodo(c *fiber.Ctx) error {
	db := database.DBConn
	todo := new(Todo)
	err := c.BodyParser(todo)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "massage": "Check your input", "data": err})
	}

	// Generate Nanoid string
	id, ex := gonanoid.New()
	if ex != nil {
		panic(ex)
	}
	todo.ID = id
	err = db.Create(&todo).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "massage": "Could not create todo", "data": err})
	}
	return c.Status(fiber.StatusCreated).JSON(&todo)
}

func GetTodoById(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var todo Todo
	err := db.Find(&todo, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "massage": "Could not find todo", "data": err})
	}

	return c.Status(fiber.StatusFound).JSON(&todo)
}

func UpdateTodo(c *fiber.Ctx) error {
	type UpdateTodo struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}
	id := c.Params("id")
	db := database.DBConn
	var todo Todo
	err := db.Find(&todo, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "massage": "Could not find todo", "data": err})
	}

	var updatedTodo UpdateTodo
	err = c.BodyParser(&updatedTodo)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "massage": "Review your input", "data": err})
	}

	todo.Title = updatedTodo.Title
	todo.Description = updatedTodo.Description
	todo.Completed = updatedTodo.Completed
	db.Save(&todo)
	return c.Status(fiber.StatusAccepted).JSON(&todo)
}

func DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var todo Todo
	err := db.Find(&todo, id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "massage": "Could not find todo", "data": err})
	}
	db.Delete(&todo)
	return c.Status(200).JSON(fiber.Map{"status": "success", "massage": "Delete completed", "data": nil})
}
