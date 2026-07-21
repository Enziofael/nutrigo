package v1

import (
	"context"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
)

type CommonService struct {
	repo repository.CommonRepository
}

func NewCommonService(repo repository.CommonRepository) *CommonService {
	return &CommonService{repo: repo}
}

func (s *CommonService) Ping(ctx context.Context) (OK bool, err error) {
	return s.repo.Ping(ctx)
}