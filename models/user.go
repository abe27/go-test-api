package models

import (
	"time"

	"github.com/abe27/api/v2/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
	"github.com/golang-jwt/jwt/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string    `gorm:"primarykey;size:21" json:"id"`
	UserName  string    `gorm:"unique;not null;;size:10" json:"username"`
	Email     string    `gorm:"default:null;size:25" json:"email"`
	Password  string    `gorm:"not null;size:255" json:"password"`
	IsVerify  bool      `json:"is_verify" default:"false"`
	CreatedAt time.Time `json:"created_at" default:"now"`
	UpdatedAt time.Time `json:"updated_at" default:"now"`
}

type Auth struct {
	AuthType string `json:"auth_type"`
	Token    string `json:"token"`
	UserName string `json:"user_name"`
}

const (
	jwtSecret = "ADSads123"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthSuccess(c *fiber.Ctx) error {
	c.Next()
	return nil
}

func CreateToken(name string) string {
	/// Create Token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = utils.UUID()
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		panic(err)
	}

	return t
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
	password := user.Password
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity
	// match := CheckPasswordHash(password, hash)
	// if match {
	// 	user.Password = hash
	// }
	user.Password = hash

	var auth Auth
	auth.AuthType = "Bearer"
	auth.Token = CreateToken(user.UserName)
	auth.UserName = user.UserName

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
		"data":    &auth,
	})
}

func Login(c *fiber.Ctx) error {
	db := database.DBConn
	var login User
	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "กรุณาตรวจสอบความถูกต้องของข้อมูลด้วย",
			"data":    nil,
		})
	}

	hand_check_passwd := login.Password
	err = db.Where("user_name=?", login.UserName).First(&login).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "ไม่พบข้อมูลผู้ใช้งาน",
			"data":    nil,
		})
	}

	password := hand_check_passwd
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity
	match := CheckPasswordHash(password, hash)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  match,
		"message": hash,
		"data":    &login,
	})
}
