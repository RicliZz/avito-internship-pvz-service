package receptionService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
)

type ReceptionService struct {
	ReceptionRepo repositories.ReceptionRepo
}

func NewReceptionService(receptionRepo repositories.ReceptionRepo) *ReceptionService {
	return &ReceptionService{
		ReceptionRepo: receptionRepo,
	}
}
