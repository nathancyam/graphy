package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"graphy/transport/graphql/dataloader"
	"graphy/transport/graphql/generated"
	"graphy/transport/graphql/model"
)

func (r *gradeResolver) Rounds(ctx context.Context, obj *model.Grade) ([]*model.Round, error) {
	return dataloader.For(ctx).RoundsByID.Load(obj.ID)
}

func (r *mutationResolver) UpdateRound(ctx context.Context, id string, input model.RoundInput) (*model.RoundUpdateResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Grade(ctx context.Context, id *string) (*model.Grade, error) {
	l := r.LogDuration(ctx, "Grades")
	defer l()

	grade, err := r.GradeSvc.FindByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	g := model.Grade(*grade)
	return &g, nil
}

func (r *queryResolver) Rounds(ctx context.Context, ids []string) ([]*model.Round, error) {
	l := r.LogDuration(ctx, "RoundService")
	defer l()

	rounds, err := r.RoundService.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	var rl = make([]*model.Round, len(rounds))
	for i, round := range rounds {
		rr := model.Round(round)
		rl[i] = &rr
	}

	return rl, nil
}

// Grade returns generated.GradeResolver implementation.
func (r *Resolver) Grade() generated.GradeResolver { return &gradeResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type gradeResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
