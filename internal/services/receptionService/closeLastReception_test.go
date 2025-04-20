package receptionService

import (
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

type MockReceptionRepository struct {
	mock.Mock
}

func (m *MockReceptionRepository) FindLastActiveReception(PVZId uuid.UUID) (error, uuid.UUID) {
	//TODO implement me
	panic("implement me")
}

func (m *MockReceptionRepository) CloseLastReception(PVZId uuid.UUID) (error, *models.Reception) {
	args := m.Called(PVZId)
	return args.Error(0), args.Get(1).(*models.Reception)
}

func TestCloseLastReception(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)

	mockReceptionRepo := new(MockReceptionRepository)

	service := ReceptionService{ReceptionRepo: mockReceptionRepo}

	pvzID := uuid.New()
	closedReception := &models.Reception{
		ID:     uuid.New(),
		PVZId:  pvzID,
		Status: "closed",
	}

	req := httptest.NewRequest(http.MethodPost, "/receptions/"+pvzID.String()+"/close", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	c.Params = gin.Params{gin.Param{Key: "pvzId", Value: pvzID.String()}}

	mockReceptionRepo.On("CloseLastReception", pvzID).Return(nil, closedReception)

	service.CloseLastReception(c)

	require.Equal(t, http.StatusOK, w.Code)

	var response models.Reception
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Equal(t, closedReception.ID, response.ID)
	require.Equal(t, closedReception.Status, response.Status)
	require.Equal(t, closedReception.PVZId, response.PVZId)

	mockReceptionRepo.AssertExpectations(t)
}

func TestCloseLastReception_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)

	mockReceptionRepo := new(MockReceptionRepository)
	service := ReceptionService{ReceptionRepo: mockReceptionRepo}

	invalidPVZID := "bad_uuid"

	req := httptest.NewRequest(http.MethodPost, "/receptions/"+invalidPVZID+"/close", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "pvzId", Value: invalidPVZID}}

	service.CloseLastReception(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	require.Contains(t, w.Body.String(), "Invalid request")

	// Убедимся, что mock не вызывался, потому что UUID невалиден
	mockReceptionRepo.AssertNotCalled(t, "CloseLastReception")
}
