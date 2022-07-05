package models

import (
	"github.com/abe27/api/v2/database"
	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func GetTodos(c *fiber.Ctx) error {
	db := database.DBConn
	var todos []Todo
	db.Find(&todos)
	return c.JSON(&todos)
}
