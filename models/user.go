package models

import (
	"time"

	"github.com/abe27/api/v2/database"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type User struct {
	ID        string    `gorm:"primarykey;length:21" json:"id"`
	UserName  string    `gorm:"uniqueIndex;type:length:10" json:"username"`
	Email     string    `gorm:"default:null;length:25" json:"email"`
	Password  string    `gorm:"length:25" json:"password"`
	IsVerify  bool      `json:"is_verify" default:"false"`
	CreatedAt time.Time `json:"created_at" default:"now"`
	UpdatedAt time.Time `json:"updated_at" default:"now"`
}

func Register(c *fiber.Ctx) error {
	db := database.DBConn
	user := new(User)
	err := c.BodyParser(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "กรุณาตรวจสอบก่อนบันทึกข้อมูลด้วย",
			"data":    err,
		})
	}

	// Generate Nanoid
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}

	user.ID = id
	err = db.Create(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "ไม่สามารถบันทึกข้อมูลได้\nกรุณาติดต่อผู้ดูแลระบบเพื่อทำการแก้ไขข้อผิดพลาดนี้",
			"data":    err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  true,
		"message": "บันทึกข้อมูลเรียบร้อยแล้ว",
		"data":    &user,
	})
}
