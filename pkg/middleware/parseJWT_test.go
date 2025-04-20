package middleware

import (
	"github.com/RicliZz/avito-internship-pvz-service/pkg/JWT"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type TestsCases struct {
	Header string
	Status int
}

func TestCheckRoleMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "super_secret_key")
	gin.SetMode(gin.TestMode)
	secret := "super_secret_key"
	expectedRole := "moderator"

	testHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	}
	commonResp, _ := JWT.CreateJWT(expectedRole)

	testsCases := []TestsCases{
		{Header: "", Status: 401},
		{Header: "wrongHeader", Status: 401},
		{Header: "Bearer " + commonResp, Status: 200},
	}
	for _, testCase := range testsCases {
		t.Run(testCase.Header, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.GET("/protected", CheckRoleMiddleware(secret, expectedRole), testHandler)

			req, _ := http.NewRequest("GET", "/protected", nil)
			if testCase.Header != "" {
				req.Header.Set("Authorization", testCase.Header)
			}
			c.Request = req
			r.ServeHTTP(w, req)

			require.Equal(t, testCase.Status, w.Code)
		})
	}
}
