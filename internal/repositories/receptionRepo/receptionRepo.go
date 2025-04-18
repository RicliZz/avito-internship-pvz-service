package receptionRepo

import (
	"context"
	"errors"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/google/uuid"
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
    			SELECT "ID" FROM reception
    			WHERE "pvzID" = $1 AND status = 'in_progress'
        		FOR UPDATE
				)
				INSERT INTO reception ("pvzID")
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

func (r *ReceptionRepository) FindLastActiveReception(PVZId uuid.UUID) (error, uuid.UUID) {
	logger.Logger.Info("FindLastActiveReception repository was started")
	var receptionID uuid.UUID
	sqlQuery := `SELECT "ID" FROM reception
				 WHERE "pvzID" = $1 AND status = 'in_progress' 
				 ORDER BY "dateTime" DESC LIMIT 1;`
	err := r.db.QueryRow(context.Background(), sqlQuery, PVZId).Scan(&receptionID)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Debugw("Don't find active reception with",
				"pvzID", PVZId)
			return errors.New("Нет открытых приёмок, добавлять товар некуда"), uuid.Nil
		}
		return err, uuid.Nil
	}
	return nil, receptionID
}

func (r *ReceptionRepository) DeleteLastProduct(pvzID uuid.UUID) error {
	logger.Logger.Info("DeleteLastProduct repository was started")
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	sqlQuery := `WITH reception_id AS (
					SELECT "ID" FROM reception
					WHERE "pvzID" = $1 AND status = 'in_progress'
					FOR UPDATE
				)
				
				DELETE FROM products
				WHERE "ID" = (
					SELECT "ID" FROM products
					WHERE "receptionID" = (SELECT "ID" FROM reception_id)
					ORDER BY "dateTime" DESC LIMIT 1
				)`
	tag, err := tx.Exec(context.Background(), sqlQuery, pvzID)
	if err != nil {
		logger.Logger.Error("Error when SQL deleting product from reception")
		return err
	}
	if tag.RowsAffected() == 0 {
		logger.Logger.Debugw("Don't find active reception, deleted 0 rows")
		return errors.New("Нет открытых приёмок ")
	}
	if err = tx.Commit(context.Background()); err != nil {
		logger.Logger.Error("Failed commit transaction")
		return err
	}
	return nil
}

func (r *ReceptionRepository) CloseLastReception(pvzID uuid.UUID) (error, *models.Reception) {
	logger.Logger.Info("CloseLastReception repository was started")
	var updatedReception models.Reception
	sqlQuery := `UPDATE reception
				SET status = 'close'
				WHERE "pvzID" = $1 AND status = 'in_progress'
				RETURNING *`
	err := r.db.QueryRow(context.Background(), sqlQuery,
		pvzID).Scan(&updatedReception.ID, &updatedReception.DateTime, &updatedReception.PVZId, &updatedReception.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Debugw("Don't find active reception, closed 0 rows")
			return err, nil
		}
		logger.Logger.Error("Error when updating reception")
		return err, nil
	}
	return nil, &updatedReception
}
