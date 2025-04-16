package repositories

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
)

type AuthenticationRepository interface {
	Register(payload models.RegisterParams) error
	GetUserByEmail(email string) (error, string, string)
}

type PVZRepository interface {
	CreatePVZ(pvz models.CreatePVZRequest) (error, *models.PVZ)
}
