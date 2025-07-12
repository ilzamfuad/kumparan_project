package service_test

import (
	"context"
	"errors"
	"testing"

	"kumparan_project/model"
	"kumparan_project/service"
	mock_repository "kumparan_project/tests/mock/repository"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestArticleService_GetArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockArticleRepositoryInterface(ctrl)
	articleService := service.NewArticleService(mockRepo)

	t.Run("success", func(t *testing.T) {
		queryParams := model.ArticleQueryParams{
			Query:  nil,
			Author: nil,
			Limit:  intPtr(10),
			Offset: intPtr(0),
		}

		mockArticles := []model.Article{
			{ID: uuid.New(), Title: "Test Article 1", Body: "Test Body 1"},
			{ID: uuid.New(), Title: "Test Article 2", Body: "Test Body 2"},
		}

		mockRepo.EXPECT().Search(context.Background(), queryParams).Return(mockArticles, 2, nil)

		articles, total, err := articleService.GetArticles(context.Background(), queryParams)

		assert.NoError(t, err)
		assert.Equal(t, 2, total)
		assert.Equal(t, mockArticles, articles)
	})

	t.Run("error from repository", func(t *testing.T) {
		queryParams := model.ArticleQueryParams{
			Query:  nil,
			Author: nil,
			Limit:  intPtr(10),
			Offset: intPtr(0),
		}

		mockRepo.EXPECT().Search(context.Background(), queryParams).Return(nil, 0, errors.New("repository error"))

		articles, total, err := articleService.GetArticles(context.Background(), queryParams)

		assert.Error(t, err)
		assert.Nil(t, articles)
		assert.Equal(t, 0, total)
	})
}

func TestArticleService_CreateArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockArticleRepositoryInterface(ctrl)
	articleService := service.NewArticleService(mockRepo)

	t.Run("success", func(t *testing.T) {
		createParams := model.ArticleCreateParams{
			AuthorID: uuid.New(),
			Title:    "New Article",
			Body:     "New Body",
		}
		articleParams := model.Article{
			AuthorID: createParams.AuthorID,
			Title:    createParams.Title,
			Body:     createParams.Body,
		}

		mockArticle := &model.Article{
			ID:       uuid.New(),
			AuthorID: createParams.AuthorID,
			Title:    createParams.Title,
			Body:     createParams.Body,
		}

		mockRepo.EXPECT().Create(gomock.Any(), articleParams).Return(mockArticle, nil)

		article, err := articleService.CreateArticles(context.Background(), createParams)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle, article)
	})

	t.Run("error from repository", func(t *testing.T) {
		createParams := model.ArticleCreateParams{
			AuthorID: uuid.New(),
			Title:    "New Article",
			Body:     "New Body",
		}

		articleParams := model.Article{
			AuthorID: createParams.AuthorID,
			Title:    createParams.Title,
			Body:     createParams.Body,
		}

		mockRepo.EXPECT().Create(gomock.Any(), articleParams).Return(nil, errors.New("repository error"))

		article, err := articleService.CreateArticles(context.Background(), createParams)

		assert.Error(t, err)
		assert.Nil(t, article)
	})
}

func intPtr(i int) *int {
	return &i
}
