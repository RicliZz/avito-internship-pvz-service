package grpcServices

import (
	"context"
	pvz_v1 "github.com/RicliZz/avito-internship-pvz-service/pkg/proto"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type GRPCserver struct {
	pvz_v1.UnimplementedPVZServiceServer
	Db *pgx.Conn
}

func (s *GRPCserver) GetPVZList(ctx context.Context, in *pvz_v1.GetPVZListRequest) (*pvz_v1.GetPVZListResponse, error) {
	sqlQuery := `SELECT "ID", "registrationDate", city FROM "PVZ"`
	rows, err := s.Db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pvzs := []*pvz_v1.PVZ{}
	for rows.Next() {
		var (
			ID               string
			registrationDate time.Time
			city             string
		)
		if err = rows.Scan(&ID, &registrationDate, &city); err != nil {
			return nil, err
		}
		pvzs = append(pvzs, &pvz_v1.PVZ{
			Id:               ID,
			RegistrationDate: timestamppb.New(registrationDate),
			City:             city,
		})
	}
	return &pvz_v1.GetPVZListResponse{Pvzs: pvzs}, nil
}
