package models

type Post struct {
	ID      *int   `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  int    `json:"user_id" binding:"required"`
}

type PostUpdate struct {
	ID      *int   `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
