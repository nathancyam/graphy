package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphy/transport/graphql/generated"
	"graphy/transport/graphql/model"
)

func (r *mutationResolver) UpdateRound(ctx context.Context, id string, input model.RoundInput) (*model.RoundUpdateResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Rounds(ctx context.Context, ids []string) ([]*model.Round, error) {
	l := r.LogDuration(ctx, "Rounds")
	defer l()

	c, err := r.RoundRepo.FindRoundsByID(ctx, ids)
	if err != nil {
		return nil, err
	}

	var rl []*model.Round
	for _, i := range c {
		rr := model.Round(i)
		rl = append(rl, &rr)
	}

	return rl, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
