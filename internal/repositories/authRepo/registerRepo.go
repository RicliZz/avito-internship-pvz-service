package authRepo

import (
	"context"
	"fmt"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	"github.com/jackc/pgx/v5"
)

type AuthRepository struct {
	db *pgx.Conn
}

func NewAuthRepository(db *pgx.Conn) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Register(payload models.RegisterParams) (error, *models.User) {
	logger.Logger.Info("Register repository was start")
	var newUser models.User
	err := r.db.QueryRow(context.Background(),
		`INSERT INTO users (email, password, role) VALUES ($1, $2, $3) RETURNING "ID", email, role`,
		payload.Email, payload.Password, payload.Role).Scan(&newUser.ID, newUser.Email, newUser.Role)

	if err != nil {
		logger.Logger.Errorw("Failed register user",
			"email", payload.Email,
			"role", payload.Role)
		return err, nil
	}
	return nil, &newUser
}

func (r *AuthRepository) GetUserByEmail(email string) (error, string, string) {
	logger.Logger.Info("GetUserByEmail repository was started")
	var password string
	var role string

	err := r.db.QueryRow(context.Background(), `SELECT password, role FROM users WHERE email = $1`,
		email).Scan(&password, &role)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Debugw("Don't find user with email", "email", email)
			return fmt.Errorf("пользователь с email %s не найден", email), "", ""
		}
		return err, "", ""
	}

	return nil, password, role
}
