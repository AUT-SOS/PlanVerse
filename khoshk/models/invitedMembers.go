package models

type InvitedMembers struct {
	UserID    uint `gorm:"primaryKey"`
	ProjectID uint `gorm:"primaryKey"`
}
