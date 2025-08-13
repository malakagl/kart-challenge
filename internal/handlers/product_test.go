package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/malakagl/kart-challenge/internal/models"
	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) FindAll() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) FindByID(productID string) (*models.Product, error) {
	args := m.Called(productID)
	return args.Get(0).(*models.Product), args.Error(1)
}

func TestListProducts_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productHandler := NewProductHandler(mockRepo)

	products := []models.Product{
		{ID: "123", Name: "Product 1", Price: 100},
		{ID: "456", Name: "Product 2", Price: 200},
	}
	mockRepo.On("FindAll").Return(products, nil)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	productHandler.ListProducts(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedProducts []models.Product
	err := json.NewDecoder(resp.Body).Decode(&fetchedProducts)
	assert.NoError(t, err)
	assert.Equal(t, products, fetchedProducts)
}

func TestListProducts_Failure(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productHandler := NewProductHandler(mockRepo)

	mockRepo.On("FindAll").Return([]models.Product{}, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	productHandler.ListProducts(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetProductByID_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productHandler := NewProductHandler(mockRepo)

	product := &models.Product{ID: "123", Name: "Product 1", Price: 100}
	mockRepo.On("FindByID", "123").Return(product, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("productID", "123")
	req := httptest.NewRequest(http.MethodGet, "/products/123", nil)
	req = req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	productHandler.GetProductByID(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var fetchedProduct models.Product
	err := json.NewDecoder(resp.Body).Decode(&fetchedProduct)
	assert.NoError(t, err)
	assert.Equal(t, product, &fetchedProduct)
}

func TestGetProductByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productHandler := NewProductHandler(mockRepo)

	mockRepo.On("FindByID", "123").Return(&models.Product{}, constants.ErrProductNotFound)

	req := httptest.NewRequest(http.MethodGet, "/products/123", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("productID", "123")
	req = req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	productHandler.GetProductByID(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetProductByID_InvalidID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	productHandler := NewProductHandler(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/products/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("productID", "")
	req = req.WithContext(context.WithValue(context.Background(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	productHandler.GetProductByID(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
