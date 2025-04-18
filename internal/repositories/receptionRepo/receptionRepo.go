package receptionRepo

import (
	"context"
	"errors"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
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
	var receptionID uuid.UUID
	sqlQuery := `SELECT "ID" FROM reception
				 WHERE "pvzID" = $1 AND status = 'in_progress' 
				 ORDER BY "dateTime" DESC LIMIT 1;`
	err := r.db.QueryRow(context.Background(), sqlQuery, PVZId).Scan(&receptionID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("Нет открытых приёмок, добавлять товар некуда"), uuid.Nil
		}
		return err, uuid.Nil
	}
	return nil, receptionID
}

func (r *ReceptionRepository) DeleteLastProduct(pvzID uuid.UUID) error {
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
		log.Println("Ошибка в SQL при удалении товара")
		return err
	}
	if tag.RowsAffected() == 0 {
		log.Println("Нет открытых приёмок, удалено 0 строк")
		return errors.New("Нет открытых приёмок ")
	}
	if err = tx.Commit(context.Background()); err != nil {
		log.Println("Ошибка при коммите")
		return err
	}
	return nil
}

func (r *ReceptionRepository) CloseLastReception(pvzID uuid.UUID) (error, *models.Reception) {
	var updatedReception models.Reception
	sqlQuery := `UPDATE reception
				SET status = 'close'
				WHERE "pvzID" = $1 AND status = 'in_progress'
				RETURNING *`
	err := r.db.QueryRow(context.Background(), sqlQuery,
		pvzID).Scan(&updatedReception.ID, &updatedReception.DateTime, &updatedReception.PVZId, &updatedReception.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Нет открытых приёмок, обновлено 0 строк")
			return err, nil
		}
		log.Println("Ошибка при SQL запросе")
		return err, nil
	}
	return nil, &updatedReception
}
