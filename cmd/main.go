package main

import (
	"context"
	"errors"
	"flag"
	"github.com/guil95/ports/config/storages/mongo"

	"github.com/guil95/ports/config/logger"
	"github.com/guil95/ports/internal/ports"
	"github.com/guil95/ports/internal/ports/infra/filereader"
	"github.com/guil95/ports/internal/ports/infra/repository"
	"github.com/guil95/ports/internal/ports/usecase"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
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

	err := uc.SavePorts(ctx, file)
	if err != nil {
		if errors.Is(err, ports.EofError{}) {
			zap.S().Info("end of save ports")
			return
		}

		zap.S().Error(err)
		return
	}
}
