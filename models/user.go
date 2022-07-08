package models

import (
	"time"

	"github.com/abe27/api/v2/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/utils"
	"github.com/golang-jwt/jwt/v4"
	gnid "github.com/matoous/go-nanoid/v2"
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
	jwtSecret = "b86cb8ef1c43e44a32acd56e60666ec533dbe58eaf3bbe4dde461f522568309b"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	// err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
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
	var r Response
	user := new(User)
	err := c.BodyParser(user)

	if err != nil {
		r.Status = false
		r.Message = "กรุณาตรวจสอบก่อนบันทึกข้อมูลด้วย"
		r.Data = err
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}

	// Generate Nanoid
	id, err := gnid.New()
	if err != nil {
		panic(err)
	}

	user.ID = id
	password := user.Password
	hash, _ := HashPassword(password) // ignore error for the sake of simplicity
	user.Password = hash

	var auth Auth
	auth.AuthType = "Bearer"
	auth.Token = CreateToken(user.UserName)
	auth.UserName = user.UserName

	err = db.Create(&user).Error
	if err != nil {
		r.Status = false
		r.Message = "ไม่สามารถบันทึกข้อมูลได้\nกรุณาติดต่อผู้ดูแลระบบเพื่อทำการแก้ไขข้อผิดพลาดนี้"
		r.Data = err
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}

	r.Status = true
	r.Message = "บันทึกข้อมูลเรียบร้อยแล้ว"
	r.Data = &auth
	return c.Status(fiber.StatusInternalServerError).JSON(r)
}

func Login(c *fiber.Ctx) error {
	var r Response
	db := database.DBConn
	var login User
	err := c.BodyParser(&login)
	if err != nil {
		r.Status = false
		r.Message = "กรุณาตรวจสอบความถูกต้องของข้อมูลด้วย"
		r.Data = err
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}

	hand_check_passwd := login.Password
	err = db.Where("user_name=?", login.UserName).First(&login).Error
	if err != nil {
		r.Status = false
		r.Message = "ไม่พบข้อมูลผู้ใช้งาน"
		r.Data = err
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}

	match := CheckPasswordHash(hand_check_passwd, login.Password)
	if !match {
		r.Status = false
		r.Message = "กรุณาระบบรหัสผ่านให้ถูกต้องด้วย"
		r.Data = nil
		return c.Status(fiber.StatusInternalServerError).JSON(r)
	}

	var auth Auth
	auth.AuthType = "Bearer"
	auth.Token = CreateToken(login.UserName)
	auth.UserName = login.ID

	// Create cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    auth.Token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	r.Status = false
	r.Message = "ยินดีต้อนรับเข้าสู่ระบบ API Service By Golang"
	r.Data = &auth
	return c.Status(fiber.StatusInternalServerError).JSON(r)
}

func Logout(c *fiber.Ctx) error {
	var r Response
	// Remove cookie
	// -time.Hour = expires before one hour
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	r.Status = false
	r.Message = "ออกจากระบบ API Service By Golang 😘 เรียบร้อยแล้ว"
	r.Data = nil
	return c.Status(fiber.StatusInternalServerError).JSON(r)
}
