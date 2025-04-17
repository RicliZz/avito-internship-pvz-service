package pvzRepo

import (
	"context"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/jackc/pgx/v5"
	"log"
)

type PVZRepository struct {
	db *pgx.Conn
}

func NewPVZRepository(db *pgx.Conn) *PVZRepository {
	return &PVZRepository{
		db: db,
	}
}

func (r *PVZRepository) CreatePVZ(payload models.CreatePVZRequest) (error, *models.PVZ) {
	log.Println("Репозиторий создания нового ПВЗ")
	var newPVZ models.PVZ
	err := r.db.QueryRow(context.Background(), `
		INSERT INTO "PVZ" (city) VALUES ($1)
		RETURNING "ID", "registrationDate", city`,
		payload.City).Scan(&newPVZ.ID, &newPVZ.RegistrationDate, &newPVZ.City)
	if err != nil {
		log.Println("Ошибка SQL запроса на создание нового ПВЗ")
		return err, nil
	}
	return nil, &newPVZ
}
