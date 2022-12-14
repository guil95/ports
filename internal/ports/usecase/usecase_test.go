package usecase

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"github.com/google/uuid"
	"github.com/guil95/ports/internal/ports"
	"github.com/guil95/ports/mocks"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

const (
	savePort            = "SavePort"
	findByIdempotencyID = "FindByIdempotencyID"
	updatePort          = "UpdatePort"

	readFileStream = "ReadFileStream"

	genericError = "generic error"
)

func TestSaveWithoutError(t *testing.T) {
	repoMock := new(mocks.Repository)
	fileReaderMock := new(mocks.FileReader)
	ctx := context.Background()

	idempotencyID := getIdempotencyID(ports.Ports{City: "Maringá"})
	portsModel := ports.Ports{City: "Maringá", IdempotencyID: idempotencyID}
	filePath := "ports_test.json"

	portChan := make(chan ports.Ports)
	errChan := make(chan error)

	fileReaderMock.On(readFileStream, filePath, portChan, errChan)
	repoMock.On(savePort, ctx, &portsModel).Return(nil)
	repoMock.On(findByIdempotencyID, ctx, idempotencyID).Return(nil, nil)

	uc := NewUseCase(repoMock, fileReaderMock)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		err := uc.SavePorts(ctx, filePath, portChan, errChan)
		assert.ErrorIs(t, err, ports.EofError{})
		wg.Done()
	}()

	portChan <- portsModel
	time.Sleep(time.Millisecond * 10)
	errChan <- ports.EofError{}

	wg.Wait()
}

func TestSaveWithError(t *testing.T) {
	repoMock := new(mocks.Repository)
	fileReaderMock := new(mocks.FileReader)
	ctx := context.Background()

	idempotencyID := getIdempotencyID(ports.Ports{City: "Maringá"})
	portsModel := ports.Ports{City: "Maringá", IdempotencyID: idempotencyID}
	filePath := "ports_test.json"

	portChan := make(chan ports.Ports)
	errChan := make(chan error)

	fileReaderMock.On(readFileStream, filePath, portChan, errChan)
	repoMock.On(savePort, ctx, &portsModel).Return(errors.New(genericError))
	repoMock.On(findByIdempotencyID, ctx, idempotencyID).Return(nil, nil)

	uc := NewUseCase(repoMock, fileReaderMock)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		err := uc.SavePorts(ctx, filePath, portChan, errChan)
		assert.Error(t, err)
		wg.Done()
	}()

	portChan <- portsModel
	time.Sleep(time.Millisecond * 10)

	wg.Wait()
}

func TestSaveWithGenericErrorOnUpdate(t *testing.T) {
	repoMock := new(mocks.Repository)
	fileReaderMock := new(mocks.FileReader)
	ctx := context.Background()

	idempotencyID := getIdempotencyID(ports.Ports{City: "Maringá"})
	portsModel := ports.Ports{City: "Maringá", IdempotencyID: idempotencyID}
	filePath := "ports_test.json"

	portChan := make(chan ports.Ports)
	errChan := make(chan error)

	fileReaderMock.On(readFileStream, filePath, portChan, errChan)
	repoMock.On(savePort, ctx, &portsModel).Return(nil)
	repoMock.On(findByIdempotencyID, ctx, idempotencyID).Return(&portsModel, nil)
	repoMock.On(updatePort, ctx, &portsModel).Return(errors.New(genericError))

	uc := NewUseCase(repoMock, fileReaderMock)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		err := uc.SavePorts(ctx, filePath, portChan, errChan)
		assert.Error(t, err)
		wg.Done()
	}()

	portChan <- portsModel

	wg.Wait()
}

func TestSaveWithGenericErrorOnFind(t *testing.T) {
	repoMock := new(mocks.Repository)
	fileReaderMock := new(mocks.FileReader)
	ctx := context.Background()

	idempotencyID := getIdempotencyID(ports.Ports{City: "Maringá"})
	portsModel := ports.Ports{City: "Maringá", IdempotencyID: idempotencyID}
	filePath := "ports_test.json"

	portChan := make(chan ports.Ports)
	errChan := make(chan error)

	fileReaderMock.On(readFileStream, filePath, portChan, errChan)
	repoMock.On(savePort, ctx, &portsModel).Return(nil)
	repoMock.On(findByIdempotencyID, ctx, idempotencyID).Return(nil, errors.New(genericError))

	uc := NewUseCase(repoMock, fileReaderMock)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		err := uc.SavePorts(ctx, filePath, portChan, errChan)
		assert.Error(t, err)
		wg.Done()
	}()

	portChan <- portsModel

	wg.Wait()
}

func TestUpdateWithoutError(t *testing.T) {
	repoMock := new(mocks.Repository)
	fileReaderMock := new(mocks.FileReader)
	ctx := context.Background()

	idempotencyID := getIdempotencyID(ports.Ports{City: "Maringá"})
	portsModel := ports.Ports{City: "Maringá", IdempotencyID: idempotencyID}
	filePath := "ports_test.json"

	portChan := make(chan ports.Ports)
	errChan := make(chan error)

	fileReaderMock.On(readFileStream, filePath, portChan, errChan)
	repoMock.On(savePort, ctx, &portsModel).Return(nil)
	repoMock.On(findByIdempotencyID, ctx, idempotencyID).Return(&portsModel, nil)
	repoMock.On(updatePort, ctx, &portsModel).Return(nil)

	uc := NewUseCase(repoMock, fileReaderMock)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		err := uc.SavePorts(ctx, filePath, portChan, errChan)
		assert.ErrorIs(t, err, ports.EofError{})
		wg.Done()
	}()

	portChan <- portsModel
	time.Sleep(time.Millisecond * 10)
	errChan <- ports.EofError{}

	wg.Wait()
}

func TestUpdateWithError(t *testing.T) {
	repoMock := new(mocks.Repository)
	fileReaderMock := new(mocks.FileReader)
	ctx := context.Background()

	idempotencyID := getIdempotencyID(ports.Ports{City: "Maringá"})
	portsModel := ports.Ports{City: "Maringá", IdempotencyID: idempotencyID}
	filePath := "ports_test.json"

	portChan := make(chan ports.Ports)
	errChan := make(chan error)

	fileReaderMock.On(readFileStream, filePath, portChan, errChan)
	repoMock.On(savePort, ctx, &portsModel).Return(nil)
	repoMock.On(findByIdempotencyID, ctx, idempotencyID).Return(&portsModel, nil)
	repoMock.On(updatePort, ctx, &portsModel).Return(errors.New(genericError))

	uc := NewUseCase(repoMock, fileReaderMock)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		err := uc.SavePorts(ctx, filePath, portChan, errChan)
		assert.Error(t, err, ports.EofError{})
		wg.Done()
	}()

	portChan <- portsModel
	time.Sleep(time.Millisecond * 10)

	wg.Wait()
}

func getIdempotencyID(v interface{}) string {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(v)

	return uuid.NewSHA1(uuid.Nil, buf.Bytes()).String()
}
