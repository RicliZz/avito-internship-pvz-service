package authRepo

import (
	"context"
	"fmt"
	"github.com/RicliZz/avito-internship-pvz-service/internal/models"
	"github.com/jackc/pgx/v5"
	"log"
)

type AuthRepository struct {
	db *pgx.Conn
}

func NewAuthRepository(db *pgx.Conn) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Register(payload models.RegisterParams) error {
	log.Println("Register rep is start")
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO users (email, password, role) VALUES ($1, $2, $3)`,
		payload.Email, payload.Password, payload.Role)

	if err != nil {
		log.Println("Error with INSERT in USERS")
		return err
	}
	return nil
}

func (r *AuthRepository) GetUserByEmail(email string) (error, string, string) {
	log.Println("GetUserByEmail is start")
	var password string
	var role string

	err := r.db.QueryRow(context.Background(), `SELECT password, role FROM users WHERE email = $1`,
		email).Scan(&password, &role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("пользователь с email %s не найден", email), "", ""
		}
		return err, "", ""
	}

	return nil, password, role
}
