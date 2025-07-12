package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kumparan_project/controller"
	"kumparan_project/model"
	mock_service "kumparan_project/tests/mock/service"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestArticleController_GetArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis, err := miniredis.Run()
	assert.NoError(t, err)
	defer mockRedis.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(), // like "localhost:6379"
	})

	mockService := mock_service.NewMockArticleServiceInterface(ctrl)
	articleController := controller.NewArticleController(mockService, rdb)

	router := gin.Default()
	router.GET("/articles", articleController.GetArticles)

	t.Run("success with cache", func(t *testing.T) {
		cacheKey := "articles:::page:1:limit:10"
		err = rdb.Set(cacheKey, `{"data":[],"pagination":{"page":1,"limit":10,"total":0,"pages":0}}`, 1*time.Minute).Err()
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodGet, "/articles?page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		rdb.FlushAll()
	})

	t.Run("success with cache without page and limit", func(t *testing.T) {
		cacheKey := "articles:::page:1:limit:10"
		err = rdb.Set(cacheKey, `{"data":[],"pagination":{"page":1,"limit":10,"total":0,"pages":0}}`, 1*time.Minute).Err()
		assert.NoError(t, err)

		req, _ := http.NewRequest(http.MethodGet, "/articles?page=0&limit=0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		rdb.FlushAll()
	})

	t.Run("success without cache", func(t *testing.T) {
		mockService.EXPECT().GetArticles(gomock.Any(), gomock.Any()).Return([]model.Article{}, 0, nil)

		req, _ := http.NewRequest(http.MethodGet, "/articles?page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		rdb.FlushAll()
	})

	t.Run("error from service", func(t *testing.T) {
		mockService.EXPECT().GetArticles(gomock.Any(), gomock.Any()).Return(nil, 0, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, "/articles?page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		rdb.FlushAll()
	})
}

func TestArticleController_CreateArticles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis, err := miniredis.Run()
	assert.NoError(t, err)
	defer mockRedis.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(), // like "localhost:6379"
	})

	mockService := mock_service.NewMockArticleServiceInterface(ctrl)
	articleController := controller.NewArticleController(mockService, rdb)

	router := gin.Default()
	router.POST("/articles", articleController.CreateArticles)

	t.Run("success", func(t *testing.T) {
		uuidArticle := uuid.New()
		body := model.ArticleCreateParams{
			AuthorID: uuid.New(),
			Title:    "Test Title",
			Body:     "Test Content",
		}
		bodyJSON, _ := json.Marshal(body)

		mockService.EXPECT().CreateArticles(gomock.Any(), body).Return(&model.Article{ID: uuidArticle, Title: "Test Title", Body: "Test Content"}, nil)

		req, _ := http.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("bad request", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("error from service", func(t *testing.T) {
		body := model.ArticleCreateParams{
			AuthorID: uuid.New(),
			Title:    "Test Title",
			Body:     "Test Content",
		}
		bodyJSON, _ := json.Marshal(body)

		mockService.EXPECT().CreateArticles(gomock.Any(), body).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
