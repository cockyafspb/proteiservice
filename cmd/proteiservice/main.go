package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"proteiservice/internal/app"
	"proteiservice/internal/config"
	"proteiservice/internal/domain/models"
	"syscall"
)

const (
	envLocal  = "local"
	envDev    = "dev"
	envProd   = "prod"
	NoAbsence = -1
)

func main() {
	cfg := config.MustLoad()
	logger := createLogger(cfg.Env)

	logger.Info("Setup succeeded")

	requestQueue := make(chan models.Request, cfg.QueueSize)
	resultQueue := make(chan models.ResultRequest, cfg.QueueSize)

	application := app.New(logger, *cfg, requestQueue, resultQueue)

	emojis := cfg.Emojis

	// Start worker goroutines to process requests from the queue
	for i := 0; i < cfg.WorkersNumber; i++ {
		go func() {
			for req := range requestQueue {
				employee, err := req.EmployeeGetter.GetEmployee(req.Email)
				if err != nil {
					resultQueue <- models.ResultRequest{Err: err}
				}
				id, err := req.AbsenceGetter.GetAbsence(employee)
				if err != nil {
					resultQueue <- models.ResultRequest{Err: err}
				}
				if id == NoAbsence {
					resultQueue <- models.ResultRequest{Name: employee.Name}
				}
				resultQueue <- models.ResultRequest{Name: employee.Name + emojis[id], Ok: true}
			}
		}()
	}

	go func() { application.GRPCServer.MustRun() }()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	logger.Info("Stopping application", zap.String("signal", sign.String()))

	application.GRPCServer.Stop()

	logger.Info("Application stopped")
}

func createLogger(env string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	var config zap.Config

	switch env {
	case envLocal:
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "console",
			EncoderConfig:     encoderCfg,
			OutputPaths: []string{
				"stderr",
			},
			ErrorOutputPaths: []string{
				"stderr",
			},
		}
	case envDev:
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths: []string{
				"stderr",
			},
			ErrorOutputPaths: []string{
				"stderr",
			},
		}
	case envProd:
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Sampling:          nil,
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths: []string{
				"stderr",
			},
			ErrorOutputPaths: []string{
				"stderr",
			},
		}
	}

	return zap.Must(config.Build())
}
