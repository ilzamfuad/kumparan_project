package controller

import (
	"encoding/json"
	"kumparan_project/model"
	"kumparan_project/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const (
	cacheExpiration = 5 * time.Minute
)

type ArticleController struct {
	articleService service.ArticleServiceInterface
	redis          *redis.Client
}

func NewArticleController(articleService service.ArticleServiceInterface, redis *redis.Client) *ArticleController {
	return &ArticleController{articleService, redis}
}

func (ac *ArticleController) GetArticles(c *gin.Context) {
	query := c.Query("query")
	author := c.Query("author")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	queryParams := model.ArticleQueryParams{
		Query:  &query,
		Author: &author,
		Limit:  &limit,
		Offset: &offset,
	}

	cacheKey := "articles:" + query + ":" + author + ":page:" + strconv.Itoa(page) + ":limit:" + strconv.Itoa(limit)
	cachedData, err := ac.redis.Get(cacheKey).Result()
	if err == nil {
		// If cache exists, return the cached response
		var cachedResponse struct {
			Data       []model.Article `json:"data"`
			Pagination gin.H           `json:"pagination"`
		}
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			c.JSON(http.StatusOK, cachedResponse)
			return
		}
	}

	articles, total, err := ac.articleService.GetArticles(c.Request.Context(), queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	response := gin.H{
		"data": articles,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
			"pages": totalPages,
		},
	}

	// Cache the response
	responseJSON, _ := json.Marshal(response)
	ac.redis.Set(cacheKey, responseJSON, cacheExpiration)

	c.JSON(http.StatusOK, response)
}

func (ac *ArticleController) CreateArticles(c *gin.Context) {
	var body model.ArticleCreateParams

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := ac.articleService.CreateArticles(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": model.ArticleCreateResponse{
		ID:        article.ID,
		AuthorID:  article.AuthorID,
		Title:     article.Title,
		Body:      article.Body,
		CreatedAt: article.CreatedAt,
	}})
}
