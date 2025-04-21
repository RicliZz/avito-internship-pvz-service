package receptionRepo

import (
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestCloseLastReception(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &ReceptionRepository{db: mock}

	pvzID := uuid.New()
	receptionID := uuid.New()
	receptionDateTime := time.Now()

	rows := pgxmock.NewRows([]string{"ID", "dateTime", "pvzID", "status"}).
		AddRow(receptionID, receptionDateTime, pvzID, "close")

	mock.ExpectQuery(`UPDATE reception SET status = 'close' WHERE "pvzID" = \$1 AND status = 'in_progress' RETURNING \*`).
		WithArgs(pvzID).
		WillReturnRows(rows)

	reception, err := repo.CloseLastReception(pvzID)

	assert.NoError(t, err)
	assert.NotNil(t, reception)
	assert.Equal(t, receptionID, reception.ID)
	assert.Equal(t, pvzID, reception.PVZId)
	assert.Equal(t, "close", reception.Status)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCloseLastReception_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &ReceptionRepository{db: mock}

	pvzID := uuid.New()

	//не нашлось строк для обновления, все приёмки и так уже закрыты
	mock.ExpectQuery(`UPDATE reception SET status = 'close' WHERE "pvzID" = \$1 AND status = 'in_progress' RETURNING \*`).
		WithArgs(pvzID).
		WillReturnError(pgx.ErrNoRows)

	reception, err := repo.CloseLastReception(pvzID)

	assert.Error(t, err)
	assert.Nil(t, reception)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLastActiveReception_Success(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &ReceptionRepository{db: mock}

	pvzID := uuid.New()
	receptionID := uuid.New()

	rows := pgxmock.NewRows([]string{"ID"}).
		AddRow(receptionID)

	mock.ExpectQuery(`SELECT "ID" FROM reception WHERE "pvzID" = \$1 AND status = 'in_progress' ORDER BY "dateTime" DESC LIMIT 1`).
		WithArgs(pvzID).
		WillReturnRows(rows)

	result, err := repo.FindLastActiveReception(pvzID)

	assert.NoError(t, err)
	assert.Equal(t, receptionID, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindLastActiveReception_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &ReceptionRepository{db: mock}

	pvzID := uuid.New()

	//нет активных приёмок
	mock.ExpectQuery(`SELECT "ID" FROM reception WHERE "pvzID" = \$1 AND status = 'in_progress' ORDER BY "dateTime" DESC LIMIT 1`).
		WithArgs(pvzID).
		WillReturnError(pgx.ErrNoRows)

	receptionID, err := repo.FindLastActiveReception(pvzID)

	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, receptionID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
