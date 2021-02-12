package rounds

import "context"

type Repository interface {
	FindRoundsByID(ctx context.Context, roundIDs []string) ([]Round, error)
	FindGradeRounds(ctx context.Context, gradeIDs []string) (map[string][]Round, error)
}
