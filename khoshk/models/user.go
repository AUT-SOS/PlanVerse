package models

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string    `gorm:"not null"`
	Password   string    `gorm:"not null"`
	Email      string    `gorm:"unique;not null"`
	ProfilePic string    `gorm:"not null"`
	IsVerified bool      `gorm:"not null"`
	Projects   []Project `gorm:"many2many:projects_members"`
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

type UserResponse struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

type GetUserResponse struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
}
