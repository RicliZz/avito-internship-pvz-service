package repositories

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/google/uuid"
)

type AuthenticationRepo interface {
	Register(payload models.RegisterParams) (*models.User, error)
	GetUserByEmail(email string) (string, string, error)
}

type PVZRepo interface {
	CreatePVZ(pvz models.CreatePVZRequest) (*models.PVZ, error)
	GetListPVZ(filters models.QueryParamForGetPVZList) ([]models.ListPVZResponse, error)
}

type ReceptionRepo interface {
	CreateReception(reception models.CreateReceptionRequest) (*models.Reception, error)
	FindLastActiveReception(PVZId uuid.UUID) (uuid.UUID, error)
	DeleteLastProduct(PVZId uuid.UUID) error
	CloseLastReception(PVZId uuid.UUID) (*models.Reception, error)
}

type ProductRepo interface {
	AddProductInActiveReception(receptionID uuid.UUID, productType string) (*models.Product, error)
}
