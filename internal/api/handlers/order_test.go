package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) Create(order models.Order) (string, error) {
	args := m.Called(order)
	return args.String(0), args.Error(1)
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) FindByID(productID uint) (*models.Product, error) {
	args := m.Called(productID)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) FindAll() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

type MockCouponValidator struct {
	mock.Mock
}

func (m *MockCouponValidator) ValidateCouponCode(code string) bool {
	args := m.Called(code)
	return args.Bool(0)
}

func TestCreateOrder_Success(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)
	mockCouponValidator := new(MockCouponValidator)

	orderHandler := NewOrderHandler(mockOrderService, mockProductService, mockCouponValidator)

	orderReq := models.OrderRequest{
		Items: []models.OrderProduct{
			{ProductID: "123", Quantity: 2},
		},
		CouponCode: "FIFTYOFF",
	}

	product := &models.Product{ID: 123, Name: "Test Product", Price: 100}
	mockProductService.On("FindByID", uint(123)).Return(product, nil)
	mockOrderService.On("Create", mock.Anything).Return("order123", nil)
	mockCouponValidator.On("ValidateCouponCode", "FIFTYOFF").Return(true)

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var orderResp models.OrderResponse
	err := json.NewDecoder(resp.Body).Decode(&orderResp)
	assert.NoError(t, err)
	assert.Equal(t, "order123", orderResp.ID)
	assert.Equal(t, []*models.Product{product}, orderResp.Products)
}

func TestCreateOrder_InvalidRequestBody(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)

	orderHandler := NewOrderHandler(mockOrderService, mockProductService, nil)

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateOrder_InvalidCouponCode(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)
	mockCouponValidator := new(MockCouponValidator)

	orderHandler := NewOrderHandler(mockOrderService, mockProductService, mockCouponValidator)

	orderReq := models.OrderRequest{
		Items: []models.OrderProduct{
			{ProductID: "123", Quantity: 2},
		},
		CouponCode: "INVALIDCODE",
	}

	product := &models.Product{ID: 123, Name: "Test Product", Price: 100}
	mockProductService.On("FindByID", 123).Return(product, nil)
	mockOrderService.On("Create", mock.Anything).Return("order123", nil)
	mockCouponValidator.On("ValidateCouponCode", "INVALIDCODE").Return(false)

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestCreateOrder_InvalidProductID(t *testing.T) {
	mockOrderService := new(MockOrderService)
	mockProductService := new(MockProductService)
	mockCouponValidator := new(MockCouponValidator)
	mockCouponValidator.On("ValidateCouponCode", "VALIDCODE").Return(true)
	orderHandler := NewOrderHandler(mockOrderService, mockProductService, mockCouponValidator)

	orderReq := models.OrderRequest{
		CouponCode: "VALIDCODE",
		Items: []models.OrderProduct{
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
	mockCouponValidator := new(MockCouponValidator)
	mockCouponValidator.On("ValidateCouponCode", "VALIDCODE").Return(true)
	orderHandler := NewOrderHandler(mockOrderService, mockProductService, mockCouponValidator)

	orderReq := models.OrderRequest{
		CouponCode: "VALIDCODE",
		Items: []models.OrderProduct{
			{ProductID: "123", Quantity: 2},
		},
	}

	mockProductService.On("FindByID", uint(123)).Return(&models.Product{}, constants.ErrProductNotFound)

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
	mockCouponValidator := new(MockCouponValidator)
	mockCouponValidator.On("ValidateCouponCode", "VALIDCODE").Return(true)
	orderHandler := NewOrderHandler(mockOrderService, mockProductService, mockCouponValidator)

	orderReq := models.OrderRequest{
		CouponCode: "VALIDCODE",
		Items: []models.OrderProduct{
			{ProductID: "123", Quantity: 2},
		},
	}

	product := &models.Product{ID: 123, Name: "Test Product", Price: 100}
	mockProductService.On("FindByID", uint(123)).Return(product, nil)
	mockOrderService.On("Create", mock.Anything).Return("", assert.AnError)

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	orderHandler.CreateOrder(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
