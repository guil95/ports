package usecase

import (
	"context"
	"github.com/guil95/ports/internal/ports"
	"go.uber.org/zap"
)

type useCase struct {
	repo       ports.Repository
	fileReader ports.FileReader
}

func NewUseCase(repo ports.Repository, fileReader ports.FileReader) *useCase {
	return &useCase{repo: repo, fileReader: fileReader}
}

func (uc *useCase) SavePorts(
	ctx context.Context,
	filePath string,
	portChan chan ports.Ports,
	errChan chan error,
) error {
	go uc.fileReader.ReadFileStream(filePath, portChan, errChan)

	for {
		select {
		case p := <-portChan:
			err := uc.saveOrUpdate(ctx, &p)
			if err != nil {
				zap.L().Error("error to save or update port", zap.Any("error", err))
				return err
			}
		case err := <-errChan:
			return err
		}
	}
}

func (uc *useCase) saveOrUpdate(ctx context.Context, model *ports.Ports) error {
	port, err := uc.repo.FindByIdempotencyID(ctx, model.IdempotencyID)
	if err != nil {
		return err
	}

	if port != nil {
		err = uc.repo.UpdatePort(ctx, port)
		if err != nil {
			return err
		}

		return nil
	}

	err = uc.repo.SavePort(ctx, model)
	if err != nil {
		return err
	}

	return nil
}
