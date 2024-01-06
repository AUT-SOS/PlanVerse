package models

import "time"

type User struct {
	ID             uint      `gorm:"primaryKey"`
	Username       string    `gorm:"unique;not null"`
	Password       string    `gorm:"not null"`
	Email          string    `gorm:"unique;not null"`
	ProfilePic     string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
	DeletedAt      time.Time `gorm:"not null"`
	MemberProjects []Project `gorm:"many2many:projects_members"`
	AdminProjects  []Project `gorm:"many2many:projects_admins"`
}
