package models

type State struct {
	ID              uint    `gorm:"primaryKey"`
	Title           string  `gorm:"not null"`
	BackGroundColor string  `gorm:"not null"`
	AdminAccess     bool    `gorm:"not null"`
	ProjectID       int     `gorm:"not null"`
	Project         Project `gorm:"foreignKey:ProjectID"`
	Tasks           []Task  `gorm:"foreignKey:StateID"`
}
