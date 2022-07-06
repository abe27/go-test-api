package models

import (
	"time"

	"github.com/abe27/api/v2/database"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Whs struct {
	ID          string    `gorm:"primarykey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" default:"false"`
	CreatedAt   time.Time `json:"created_at" default:"now"`
	UpdatedAt   time.Time `json:"updated_at" default:"now"`
}

func GetWhs(c *fiber.Ctx) error {
	db := database.DBConn
	var whs []Whs
	db.Find(&whs)
	return c.Status(fiber.StatusOK).JSON(&whs)
}

func CreateWhs(c *fiber.Ctx) error {
	db := database.DBConn
	whs := new(Whs)
	err := c.BodyParser(whs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "กรุณาตรวจสอบก่อนบันทึกข้อมูลด้วย",
			"data":    err,
		})
	}

	// Generate nanoid
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	whs.ID = id
	err = db.Create(&whs).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "เกิดข้อผิดพลาดรหว่างทำการบันทึกข้อมูล\nกรุณาติดต่อผู้ดูแลระบบเพื่อทำการแก้ไขข้อผิดพลาดนี้!",
			"data":    err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": "บันทึกข้อมูลเรียบร้อยแล้ว",
		"data":    &whs,
	})
}

func ShowWhsById(c *fiber.Ctx) error {
	// Get Params from url
	id := c.Params("id")
	db := database.DBConn
	var whs Whs
	err := db.First(&whs, id).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "ไม่พบข้อมูล!",
			"data":    err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "บันทึกข้อมูลเรียบร้อยแล้ว",
		"data":    &whs,
	})
}
