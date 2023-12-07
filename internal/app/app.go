package app

import (
	"go.uber.org/zap"
	grpcapp "proteiservice/internal/app/grpc"
	"proteiservice/internal/config"
	"proteiservice/internal/domain/models"
	"proteiservice/internal/services/absences"
	"proteiservice/internal/transport"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *zap.Logger, cfg config.Config, requestQueue chan models.Request, resultQueue chan models.ResultRequest) *App {
	httpManager := transport.New(log, cfg.HTTP.Ip, cfg.HTTP.Port, cfg.AuthData.Login, cfg.AuthData.Password)
	absenceService := absences.New(log, httpManager, httpManager, cfg.Emojis, requestQueue, resultQueue)
	grpcApp := grpcapp.New(log, absenceService, cfg.GRPC.Port)

	return &App{GRPCServer: grpcApp}
}
