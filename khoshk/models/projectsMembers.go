package models

type ProjectMembers struct {
	UserID    int  `gorm:"primaryKey"`
	ProjectID int  `gorm:"primaryKey"`
	IsAdmin   bool `gorm:"not null"`
}
