package models

import "time"

type GediType struct {
	ID          string    `gorm:"primarykey;size:21"    json:"id"`
	Title       string    `gorm:"size:50"               json:"title"`
	Description string    `gorm:"size:255"              json:"description"`
	IsActive    bool      `json:"is_active"             default:"false"`
	CreatedAt   time.Time `json:"created_at"            default:"now"`
	UpdatedAt   time.Time `json:"updated_at"            default:"now"`
}
