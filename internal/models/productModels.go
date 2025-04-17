package models

import (
	"github.com/google/uuid"
	"time"
)

type AddProductRequest struct {
	Type  string    `json:"type" binding:"required,oneof=электроника одежда обувь"`
	PvzID uuid.UUID `json:"pvzId" binding:"required"`
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	ProductType string    `json:"type"`
	ReceptionId uuid.UUID `json:"receptionId"`
}
