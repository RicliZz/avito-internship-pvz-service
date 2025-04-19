package pvzService

import (
	"github.com/RicliZz/avito-internship-pvz-service/internal/repositories"
	pvz_v1 "github.com/RicliZz/avito-internship-pvz-service/pkg/proto"
)

type PVZService struct {
	PVZRepo          repositories.PVZRepo
	PVZServiceClient pvz_v1.PVZServiceClient
}

func NewPVZService(PVZRepo repositories.PVZRepo, PVZServiceClient pvz_v1.PVZServiceClient) *PVZService {
	return &PVZService{
		PVZRepo:          PVZRepo,
		PVZServiceClient: PVZServiceClient,
	}
}
