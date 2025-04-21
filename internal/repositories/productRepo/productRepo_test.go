package productRepo

import (
	"fmt"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddProductInActiveReception(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &ProductRepository{db: mock}

	receptionID := uuid.New()
	productID := uuid.New()
	productType := "одежда"
	productDateTime := time.Now()

	rows := pgxmock.NewRows([]string{"ID", "dateTime", "type", "receptionID"}).
		AddRow(productID, productDateTime, productType, receptionID)

	mock.ExpectQuery(`INSERT INTO products \("type", "receptionID"\) VALUES \(\$1, \$2\) RETURNING "ID", "dateTime", "type", "receptionID"`).
		WithArgs(productType, receptionID).
		WillReturnRows(rows)

	product, err := repo.AddProductInActiveReception(receptionID, productType)
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, productID, product.ID)
	assert.Equal(t, productType, product.ProductType)
	assert.Equal(t, receptionID, product.ReceptionId)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddProductInActiveReception_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &ProductRepository{db: mock}

	receptionID := uuid.New()
	invalidProductType := "НЕ одежда обувь электроника"

	mock.ExpectQuery(`INSERT INTO products \("type", "receptionID"\) VALUES \(\$1, \$2\) RETURNING "ID", "dateTime", "type", "receptionID"`).
		WithArgs(invalidProductType, receptionID).
		WillReturnError(fmt.Errorf("invalid input value for enum product_type: \"%s\"", invalidProductType))

	product, err := repo.AddProductInActiveReception(receptionID, invalidProductType)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "invalid input value for enum")
	assert.NoError(t, mock.ExpectationsWereMet())
}
