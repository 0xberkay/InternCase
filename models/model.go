package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name     string `gorm:"not null" validate:"required"`
	Email    string `gorm:"not null" validate:"required,email"`
	Password string `gorm:"not null" validate:"required"`
}

type Plan struct {
	gorm.Model
	StudentID   uint      `gorm:"not null" validate:"required"`
	Description string    `gorm:"not null" validate:"required"`
	Start       time.Time `gorm:"not null" validate:"required"`
	End         time.Time `gorm:"not null" validate:"required"`
	StatuDataID uint      `gorm:"not null" validate:"required"`
}

type StatuData struct {
	gorm.Model
	StatuName string `gorm:"not null" validate:"required"`
}

type StudentLoginForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type Claims struct {
	ID uint
	jwt.RegisteredClaims
}
