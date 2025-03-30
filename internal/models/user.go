package models

type User struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UserUpdate struct {
	ID    *uint64 `json:"id"`
	Name  string  `json:"name"  binding:"required"`
	Email string  `json:"email" binding:"required,email"`
}
