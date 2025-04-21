package pvzRepo

import (
	"context"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PVZRepository struct {
	db *pgxpool.Pool
}

func NewPVZRepository(db *pgxpool.Pool) *PVZRepository {
	return &PVZRepository{
		db: db,
	}
}

func (r *PVZRepository) CreatePVZ(payload models.CreatePVZRequest) (*models.PVZ, error) {
	logger.Logger.Info("CreatePVZ repository was started")
	var newPVZ models.PVZ
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO "PVZ" (city) VALUES ($1)
		RETURNING "ID", "registrationDate", city`,
		payload.City).Scan(&newPVZ.ID, &newPVZ.RegistrationDate, &newPVZ.City)
	if err != nil {
		logger.Logger.Error("Failed create new PVZ")
		return nil, err
	}
	return &newPVZ, nil
}

func (r *PVZRepository) GetListPVZ(params models.QueryParamForGetPVZList) ([]models.ListPVZResponse, error) {
	logger.Logger.Info("GetListPVZ repository was started")
	offset := (params.Page - 1) * params.Limit

	sqlQuery := `
		SELECT "pvzID", "registrationDate", city,
		       r."ID", r."dateTime", r.status,
		       p."ID", p."dateTime", p.type
		FROM "PVZ" pvz
		JOIN reception r ON pvz."ID" = r."pvzID"
		JOIN products p ON p."receptionID" = r."ID"
		WHERE r."dateTime" BETWEEN $1 AND $2
		ORDER BY r."dateTime" DESC 
		LIMIT $3
		OFFSET $4
	`

	rows, err := r.db.Query(context.Background(), sqlQuery, params.StartDate, params.EndDate, params.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	forDeleteDuplicatesPVZ := make(map[uuid.UUID]*models.ListPVZResponse)
	forDeleteDuplicatesReception := make(map[uuid.UUID]*models.ListReceptionResponse)
	receptionAlreadyAdded := make(map[uuid.UUID]bool)

	var res []models.ListPVZResponse

	for rows.Next() {
		var (
			pvzID            uuid.UUID
			registrationDate time.Time
			city             string
			receptionID      uuid.UUID
			receptionDate    time.Time
			status           string
			productID        uuid.UUID
			productDate      time.Time
			productType      string
		)

		err = rows.Scan(&pvzID, &registrationDate, &city,
			&receptionID, &receptionDate, &status,
			&productID, &productDate, &productType)
		if err != nil {
			return nil, err
		}

		findPVZ, ok := forDeleteDuplicatesPVZ[pvzID]
		if !ok {
			findPVZ = &models.ListPVZResponse{
				PVZ: models.PVZ{
					ID:               pvzID,
					RegistrationDate: registrationDate,
					City:             city,
				},
				Receptions: []*models.ListReceptionResponse{},
			}
			forDeleteDuplicatesPVZ[pvzID] = findPVZ
		}

		findReception, ok := forDeleteDuplicatesReception[receptionID]
		if !ok {
			findReception = &models.ListReceptionResponse{
				Reception: models.Reception{
					ID:       receptionID,
					DateTime: receptionDate,
					PVZId:    pvzID,
					Status:   status,
				},
				Products: []models.Product{},
			}
			forDeleteDuplicatesReception[receptionID] = findReception
		}

		findReception.Products = append(findReception.Products, models.Product{
			ID:          productID,
			ProductType: productType,
			DateTime:    productDate,
		})

		if !receptionAlreadyAdded[receptionID] {
			findPVZ.Receptions = append(findPVZ.Receptions, findReception)
			receptionAlreadyAdded[receptionID] = true
		}
	}

	for _, v := range forDeleteDuplicatesPVZ {
		res = append(res, *v)
	}

	return res, nil
}
