package productService

import (
	"bytes"
	"encoding/json"
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

func (m *MockReceptionRepository) CreateReception(reception models.CreateReceptionRequest) (error, *models.Reception) {
	//TODO implement me
	panic("implement me")
}

func (m *MockReceptionRepository) FindLastActiveReception(PVZId uuid.UUID) (error, uuid.UUID) {
	args := m.Called(PVZId)
	return args.Error(0), args.Get(1).(uuid.UUID)
}

func (m *MockReceptionRepository) DeleteLastProduct(PVZId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockReceptionRepository) CloseLastReception(PVZId uuid.UUID) (error, *models.Reception) {
	//TODO implement me
	panic("implement me")
}

func (m *MockProductRepository) AddProductInActiveReception(receptionID uuid.UUID, productType string) (error, *models.Product) {
	args := m.Called(receptionID, productType)
	return nil, args.Get(1).(*models.Product)
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

	mockReceptionRepo.On("FindLastActiveReception", pvzID).Return(nil, receptionID)
	mockProductRepo.On("AddProductInActiveReception", receptionID, product.ProductType).Return(nil, product)
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
