package authService

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/pass"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAuthenticationRepository struct {
	mock.Mock
}

func (m *MockAuthenticationRepository) Register(payload models.RegisterParams) (*models.User, error) {
	args := m.Called(payload)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthenticationRepository) GetUserByEmail(email string) (error, string, string) {
	args := m.Called(email)
	return args.Error(0), args.String(1), args.String(2)
}

func TestLoginService_Login(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	gin.SetMode(gin.TestMode)

	paramPass, _ := pass.CreateHash("password123")

	mockRepo := new(MockAuthenticationRepository)

	loginService := &AuthLoginService{
		authDB: mockRepo,
	}

	t.Run("Valid credentials", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "user@example.com").Return(nil, paramPass, "moderator")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		loginParams := models.LoginParams{
			Email:    "user@example.com",
			Password: "password123", // Правильный пароль
		}
		loginParamsJSON, err := json.Marshal(loginParams)
		if err != nil {
			t.Fatalf("Failed to marshal login params: %v", err)
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginParamsJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		loginService.Login(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Invalid password", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "user@example.com").Return(nil, paramPass, "moderator")

		// Создаем запрос
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		loginParams := models.LoginParams{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}
		loginParamsJSON, err := json.Marshal(loginParams)
		if err != nil {
			t.Fatalf("Failed to marshal login params: %v", err)
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginParamsJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		loginService.Login(c)

		assert.Equal(t, 401, w.Code)
		assert.JSONEq(t, `{"message":"Invalid credentials"}`, w.Body.String())
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "user@example.com").Return(errors.New("user not found"), "", "")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		loginParams := models.LoginParams{
			Email: "user@example.com",
		}
		loginParamsJSON, err := json.Marshal(loginParams)
		if err != nil {
			t.Fatalf("Failed to marshal login params: %v", err)
		}
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(loginParamsJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		loginService.Login(c)

		assert.Equal(t, 401, w.Code)
		assert.JSONEq(t, `{"message":"Invalid credentials"}`, w.Body.String())

	})
}
