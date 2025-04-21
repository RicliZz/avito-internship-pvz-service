package models

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type PVZ struct {
	ID               uuid.UUID `json:"id,omitempty"`
	RegistrationDate time.Time `json:"registrationDate,omitempty"`
	City             string    `json:"city"`
}

type CreatePVZRequest struct {
	City string `json:"city" binding:"required,oneof=Москва Санкт-Петербург Казань"`
}

type ListPVZResponse struct {
	PVZ
	Receptions []*ListReceptionResponse
}

type QueryParamForGetPVZList struct {
	StartDate time.Time `form:"startDate" time_format:"2006-01-02T15:04:05Z"`
	EndDate   time.Time `form:"endDate" time_format:"2006-01-02T15:04:05Z"`
	Page      int       `form:"page" binding:"min=1"`
	Limit     int       `form:"limit" binding:"min=1,max=30"`
}

var DatesForGetPVZList validator.Func = func(fl validator.FieldLevel) bool {
	params, ok := fl.Field().Interface().(QueryParamForGetPVZList)
	if !ok {
		return false
	}

	if params.StartDate.IsZero() || params.EndDate.IsZero() {
		return true
	}

	return !params.StartDate.After(params.EndDate)
}
