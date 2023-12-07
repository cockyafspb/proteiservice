package app

import (
	"go.uber.org/zap"
	grpcapp "proteiservice/internal/app/grpc"
	"proteiservice/internal/config"
	"proteiservice/internal/services/absences"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *zap.Logger, cfg config.Config) *App {
	// TODO: init absences service
	absenceService := absences.New(log, cfg.HTTP.Ip, cfg.HTTP.Port, cfg.AuthData.Login, cfg.AuthData.Password)
	grpcApp := grpcapp.New(log, absenceService, cfg.GRPC.Port)

	return &App{GRPCServer: grpcApp}
}
