package repositories

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
)

type AuthenticationRepository interface {
	Register(payload models.RegisterParams) error
	GetUserByEmail(email string) (string, string, error)
}
