package rounds

import "context"

type Repository interface {
	FindRoundsByID(ctx context.Context, roundIDs []string) ([]Round, error)
}
