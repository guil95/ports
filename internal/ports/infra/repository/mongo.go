package repository

import (
	"context"
	"os"
	"time"

	"github.com/guil95/ports/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type mongoRepository struct {
	db *mongo.Client
}

func NewMongo(db *mongo.Client) ports.Repository {
	return &mongoRepository{db}
}

const (
	collectionPorts          = "ports"
	collectionPortsDuplicate = "ports_duplicate"
)

func (m mongoRepository) SavePort(ctx context.Context, ports *ports.Ports) error {
	ports.CreatedAt = time.Now()
	ports.UpdatedAt = time.Now()
	db := m.db.Database(os.Getenv("DB_DATABASE"))

	_, err := db.Collection(collectionPorts).InsertOne(ctx, ports)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	return nil
}

func (m mongoRepository) FindByIdempotencyID(ctx context.Context, idempotencyID string) (*ports.Ports, error) {
	db := m.db.Database(os.Getenv("DB_DATABASE"))

	var p ports.Ports

	err := db.Collection(collectionPorts).FindOne(
		ctx,
		bson.D{{"idempotency_id", idempotencyID}},
	).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return &p, nil
}

func (m mongoRepository) UpdatePort(ctx context.Context, model *ports.Ports) error {
	session, err := m.db.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	transactional := func(sessCtx mongo.SessionContext) (interface{}, error) {
		db := m.db.Database(os.Getenv("DB_DATABASE"))
		portsColl := db.Collection(collectionPortsDuplicate)
		portsDuplicateColl := db.Collection(collectionPortsDuplicate)

		_, err = portsDuplicateColl.InsertOne(ctx, model)
		if err != nil {
			zap.S().Error(err)
			return nil, err
		}

		model.UpdatedAt = time.Now()
		updateFields := bson.M{"$set": bson.M{
			"city":           model.City,
			"code":           model.Code,
			"coordinates":    model.Coordinates,
			"country":        model.Country,
			"idempotency_id": model.IdempotencyID,
			"name":           model.Name,
			"province":       model.Province,
			"timezone":       model.Timezone,
			"unlocs":         model.Unlocs,
		}}
		_, err = portsColl.UpdateOne(ctx, bson.D{{"idempotency_id", model.IdempotencyID}}, updateFields)
		if err != nil {
			zap.S().Error(err)
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, transactional)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	return nil
}
