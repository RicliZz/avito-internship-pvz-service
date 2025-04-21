package pvzRepo

import (
	"errors"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestCreatePVZ(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PVZRepository{db: mock}

	city := "Москва"
	pvzID := uuid.New()
	registrationDate := time.Now()

	rows := pgxmock.NewRows([]string{"ID", "registrationDate", "city"}).
		AddRow(pvzID, registrationDate, city)

	mock.ExpectQuery(`INSERT INTO "PVZ" \(city\) VALUES \(\$1\) RETURNING "ID", "registrationDate", city`).
		WithArgs(city).
		WillReturnRows(rows)

	result, err := repo.CreatePVZ(models.CreatePVZRequest{City: city})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, pvzID, result.ID)
	assert.WithinDuration(t, registrationDate, result.RegistrationDate, time.Second)
	assert.Equal(t, city, result.City)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePVZ_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &PVZRepository{db: mock}

	city := "Самара" //неверный город

	mock.ExpectQuery(`INSERT INTO "PVZ" \(city\) VALUES \(\$1\) RETURNING "ID", "registrationDate", city`).
		WithArgs(city).
		WillReturnError(errors.New("city constraint violation"))

	result, err := repo.CreatePVZ(models.CreatePVZRequest{City: city})

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}
