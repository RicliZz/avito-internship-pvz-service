package receptionService

import (
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (m *MockReceptionRepository) DeleteLastProduct(PVZId uuid.UUID) error {
	args := m.Called(PVZId)
	return args.Error(0)
}

func TestDeleteLastProductInActiveReception(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)
	mockReceptionRepo := new(MockReceptionRepository)
	service := ReceptionService{ReceptionRepo: mockReceptionRepo}
	pvzID := uuid.New()

	mockReceptionRepo.On("DeleteLastProduct", pvzID).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "pvzId", Value: pvzID.String()}}

	service.DeleteLastProductInActiveReception(c)

	require.Equal(t, http.StatusOK, w.Code)

	require.Equal(t, `"Товар удалён"`, w.Body.String())

	mockReceptionRepo.AssertExpectations(t)
}

func TestDeleteLastProductInActiveReception_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)

	service := ReceptionService{ReceptionRepo: new(MockReceptionRepository)}
	invalidPVZID := "bad_uuid"
	req := httptest.NewRequest(http.MethodPost, "/pvz/"+invalidPVZID+"/delete_last_product", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "pvzId", Value: "invalid-uuid"}}

	service.DeleteLastProductInActiveReception(c)

	require.Equal(t, 400, w.Code)
	require.Contains(t, w.Body.String(), "Invalid request")
}
