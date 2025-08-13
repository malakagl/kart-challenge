package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models2 "github.com/malakagl/kart-challenge/pkg/models"

	"github.com/malakagl/kart-challenge/pkg/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) Create(order models2.Order) (string, error) {
	args := m.Called(order)
	return args.String(0), args.Error(1)
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) FindByID(productID string) (*models2.Product, error) {
	args := m.Called(productID)
	return args.Get(0).(*models2.Product), args.Error(1)
}

func (m *MockProductService) FindAll() ([]models2.Product, error) {
	args := m.Called()
	return args.Get(0).([]models2.Product), args.Error(1)
}

func TestCreateOrder_Success(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)

	orderHandler := NewOrderHandler(mockOrderService, mockProductService)

	orderReq := models2.OrderRequest{
		Items: []models2.OrderItem{
			{ProductID: "123", Quantity: 2},
		},
		CouponCode: "DISCOUNT10",
	}

	product := &models2.Product{ID: "123", Name: "Test Product", Price: 100}
	mockProductService.On("FindByID", "123").Return(product, nil)
	mockOrderService.On("Create", mock.Anything).Return("order123", nil)

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var orderResp models2.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&orderResp)
	assert.NoError(t, err)
	assert.Equal(t, "order123", orderResp.ID)
	assert.Equal(t, orderReq.Items, orderResp.Items)
	assert.Equal(t, []*models2.Product{product}, orderResp.Products)
}

func TestCreateOrder_InvalidRequestBody(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)

	orderHandler := NewOrderHandler(mockOrderService, mockProductService)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateOrder_InvalidProductID(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)

	orderHandler := NewOrderHandler(mockOrderService, mockProductService)

	orderReq := models2.OrderRequest{
		Items: []models2.OrderItem{
			{ProductID: "", Quantity: 2},
		},
	}

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateOrder_ProductNotFound(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)

	orderHandler := NewOrderHandler(mockOrderService, mockProductService)

	orderReq := models2.OrderRequest{
		Items: []models2.OrderItem{
			{ProductID: "123", Quantity: 2},
		},
	}

	mockProductService.On("FindByID", "123").Return(&models2.Product{}, constants.ErrProductNotFound)

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateOrder_FailedToCreateOrder(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)

	orderHandler := NewOrderHandler(mockOrderService, mockProductService)

	orderReq := models2.OrderRequest{
		Items: []models2.OrderItem{
			{ProductID: "123", Quantity: 2},
		},
	}

	product := &models2.Product{ID: "123", Name: "Test Product", Price: 100}
	mockProductService.On("FindByID", "123").Return(product, nil)
	mockOrderService.On("Create", mock.Anything).Return("", assert.AnError)

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
