package productRepo

import (
	"context"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) AddProductInActiveReception(receptionID uuid.UUID, productType string) (*models.Product, error) {
	logger.Logger.Info("AddProductInActiveReception repository was started")
	var newProduct models.Product
	sqlQuery := `INSERT INTO products (type, "receptionID") VALUES ($1, $2) RETURNING "ID", "dateTime", type, "receptionID"`
	err := r.db.QueryRow(context.Background(), sqlQuery,
		productType, receptionID).Scan(&newProduct.ID, &newProduct.DateTime, &newProduct.ProductType, &newProduct.ReceptionId)
	if err != nil {
		logger.Logger.Error("AddProductInActiveReception repository failed to insert product", err)
		return nil, err
	}
	return &newProduct, nil
}
