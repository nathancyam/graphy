package rounds

import "context"

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

func (s *UpdateService) Update(ctx context.Context, req *UpdateRequest) {
}
