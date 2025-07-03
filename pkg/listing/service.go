package listing

import "context"

// RepositoryMySQL provides access to MySQL repository
type RepositoryMySQL interface {
	ReadArticles(context.Context, FilterGetArticle) ([]Article, error)
}

type Service interface {
	GetArticles(ctx context.Context, lfga FilterGetArticle) ([]Article, error)
}

type service struct {
	rmy RepositoryMySQL
}

func NewService(rmy RepositoryMySQL) Service {
	return &service{rmy}
}

func (s *service) GetArticles(ctx context.Context, lfga FilterGetArticle) ([]Article, error) {

	lars, err := s.rmy.ReadArticles(ctx, lfga)
	if err != nil {
		return lars, err
	}

	return lars, err
}
