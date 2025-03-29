package models

import "github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"

type User struct {
	ID    *int   `json:"id"`
	Name  string `json:"name"  binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (u User) ToEnt() *ent.User {
	dbUser := &ent.User{
		Name:  u.Name,
		Email: u.Email,
	}

	if u.ID != nil {
		dbUser.ID = *u.ID
	}

	return dbUser
}
