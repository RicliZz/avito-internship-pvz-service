package productService

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockProductRepository struct {
	mock.Mock
}

type MockReceptionRepository struct {
	mock.Mock
}

func (m *MockReceptionRepository) CreateReception(reception models.CreateReceptionRequest) (*models.Reception, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockReceptionRepository) FindLastActiveReception(PVZId uuid.UUID) (uuid.UUID, error) {
	args := m.Called(PVZId)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockReceptionRepository) DeleteLastProduct(PVZId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockReceptionRepository) CloseLastReception(PVZId uuid.UUID) (*models.Reception, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockProductRepository) AddProductInActiveReception(receptionID uuid.UUID, productType string) (*models.Product, error) {
	args := m.Called(receptionID, productType)
	return args.Get(0).(*models.Product), args.Error(1)
}

func TestAddProductInReception(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()

	gin.SetMode(gin.TestMode)

	mockReceptionRepo := new(MockReceptionRepository)
	mockProductRepo := new(MockProductRepository)

	service := ProductService{ProductRepo: mockProductRepo, ReceptionRepository: mockReceptionRepo}
	pvzID := uuid.New()
	receptionID := uuid.New()
	product := &models.Product{
		ID:          uuid.New(),
		ProductType: "одежда",
	}
	requestBody := models.AddProductRequest{
		Type:  product.ProductType,
		PvzID: pvzID,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockReceptionRepo.On("FindLastActiveReception", pvzID).Return(receptionID, nil)
	mockProductRepo.On("AddProductInActiveReception", receptionID, product.ProductType).Return(product, nil)
	service.AddProductInReception(c)
	require.Equal(t, http.StatusCreated, w.Code)
	var response models.Product
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, product.ID, response.ID)
	require.Equal(t, product.ProductType, response.ProductType)

	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertExpectations(t)
}

func TestAddProductInReception_NotFoundReception(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()

	gin.SetMode(gin.TestMode)

	mockReceptionRepo := new(MockReceptionRepository)
	mockProductRepo := new(MockProductRepository)

	service := ProductService{
		ProductRepo:         mockProductRepo,
		ReceptionRepository: mockReceptionRepo,
	}

	pvzID := uuid.New()
	product := &models.Product{
		ProductType: "одежда",
	}

	requestBody := models.AddProductRequest{
		Type:  product.ProductType,
		PvzID: pvzID,
	}

	bodyBytes, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Возвращаем ошибку при поиске приёмки
	mockReceptionRepo.On("FindLastActiveReception", pvzID).Return(uuid.Nil, errors.New("not found"))

	service.AddProductInReception(c)

	require.Equal(t, 400, w.Code)

	mockReceptionRepo.AssertExpectations(t)
	mockProductRepo.AssertNotCalled(t, "AddProductInActiveReception", mock.Anything, mock.Anything)
}
