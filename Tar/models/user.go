package models

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type User struct {
	ID             uint      `gorm:"primaryKey"`
	Username       string    `gorm:"unique;not null"`
	Password       string    `gorm:"not null"`
	Email          string    `gorm:"unique;not null"`
	ProfilePic     string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
	MemberProjects []Project `gorm:"many2many:projects_members"`
	AdminProjects  []Project `gorm:"many2many:projects_admins;"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type VerifyRequest struct {
	OTP string `json:"otp"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
