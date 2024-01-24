package models

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string    `gorm:"not null"`
	Password       string    `gorm:"not null"`
	Email          string    `gorm:"unique;not null"`
	ProfilePic     string    `gorm:"not null"`
	IsVerified     bool      `gorm:"not null"`
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

type Response struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}
