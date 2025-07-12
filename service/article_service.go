package service

import (
	"context"
	"kumparan_project/model"
	"kumparan_project/repository"
)

type ArticleServiceInterface interface {
	GetArticles(context.Context, model.ArticleQueryParams) ([]model.Article, int, error)
	CreateArticles(context.Context, model.ArticleCreateParams) (*model.Article, error)
}

type articleService struct {
	articleRepo repository.ArticleRepositoryInterface
}

func NewArticleService(articleRepo repository.ArticleRepositoryInterface) ArticleServiceInterface {
	return &articleService{articleRepo}
}

func (a *articleService) GetArticles(ctx context.Context, params model.ArticleQueryParams) ([]model.Article, int, error) {
	return a.articleRepo.Search(ctx, params)
}

func (a *articleService) CreateArticles(ctx context.Context, body model.ArticleCreateParams) (*model.Article, error) {
	article := model.Article{
		AuthorID: body.AuthorID,
		Title:    body.Title,
		Body:     body.Body,
	}

	return a.articleRepo.Create(ctx, article)
}
