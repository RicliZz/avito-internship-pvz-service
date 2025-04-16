package JWT

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func CreateDummyJWT(role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString, err
}
