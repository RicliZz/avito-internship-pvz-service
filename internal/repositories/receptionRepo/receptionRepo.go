package receptionRepo

import (
	"context"
	"errors"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/jackc/pgx/v5"
	"log"
)

type ReceptionRepository struct {
	db *pgx.Conn
}

func NewReceptionRepository(db *pgx.Conn) *ReceptionRepository {
	return &ReceptionRepository{
		db: db,
	}
}

func (r *ReceptionRepository) CreateReception(payload models.CreateReceptionRequest) (error, *models.Reception) {
	var newReception models.Reception
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		log.Println("Не удалось начать транзакцию")
		return err, nil
	}
	defer tx.Rollback(context.Background())
	sqlQuery := `
				WITH in_progress_reception AS (
    			SELECT id FROM reception
    			WHERE pvzid = $1 AND status = 'in_progress'
        		FOR UPDATE
				)
				INSERT INTO reception (pvzid)
				SELECT $1
				WHERE NOT EXISTS(SELECT 1 FROM in_progress_reception)
				RETURNING *;
				`
	err = tx.QueryRow(context.Background(), sqlQuery,
		payload.PVZId).Scan(&newReception.ID, &newReception.DateTime, &newReception.PVZId, &newReception.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Есть незакрытая приёмка")
			return errors.New("По данному ПВЗ уже идёт приёмка "), nil
		}
		log.Println("Ошибка в SQL при создании новой приёмки")
		return err, nil
	}
	if err = tx.Commit(context.Background()); err != nil {
		return err, nil
	}
	return nil, &newReception
}
