package receptionService

import (
	"bytes"
	"encoding/json"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (m *MockReceptionRepository) CreateReception(reception models.CreateReceptionRequest) (*models.Reception, error) {
	args := m.Called(reception)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func TestCreateReception(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)
	mockReceptionRepo := new(MockReceptionRepository)

	service := ReceptionService{ReceptionRepo: mockReceptionRepo}
	receptionRequest := models.CreateReceptionRequest{PVZId: uuid.New()}

	newReception := &models.Reception{
		ID:     uuid.New(),
		PVZId:  receptionRequest.PVZId,
		Status: "in_progress",
	}
	bodyBytes, _ := json.Marshal(receptionRequest)
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	mockReceptionRepo.On("CreateReception", receptionRequest).Return(newReception, nil)

	service.CreateReception(c)

	require.Equal(t, http.StatusCreated, w.Code)

	var response models.Reception
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, newReception.ID, response.ID)
	require.Equal(t, newReception.Status, response.Status)
	require.Equal(t, newReception.PVZId, response.PVZId)

	mockReceptionRepo.AssertExpectations(t)
}

func TestCreateReception_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)
	mockReceptionRepo := new(MockReceptionRepository)

	service := ReceptionService{ReceptionRepo: mockReceptionRepo}

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/receptions", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	service.CreateReception(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid request")
}
