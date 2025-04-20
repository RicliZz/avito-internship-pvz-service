package authService

import (
	"bytes"
	"encoding/json"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDummyLogin(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)
	params := models.DummyLoginParams{
		Role: "moderator",
	}
	body, _ := json.Marshal(&params)
	req, err := http.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	server := DummyLoginService{}
	server.DummyLogin(c)

	assert.Equal(t, 200, w.Code)

	var token string
	err = json.Unmarshal(w.Body.Bytes(), &token)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestDummyLogin_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)

	body := []byte(`{}`)

	req, err := http.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	server := DummyLoginService{}
	server.DummyLogin(c)
	assert.Equal(t, 400, w.Code)

	var resp models.Error
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Message)
}
