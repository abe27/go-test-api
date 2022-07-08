package models

import (
	"time"

	"github.com/abe27/api/v2/database"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Whs struct {
	ID          string    `gorm:"primarykey;size:21"    json:"id"`
	Title       string    `gorm:"size:50"     json:"title"`
	Description string    `gorm:"size:255"    json:"description"`
	IsActive    bool      `json:"is_active"     default:"false"`
	CreatedAt   time.Time `json:"created_at"    default:"now"`
	UpdatedAt   time.Time `json:"updated_at"    default:"now"`
}

var r Response

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
		r.Status = false
		r.Message = "กรุณาตรวจสอบก่อนบันทึกข้อมูลด้วย"
		r.Data = err
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}

	// Generate nanoid
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	whs.ID = id
	err = db.Create(&whs).Error
	if err != nil {
		r.Status = false
		r.Message = "เกิดข้อผิดพลาดรหว่างทำการบันทึกข้อมูล\nกรุณาติดต่อผู้ดูแลระบบเพื่อทำการแก้ไขข้อผิดพลาดนี้!"
		r.Data = err
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}

	r.Status = true
	r.Message = "บันทึกข้อมูลเรียบร้อยแล้ว"
	r.Data = &whs
	return c.Status(fiber.StatusInternalServerError).JSON(r)
}

func ShowWhsById(c *fiber.Ctx) error {
	// Get Params from url
	id := c.Params("id")
	db := database.DBConn
	var whs Whs
	err := db.First(&whs, id).Error
	if err != nil {
		r.Status = true
		r.Message = "ไม่พบข้อมูล"
		r.Data = err
		return c.Status(fiber.StatusNotFound).JSON(r)
	}

	r.Status = true
	r.Message = "บันทึกข้อมูลเรียบร้อยแล้ว"
	r.Data = &whs
	return c.Status(fiber.StatusOK).JSON(r)
}

func UpdateWhs(c *fiber.Ctx) error {
	type UpdateWhs struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsActive    bool   `json:"is_active"`
	}

	id := c.Params("id")
	db := database.DBConn

	var whs Whs
	err := db.First(&whs, id).Error
	if err != nil {
		r.Status = false
		r.Message = "ไม่พบข้อมูล"
		r.Data = err
		return c.Status(fiber.StatusNotFound).JSON(r)
	}

	var updateWhs UpdateWhs
	err = c.BodyParser(&updateWhs)
	if err != nil {
		r.Status = false
		r.Message = "กรุณาตรวจสอบก่อนบันทึกข้อมูลด้วย"
		r.Data = err
		return c.Status(fiber.StatusNotAcceptable).JSON(r)
	}

	// Update Data
	whs.Title = updateWhs.Title
	whs.Description = updateWhs.Description
	whs.IsActive = updateWhs.IsActive
	whs.UpdatedAt = time.Now()
	db.Save(&whs)

	r.Status = true
	r.Message = "อัพเดทข้อมูลเรียบร้อยแล้ว"
	r.Data = &whs
	return c.Status(fiber.StatusAccepted).JSON(r)
}

func DeleteWhs(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DBConn
	var whs Whs
	err := db.First(&whs, id).Error
	if err != nil {
		r.Status = false
		r.Message = "ไม่พบข้อมูล"
		r.Data = err
		return c.Status(fiber.StatusNotFound).JSON(r)
	}

	// Delete Data
	db.Delete(&whs)
	r.Status = false
	r.Message = "บันทึกข้อมูลเรียบร้อยแล้ว"
	r.Data = &whs
	return c.Status(fiber.StatusNotFound).JSON(r)
}
