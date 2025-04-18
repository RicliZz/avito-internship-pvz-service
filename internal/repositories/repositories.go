package repositories

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/google/uuid"
)

type AuthenticationRepo interface {
	Register(payload models.RegisterParams) (error, *models.User)
	GetUserByEmail(email string) (error, string, string)
}

type PVZRepo interface {
	CreatePVZ(pvz models.CreatePVZRequest) (error, *models.PVZ)
	GetListPVZ(filters models.QueryParamForGetPVZList) (error, []models.ListPVZResponse)
}

type ReceptionRepo interface {
	CreateReception(reception models.CreateReceptionRequest) (error, *models.Reception)
	FindLastActiveReception(PVZId uuid.UUID) (error, uuid.UUID)
	DeleteLastProduct(PVZId uuid.UUID) error
	CloseLastReception(PVZId uuid.UUID) (error, *models.Reception)
}

type ProductRepo interface {
	AddProductInActiveReception(receptionID uuid.UUID, productType string) (error, *models.Product)
}
