package ports

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/google/uuid"
)

type Ports struct {
	IdempotencyID string    `json:"idempotency_id" bson:"idempotency_id"`
	Name          string    `json:"name" bson:"name"`
	City          string    `json:"city" bson:"city"`
	Country       string    `json:"country" bson:"country"`
	Coordinates   []float64 `json:"coordinates" bson:"coordinates"`
	Province      string    `json:"province" bson:"province"`
	Timezone      string    `json:"timezone" bson:"timezone"`
	Unlocs        []string  `json:"unlocs" bson:"unlocs"`
	Code          string    `json:"code" bson:"code"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

func (p *Ports) SetID() {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}

	p.IdempotencyID = uuid.NewSHA1(uuid.Nil, buf.Bytes()).String()
}
