package main

import (
	"context"
	"errors"
	"flag"
	"github.com/guil95/ports/config/storages/mongo"
	"go.uber.org/zap"

	"github.com/guil95/ports/config/logger"
	"github.com/guil95/ports/internal/ports"
	"github.com/guil95/ports/internal/ports/infra/filereader"
	"github.com/guil95/ports/internal/ports/infra/repository"
	"github.com/guil95/ports/internal/ports/usecase"
	_ "github.com/joho/godotenv/autoload"
)

var file string

func init() {
	flag.StringVar(&file, "file", "ports.json", "ports file path")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.SetupLogger(ctx)

	uc := usecase.NewUseCase(repository.NewMongo(mongo.Connect(ctx)), filereader.NewReadFile())

	portChan := make(chan ports.Ports)
	errChan := make(chan error)

	go func(errChan chan error) {
		err := uc.SavePorts(ctx, file, portChan, errChan)
		if err != nil {
			errChan <- err
		}
	}(errChan)

	zap.S().Info("starting to save ports")

	for {
		select {
		case err := <-errChan:
			if errors.Is(err, ports.EofError{}) {
				zap.S().Info("end of save ports")
				return
			}

			zap.L().Error("error to save or update ports", zap.Any("error", err))
			return
		}
	}
}
