package repository

import (
	"context"
	"errors"
	"kumparan_project/model"
	"kumparan_project/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	limit  = 10
	offset = 0
)

func TestArticleRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewArticleRepository(db)

	t.Run("success", func(t *testing.T) {
		query := "test query"
		author := "test author"
		params := model.ArticleQueryParams{
			Query:  &query,
			Author: &author,
			Limit:  &limit,
			Offset: &offset,
		}

		mock.ExpectQuery("SELECT COUNT\\(\\*\\)").
			WithArgs(utils.NullString(params.Query), utils.NullString(params.Author)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT a\\.id, a\\.author_id, a\\.title, a\\.body, a\\.created_at, au\\.id, au\\.name").
			WithArgs(utils.NullString(params.Query), utils.NullString(params.Author), params.Limit, params.Offset).
			WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "title", "body", "created_at", "author_id", "author_name"}).
				AddRow(uuid.New(), uuid.New(), "Test Title", "Test Body", time.Now(), uuid.New(), "Test Author"))

		articles, total, err := repo.Search(context.Background(), params)
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Len(t, articles, 1)
		assert.Equal(t, "Test Title", articles[0].Title)
		assert.Equal(t, "Test Author", articles[0].Author.Name)
	})

	t.Run("count query error", func(t *testing.T) {
		query := "test query"
		author := "test author"
		params := model.ArticleQueryParams{
			Query:  &query,
			Author: &author,
			Limit:  &limit,
			Offset: &offset,
		}

		mock.ExpectQuery("SELECT COUNT\\(\\*\\)").
			WithArgs(utils.NullString(params.Query), utils.NullString(params.Author)).
			WillReturnError(errors.New("count query error"))

		articles, total, err := repo.Search(context.Background(), params)
		assert.Error(t, err)
		assert.Equal(t, "count articles: count query error", err.Error())
		assert.Nil(t, articles)
		assert.Equal(t, 0, total)
	})

	t.Run("data query error", func(t *testing.T) {
		query := "test query"
		author := "test author"
		params := model.ArticleQueryParams{
			Query:  &query,
			Author: &author,
			Limit:  &limit,
			Offset: &offset,
		}

		mock.ExpectQuery("SELECT COUNT\\(\\*\\)").
			WithArgs(utils.NullString(params.Query), utils.NullString(params.Author)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT a\\.id, a\\.author_id, a\\.title, a\\.body, a\\.created_at, au\\.id, au\\.name").
			WithArgs(utils.NullString(params.Query), utils.NullString(params.Author), params.Limit, params.Offset).
			WillReturnError(errors.New("data query error"))

		articles, total, err := repo.Search(context.Background(), params)
		assert.Error(t, err)
		assert.Equal(t, "query articles: data query error", err.Error())
		assert.Nil(t, articles)
		assert.Equal(t, 0, total)
	})

	t.Run("scan error", func(t *testing.T) {
		query := "test query"
		author := "test author"
		params := model.ArticleQueryParams{
			Query:  &query,
			Author: &author,
			Limit:  &limit,
			Offset: &offset,
		}

		mock.ExpectQuery("SELECT COUNT\\(\\*\\)").
			WithArgs(utils.NullString(params.Query), utils.NullString(params.Author)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT a\\.id, a\\.author_id, a\\.title, a\\.body, a\\.created_at, au\\.id, au\\.name").
			WithArgs(utils.NullString(params.Query), utils.NullString(params.Author), params.Limit, params.Offset).
			WillReturnRows(sqlmock.NewRows([]string{"id", "author_id", "title", "body", "created_at", "author_id", "author_name"}).
				AddRow("1", uuid.New(), "Test Title", "Test Body", time.Now(), uuid.New(), "Test Author"))

		articles, total, err := repo.Search(context.Background(), params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "scan article")
		assert.Nil(t, articles)
		assert.Equal(t, 0, total)
	})
}

func TestArticleRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewArticleRepository(db)

	t.Run("success", func(t *testing.T) {
		article := model.Article{
			AuthorID: uuid.New(),
			Title:    "Test Title",
			Body:     "Test Body",
		}

		mock.ExpectExec("INSERT INTO articles \\(id, author_id, title, body, created_at\\)").
			WithArgs(sqlmock.AnyArg(), article.AuthorID, article.Title, article.Body, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Create(context.Background(), article)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, article.Title, result.Title)
		assert.Equal(t, article.Body, result.Body)
		assert.Equal(t, article.AuthorID, result.AuthorID)
		assert.NotZero(t, result.ID)
		assert.NotZero(t, result.CreatedAt)
	})

	t.Run("exec error", func(t *testing.T) {
		article := model.Article{
			AuthorID: uuid.New(),
			Title:    "Test Title",
			Body:     "Test Body",
		}

		mock.ExpectExec("INSERT INTO articles \\(id, author_id, title, body, created_at\\)").
			WithArgs(sqlmock.AnyArg(), article.AuthorID, article.Title, article.Body, sqlmock.AnyArg()).
			WillReturnError(errors.New("exec error"))

		result, err := repo.Create(context.Background(), article)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "fail create article: exec error", err.Error())
	})
}
