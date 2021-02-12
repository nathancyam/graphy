package grades

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Grade, error)
}
