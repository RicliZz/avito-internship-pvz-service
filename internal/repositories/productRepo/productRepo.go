package productRepo

import (
	"context"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) AddProductInActiveReception(receptionID uuid.UUID, productType string) (error, *models.Product) {
	var newProduct models.Product
	sqlQuery := `INSERT INTO products (type, "receptionID") VALUES ($1, $2) RETURNING "ID", "dateTime", type, "receptionID"`
	err := r.db.QueryRow(context.Background(), sqlQuery,
		productType, receptionID).Scan(&newProduct.ID, &newProduct.DateTime, &newProduct.ProductType, &newProduct.ReceptionId)
	if err != nil {
		return err, nil
	}
	return nil, &newProduct
}
