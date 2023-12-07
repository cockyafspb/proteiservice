package grpcapp

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	absencesgrpc "proteiservice/internal/grpc/absences"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates new gRPC server app
func New(log *zap.Logger, absences absencesgrpc.Absences, port int) *App {
	gRPCServer := grpc.NewServer()

	absencesgrpc.Register(gRPCServer, absences)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) run() error {
	log := a.log.With(zap.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	log.Info("gRPC server is running", zap.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	a.log.Info("Stopping gRPC server", zap.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
