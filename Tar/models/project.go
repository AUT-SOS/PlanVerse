package models

type Project struct {
	ID            uint     `gorm:"primaryKey"`
	Title         string   `gorm:"not null"`
	Description   string   `gorm:"not null"`
	BackGroundPic string   `gorm:"not null"`
	OwnerID       int      `gorm:"not null"`
	Members       []User   `gorm:"many2many:projects_members"`
	Admins        []User   `gorm:"many2many:projects_admins"`
	States        []State  `gorm:"foreignKey:ProjectID"`
	JoinLink      JoinLink `gorm:"foreignKey:ProjectID"`
}
