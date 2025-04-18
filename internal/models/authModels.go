package models

type RegisterParams struct {
	Email    string ` json:"email" binding:"required,email"`
	Password string ` json:"password" binding:"required"`
	Role     string ` json:"role" binding:"required,oneof=employee moderator"`
}

type LoginParams struct {
	Email    string ` json:"email" binding:"required,email"`
	Password string ` json:"password" binding:"required"`
}

type DummyLoginParams struct {
	Role string ` json:"role" binding:"required,oneof=employee moderator"` //Возможные роли: модератор/работник
}
