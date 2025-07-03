package listing

import (
	"article/config"
	"context"
)

// RepositoryMySQL provides access to MySQL repository
type RepositoryMySQL interface {
	ReadArticles(context.Context, FilterGetArticle) ([]Article, error)
}

// RepositoryRedis provides access to Redis repository
type RepositoryRedis interface {
	CreateArticle(context.Context, []Article, FilterGetArticle)
	ReadArticles(context.Context, FilterGetArticle) ([]Article, error)
}

type Service interface {
	GetArticles(context.Context, FilterGetArticle) ([]Article, error)
}

type service struct {
	rmy RepositoryMySQL
	rre RepositoryRedis
}

func NewService(rmy RepositoryMySQL, rre RepositoryRedis) Service {
	return &service{rmy, rre}
}

func (s *service) GetArticles(ctx context.Context, lfgar FilterGetArticle) ([]Article, error) {
	lars, err := s.rre.ReadArticles(ctx, lfgar)
	if err != nil {
		if err.Error() == config.RedisErrKeyDoesNotExist {
			err = nil

			lars, err := s.rmy.ReadArticles(ctx, lfgar)
			if err != nil {
				return lars, err
			}

			s.rre.CreateArticle(ctx, lars, lfgar)
			return lars, err
		}
	}

	return lars, err
}
