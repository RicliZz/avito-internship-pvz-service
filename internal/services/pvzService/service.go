package pvzService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
)

type PVZService struct {
	PVZRepo repositories.PVZRepo
}

func NewPVZService(PVZRepo repositories.PVZRepo) *PVZService {
	return &PVZService{
		PVZRepo: PVZRepo,
	}
}
