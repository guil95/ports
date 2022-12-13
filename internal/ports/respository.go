package ports

import "context"

type Repository interface {
	SavePort(ctx context.Context, ports *Ports) error
	FindByIdempotencyID(ctx context.Context, idempotencyID string) (*Ports, error)
	UpdatePort(ctx context.Context, ports *Ports) error
}
