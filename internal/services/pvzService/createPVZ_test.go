package pvzService

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

type MockPVZRepository struct {
	mock.Mock
}

func (m *MockPVZRepository) CreatePVZ(pvz models.CreatePVZRequest) (error, *models.PVZ) {
	args := m.Called(pvz)
	return args.Error(0), args.Get(1).(*models.PVZ)
}

func (m *MockPVZRepository) GetListPVZ(filters models.QueryParamForGetPVZList) (error, []models.ListPVZResponse) {
	//TODO implement me
	panic("implement me")
}

func TestCreatePVZ(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)
	mockRepo := new(MockPVZRepository)
	service := &PVZService{PVZRepo: mockRepo}
	expectedPVZ := &models.PVZ{
		ID:   uuid.New(),
		City: "Москва",
	}
	reqBody := models.CreatePVZRequest{City: "Москва"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/pvz", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockRepo.On("CreatePVZ", reqBody).Return(nil, expectedPVZ)

	service.CreatePVZ(c)

	require.Equal(t, http.StatusCreated, w.Code)
	var response models.PVZ
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, expectedPVZ.ID, response.ID)
	require.Equal(t, expectedPVZ.City, response.City)

	mockRepo.AssertExpectations(t)
}
