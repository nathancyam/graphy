package grades

import "context"

type Service struct {
	Repository Repository
}

func (s Service) FindByID(ctx context.Context, id string) (*Grade, error) {
	return s.Repository.FindByID(ctx, id)
}

func NewService(repository Repository) *Service {
	return &Service{Repository: repository}
}
