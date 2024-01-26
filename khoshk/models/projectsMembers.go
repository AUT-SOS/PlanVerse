package models

type ProjectsMembers struct {
	UserID    uint `gorm:"primaryKey"`
	ProjectID uint `gorm:"primaryKey"`
	IsAdmin   bool
}
