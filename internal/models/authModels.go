package models

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type RegisterParams struct {
	Email    string ` json:"email" binding:"required,email"`
	Password string ` json:"password" binding:"required,min=8"`
	Role     string ` json:"role" binding:"required,oneof=employee moderator"`
}

type LoginParams struct {
	Email    string ` json:"email" binding:"required,email"`
	Password string ` json:"password" binding:"required"`
}

type DummyLoginParams struct {
	Role string ` json:"role" binding:"required,oneof=employee moderator"` //Возможные роли: модератор/работник
}
