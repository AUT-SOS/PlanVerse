package helpers

import "time"

type ProjectMember struct {
	Email string
}

type ProjectID struct {
	ProjectID int
}

type ShowRole struct {
	IsAdmin bool
}

type Owner struct {
	OwnerID int
}

type Admin struct {
	UserID        int
	PromotionTime time.Time
}

type TaskID struct {
	ID int
}

type UserID struct {
	UserID int
}
