package model

import "time"

type Recipe struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      uint      `json:"user_id" gorm:"not null"` // 👈 add this
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Ingredients string    `json:"ingredients" gorm:"not null"` // comma separated for now
	Steps       string    `json:"steps" gorm:"not null"`
	CookTime    int       `json:"cook_time"` // in minutes
	Servings    int       `json:"servings"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
