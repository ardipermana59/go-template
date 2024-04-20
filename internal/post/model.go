package post

import (
	"time"

	"github.com/ardipermana59/go-template/internal/user"
)

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      user.User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePostDTO struct {
	Title   string `json:"title" binding:"required,min=3"`
	Content string `json:"content" binding:"required,min=10"`
}

type UpdatePostDTO struct {
	Title   string `json:"title" binding:"omitempty,min=3"`
	Content string `json:"content" binding:"omitempty,min=10"`
}

type PostResponse struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	UserID    uint               `json:"user_id"`
	User      *user.UserResponse `json:"user"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func (p *Post) ToResponse() *PostResponse {
	return &PostResponse{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		UserID:    p.UserID,
		User:      p.User.ToResponse(),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
