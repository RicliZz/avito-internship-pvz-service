package integration_tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func TestIntegrationAddProductInReception(t *testing.T) {
	connStr := "postgres://postgres:123123123@localhost:5432/avito_internship_test?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	require.NoError(t, err)
	defer conn.Close(context.Background())

	cities := []string{"Москва", "Санкт-Петербург", "Казань"}
	products := []string{"одежда", "электроника", "обувь"}

	pvzCity := cities[rand.Intn(len(cities))]
	productTypes := make([]string, 50)
	for i := 0; i < 50; i++ {
		productTypes[i] = products[rand.Intn(len(products))]
	}

	var pvzID uuid.UUID

	sqlInsertPVZ := `INSERT INTO "PVZ" ("city") VALUES ($1) RETURNING "ID";`
	err = conn.QueryRow(context.Background(), sqlInsertPVZ, pvzCity).Scan(&pvzID)
	require.NoError(t, err)

	sqlInsertReception := `INSERT INTO reception ("pvzID", status) VALUES ($1, 'in_progress') RETURNING "ID";`
	var receptionID uuid.UUID
	err = conn.QueryRow(context.Background(), sqlInsertReception, pvzID).Scan(&receptionID)
	require.NoError(t, err)

	for i := 0; i < 50; i++ {

		sqlInsertProduct := `INSERT INTO products (type, "receptionID") VALUES ($1, $2) RETURNING "ID";`
		var productID uuid.UUID
		err = conn.QueryRow(context.Background(), sqlInsertProduct, productTypes[i], receptionID).Scan(&productID)
		require.NoError(t, err)

	}

	sqlCountProducts := `SELECT COUNT(*) FROM products WHERE "receptionID" = $1;`
	var productCount int
	err = conn.QueryRow(context.Background(), sqlCountProducts, receptionID).Scan(&productCount)
	require.NoError(t, err)
	require.Equal(t, 50, productCount)

	sqlUpdateReception := `UPDATE reception SET status = 'close' WHERE "ID" = $1;`
	_, err = conn.Exec(context.Background(), sqlUpdateReception, receptionID)
	require.NoError(t, err)

	sqlCheckStatus := `SELECT status FROM reception WHERE "ID" = $1;`
	var status string
	err = conn.QueryRow(context.Background(), sqlCheckStatus, receptionID).Scan(&status)
	require.NoError(t, err)
	require.Equal(t, "close", status)
}
