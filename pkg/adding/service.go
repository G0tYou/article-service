package adding

import (
	"context"
)

// RepositoryMySQL provides access to MySQL repository
type RepositoryMySQL interface {
	CreateArticle(context.Context, Article) (int, error)
}

// RepositoryRedis provides access to Redis repository
type RepositoryRedis interface {
	DeleteArticle(context.Context)
}

// Service provides adding operations
type Service interface {
	AddArticle(context.Context, Article) (int, error)
}

type service struct {
	rmy RepositoryMySQL
	rre RepositoryRedis
}

// NewService creates an adding service with the necessary dependencies
func NewService(rmy RepositoryMySQL, rre RepositoryRedis) Service {
	return &service{rmy, rre}
}

// AddArticle persists the given article to repository
func (s *service) AddArticle(ctx context.Context, a Article) (int, error) {
	id, err := s.rmy.CreateArticle(ctx, a)
	if err != nil {
		return 0, err
	}

	s.rre.DeleteArticle(ctx)

	return id, nil
}
