package adding

import (
	"context"
)

// RepositoryMySQL provides access to MySQL repository
type RepositoryMySQL interface {
	CreateArticle(context.Context, Article) (int, error)
}

// Service provides adding operations
type Service interface {
	AddArticle(context.Context, Article) (int, error)
}

type service struct {
	rmy RepositoryMySQL
}

// NewService creates an adding service with the necessary dependencies
func NewService(rmy RepositoryMySQL) Service {
	return &service{rmy}
}

// AddRewardStarPoin persists the given account & starpoin to repository
func (s *service) AddArticle(ctx context.Context, a Article) (int, error) {
	id, err := s.rmy.CreateArticle(ctx, a)
	if err != nil {
		return 0, err
	}

	return id, nil
}
