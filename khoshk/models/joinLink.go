package models

import "gorm.io/gorm"

type JoinLink struct {
	gorm.Model
	Link      string `gorm:"unique;not null"`
	ProjectID int    `gorm:"notnull"`
}
