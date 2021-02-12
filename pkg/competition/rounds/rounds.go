package rounds

import (
	"context"
)

type Service interface {
	FindByIDs(ctx context.Context, ids []string) ([]Round, error)
	FindGradeRounds(ctx context.Context, gradeIDs []string) (map[string][]Round, error)
	Update(ctx context.Context, req *UpdateRequest) (interface{}, error)
}

type ServiceImpl struct {
	UpdateService *UpdateService
	Repository    Repository
}

func (s ServiceImpl) FindByIDs(ctx context.Context, ids []string) ([]Round, error) {
	return s.Repository.FindRoundsByID(ctx, ids)
}

func (s ServiceImpl) FindGradeRounds(ctx context.Context, gradeIDs []string) (map[string][]Round, error) {
	return s.Repository.FindGradeRounds(ctx, gradeIDs)
}

func (s ServiceImpl) Update(ctx context.Context, req *UpdateRequest) (interface{}, error) {
	return s.UpdateService.Update(ctx, req)
}

func NewServiceImpl(updateService *UpdateService, repo Repository) *ServiceImpl {
	return &ServiceImpl{UpdateService: updateService, Repository: repo}
}
