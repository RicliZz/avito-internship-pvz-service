package authRepo

import (
	"fmt"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestRegister(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &AuthRepository{db: mock}

	email := "test@example.com"
	password := "hashedpassword"
	role := "user"

	userID := uuid.New()

	rows := pgxmock.NewRows([]string{"ID", "email", "role"}).
		AddRow(userID, email, role)

	mock.ExpectQuery(`INSERT INTO users \(email, password, role\) VALUES \(\$1, \$2, \$3\) RETURNING "ID", email, role`).
		WithArgs(email, password, role).
		WillReturnRows(rows)

	user, err := repo.Register(models.RegisterParams{
		Email:    email,
		Password: password,
		Role:     role,
	})

	assert.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, role, user.Role)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegister_Fail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &AuthRepository{db: mock}

	email := "test@example.com"
	password := "hashedpassword"
	role := "user"

	mock.ExpectQuery(`INSERT INTO users \(email, password, role\) VALUES \(\$1, \$2, \$3\) RETURNING "ID", email, role`).
		WithArgs(email, password, role).
		WillReturnError(fmt.Errorf("database error"))

	user, err := repo.Register(models.RegisterParams{
		Email:    email,
		Password: password,
		Role:     role,
	})

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "database error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &AuthRepository{db: mock}

	email := "test@example.com"
	password := "hashedpassword"
	role := "user"

	mock.ExpectQuery(`SELECT password, role FROM users WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(pgxmock.NewRows([]string{"password", "role"}).
			AddRow(password, role))

	resultPassword, resultRole, err := repo.GetUserByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, password, resultPassword)
	assert.Equal(t, role, resultRole)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_UserNotFound(t *testing.T) {
	logger.Logger = zap.NewNop().Sugar()
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &AuthRepository{db: mock}

	email := "нет_почты@example.com"

	mock.ExpectQuery(`SELECT password, role FROM users WHERE email = \$1`).
		WithArgs(email).
		WillReturnError(pgx.ErrNoRows)

	resultPassword, resultRole, err := repo.GetUserByEmail(email)

	assert.Error(t, err)
	assert.Equal(t, "", resultPassword)
	assert.Equal(t, "", resultRole)
	assert.Contains(t, err.Error(), "пользователь с email")

	assert.NoError(t, mock.ExpectationsWereMet())
}
