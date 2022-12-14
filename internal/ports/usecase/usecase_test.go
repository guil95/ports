package usecase

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/guil95/ports/internal/ports"
	"github.com/guil95/ports/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	savePort            = "SavePort"
	findByIdempotencyID = "FindByIdempotencyID"
	readFileStream      = "ReadFileStream"
)

func TestUseCase(t *testing.T) {
	repoMock := new(mocks.Repository)
	fileReaderMock := new(mocks.FileReader)
	ctx := context.Background()

	portChan := make(chan ports.Ports)
	errChan := make(chan error)
	idempotencyID := getIdempotencyID(ports.Ports{City: "Maringá"})
	portsModel := ports.Ports{City: "Maringá", IdempotencyID: idempotencyID}
	filePath := "ports_test.json"

	t.Run("test create ports with valid ports should return eof error", func(t *testing.T) {
		fileReaderMock.On(readFileStream, filePath, portChan, errChan)
		repoMock.On(savePort, ctx, &portsModel).Return(nil)
		repoMock.On(findByIdempotencyID, ctx, idempotencyID).Return(nil, nil)

		uc := NewUseCase(repoMock, fileReaderMock)

		go func() {
			err := uc.SavePorts(ctx, filePath, portChan, errChan)
			assert.ErrorIs(t, err, ports.EofError{})
		}()

		time.Sleep(time.Millisecond * 10)
		portChan <- portsModel
		time.Sleep(time.Millisecond * 10)
		errChan <- ports.EofError{}
	})
}

func getIdempotencyID(v interface{}) string {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(ports.Ports{City: "Maringá"})

	return uuid.NewSHA1(uuid.Nil, buf.Bytes()).String()
}
