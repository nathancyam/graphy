package rounds

import (
	"context"
)

type PolicyHandler interface {
	Do(ctx context.Context, req *UpdateRequest) (context.Context, error)
}

type UpdateRequest struct {
	RoundID string
	Payload struct{}
}

type UpdateService struct {
	Repository Repository
}

func NewUpdateService(repository Repository) *UpdateService {
	return &UpdateService{Repository: repository}
}

func (s *UpdateService) Update(ctx context.Context, req *UpdateRequest) (interface{}, error) {
	return nil, nil
}
