package models

import "github.com/danilevy1212/UserPostApi-Challenge/internal/database/repositories/postgresql/ent"

type Post struct {
	ID      *int   `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  uint64 `json:"user_id" binding:"required"`
}

func (p Post) ToEnt() *ent.Post {
	dbPost := &ent.Post{
		Title:   p.Title,
		Content: p.Content,
	}

	if p.ID != nil {
		dbPost.ID = *p.ID
	}

	return dbPost
}
