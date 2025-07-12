package model

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Author    Author    `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Article) TableName() string {
	return "articles"
}

type ArticleCreateParams struct {
	AuthorID uuid.UUID `json:"author_id" binding:"required"`
	Title    string    `json:"title" binding:"required"`
	Body     string    `json:"body" binding:"required"`
}

type ArticleQueryParams struct {
	Author *string
	Query  *string
	Limit  *int
	Offset *int
}

type ArticleCreateResponse struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
