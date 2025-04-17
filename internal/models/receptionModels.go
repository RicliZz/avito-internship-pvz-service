package models

import (
	"github.com/google/uuid"
	"time"
)

type Reception struct {
	ID       uuid.UUID `json:"id"`
	DateTime time.Time `json:"datetime" bson:"datetime"`
	PVZId    uuid.UUID `json:"pvzId" binding:"required"`
	Status   string    `json:"status" binding:"required,oneof= in_progress close"`
}

type CreateReceptionRequest struct {
	PVZId uuid.UUID `json:"pvzId" binding:"required"`
}
