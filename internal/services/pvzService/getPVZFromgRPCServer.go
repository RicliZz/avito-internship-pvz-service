package pvzService

import (
	"context"
	"github.com/RicliZz/avito-internship-pvz-service/pkg/logger"
	pvz_v1 "github.com/RicliZz/avito-internship-pvz-service/pkg/proto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (s *PVZService) GetPVZListFromRPCServer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pvzs, err := s.PVZServiceClient.GetPVZList(ctx, &pvz_v1.GetPVZListRequest{})
	if err != nil {
		logger.Logger.Info("Failed gRPC request")
	}
	c.JSON(http.StatusOK, pvzs)
}
