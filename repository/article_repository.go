package repository

import (
	"context"
	"database/sql"
	"fmt"
	"kumparan_project/model"
	"kumparan_project/utils"
	"time"

	"github.com/google/uuid"
)

type ArticleRepositoryInterface interface {
	Create(ctx context.Context, article model.Article) (*model.Article, error)
	Search(ctx context.Context, params model.ArticleQueryParams) ([]model.Article, int, error)
}

type articleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) ArticleRepositoryInterface {
	return &articleRepository{db}
}

func (r *articleRepository) Search(ctx context.Context, params model.ArticleQueryParams) ([]model.Article, int, error) {
	countQuery := `
        SELECT COUNT(*)
        FROM articles a
        JOIN authors au ON a.author_id = au.id
        WHERE
            ($1 = '' OR a.tsv @@ plainto_tsquery('english', $1)) AND
            ($2 = '' OR au.name ILIKE '%' || $2 || '%');
    `
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, utils.NullString(params.Query), utils.NullString(params.Author)).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count articles: %w", err)
	}

	dataQuery := `
        SELECT
            a.id, a.author_id, a.title, a.body, a.created_at,
            au.id, au.name
        FROM articles a
        JOIN authors au ON a.author_id = au.id
        WHERE
            ($1 = '' OR a.tsv @@ plainto_tsquery('english', $1)) AND
            ($2 = '' OR au.name ILIKE '%' || $2 || '%')
        ORDER BY a.created_at DESC
        LIMIT $3 OFFSET $4;
    `
	rows, err := r.db.QueryContext(ctx, dataQuery, utils.NullString(params.Query), utils.NullString(params.Author), params.Limit, params.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query articles: %w", err)
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var a model.Article
		err := rows.Scan(
			&a.ID,
			&a.AuthorID,
			&a.Title,
			&a.Body,
			&a.CreatedAt,
			&a.Author.ID,
			&a.Author.Name,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("scan article: %w", err)
		}
		articles = append(articles, a)
	}

	return articles, total, nil
}

func (r *articleRepository) Create(ctx context.Context, article model.Article) (*model.Article, error) {
	article.ID = uuid.New()
	article.CreatedAt = time.Now()

	query := fmt.Sprintf(`
        INSERT INTO %s (id, author_id, title, body, created_at)
        VALUES ($1, $2, $3, $4, $5);
    `, article.TableName())

	_, err := r.db.ExecContext(ctx, query,
		article.ID,
		article.AuthorID,
		article.Title,
		article.Body,
		article.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("fail create article: %w", err)
	}

	return &article, nil
}
