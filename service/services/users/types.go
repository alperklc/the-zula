package users

import "time"

type User struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Language    string    `json:"language"`
	Theme       string    `json:"theme"`
}

type UserUpdate struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Language    string `json:"language"`
	Theme       string `json:"theme"`
}
