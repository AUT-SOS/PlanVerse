package models

import "time"

type ProjectsMembers struct {
	UserID        uint `gorm:"primaryKey"`
	ProjectID     uint `gorm:"primaryKey"`
	IsAdmin       bool
	PromotionTime time.Time
}
