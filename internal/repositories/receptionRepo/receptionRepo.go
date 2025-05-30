package receptionRepo

import (
	"context"
	"errors"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
)

// Для тестов, работает и без него, если поле db структуры репозитория - *pgxpool.Pool
type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type ReceptionRepository struct {
	db Querier
}

func NewReceptionRepository(db Querier) *ReceptionRepository {
	return &ReceptionRepository{
		db: db,
	}
}

func (r *ReceptionRepository) CreateReception(payload models.CreateReceptionRequest) (*models.Reception, error) {
	var newReception models.Reception
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		log.Println("Не удалось начать транзакцию")
		return nil, err
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
			return nil, errors.New("По данному ПВЗ уже идёт приёмка ")
		}
		log.Println("Ошибка в SQL при создании новой приёмки")
		return nil, err
	}
	if err = tx.Commit(context.Background()); err != nil {
		return nil, err
	}
	return &newReception, nil
}

func (r *ReceptionRepository) FindLastActiveReception(PVZId uuid.UUID) (uuid.UUID, error) {
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
			return uuid.Nil, errors.New("Нет открытых приёмок, добавлять товар некуда")
		}
		return uuid.Nil, err
	}
	return receptionID, nil
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

func (r *ReceptionRepository) CloseLastReception(pvzID uuid.UUID) (*models.Reception, error) {
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
			return nil, err
		}
		logger.Logger.Error("Error when updating reception")
		return nil, err
	}
	return &updatedReception, nil
}
