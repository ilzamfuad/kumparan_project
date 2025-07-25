// Code generated by MockGen. DO NOT EDIT.
// Source: repository/article_repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	model "kumparan_project/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockArticleRepositoryInterface is a mock of ArticleRepositoryInterface interface.
type MockArticleRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockArticleRepositoryInterfaceMockRecorder
}

// MockArticleRepositoryInterfaceMockRecorder is the mock recorder for MockArticleRepositoryInterface.
type MockArticleRepositoryInterfaceMockRecorder struct {
	mock *MockArticleRepositoryInterface
}

// NewMockArticleRepositoryInterface creates a new mock instance.
func NewMockArticleRepositoryInterface(ctrl *gomock.Controller) *MockArticleRepositoryInterface {
	mock := &MockArticleRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockArticleRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleRepositoryInterface) EXPECT() *MockArticleRepositoryInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockArticleRepositoryInterface) Create(ctx context.Context, article model.Article) (*model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, article)
	ret0, _ := ret[0].(*model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockArticleRepositoryInterfaceMockRecorder) Create(ctx, article interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockArticleRepositoryInterface)(nil).Create), ctx, article)
}

// Search mocks base method.
func (m *MockArticleRepositoryInterface) Search(ctx context.Context, params model.ArticleQueryParams) ([]model.Article, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", ctx, params)
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Search indicates an expected call of Search.
func (mr *MockArticleRepositoryInterfaceMockRecorder) Search(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockArticleRepositoryInterface)(nil).Search), ctx, params)
}
