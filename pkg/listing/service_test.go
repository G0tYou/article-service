package listing

import (
	"article/config"
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

type repositoryMySQLMock struct {
	mock.Mock
}

func (rmym *repositoryMySQLMock) ReadArticles(ctx context.Context, lfgar FilterGetArticle) ([]Article, error) {
	args := rmym.Called(ctx, lfgar)

	return args.Get(0).([]Article), args.Error(1)
}

type repositoryRedisMock struct {
	mock.Mock
}

func (rrm *repositoryRedisMock) CreateArticle(ctx context.Context, ars []Article, lfgar FilterGetArticle) {
	rrm.Called(ctx, ars, lfgar)
}

func (rrm *repositoryRedisMock) ReadArticles(ctx context.Context, lfgar FilterGetArticle) ([]Article, error) {
	args := rrm.Called(ctx, lfgar)

	return args.Get(0).([]Article), args.Error(1)
}

var (
	ctx  = context.Background()
	rmym = new(repositoryMySQLMock)
	rrdm = new(repositoryRedisMock)
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name string
		rmym *repositoryMySQLMock
		rrdm *repositoryRedisMock
		want Service
	}{
		{
			name: "success",
			rmym: rmym,
			rrdm: rrdm,
			want: &service{rmym, rrdm},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewService(tt.rmym, tt.rrdm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetArticles(t *testing.T) {
	tests := []struct {
		name      string
		lfgar     FilterGetArticle
		setupMock func(*repositoryMySQLMock, *repositoryRedisMock)
		want      []Article
		wantErr   bool
	}{
		{
			name:  "succes get from redis",
			lfgar: FilterGetArticle{Limit: 10, Page: 0},
			setupMock: func(rmym *repositoryMySQLMock, rrdm *repositoryRedisMock) {
				rrdm.On("ReadArticles", ctx, mock.Anything).Return([]Article{{ID: 1, Author: Author{ID: 1, Name: "testname"}, Title: "titletest", Body: "bodytest"}}, nil).Once()
			},
			want:    []Article{{ID: 1, Author: Author{ID: 1, Name: "testname"}, Title: "titletest", Body: "bodytest"}},
			wantErr: false,
		},
		{
			name:  "succes key redis not found",
			lfgar: FilterGetArticle{Limit: 10, Page: 0},
			setupMock: func(rmym *repositoryMySQLMock, rrdm *repositoryRedisMock) {
				rrdm.On("ReadArticles", ctx, mock.Anything).Return([]Article{}, errors.New(config.RedisErrKeyDoesNotExist)).Once()
				rmym.On("ReadArticles", ctx, mock.Anything).Return([]Article{{ID: 1, Author: Author{ID: 1, Name: "testname"}, Title: "titletest", Body: "bodytest"}}, nil).Once()
				rrdm.On("CreateArticle", ctx, mock.Anything, mock.Anything).Return().Once()
			},
			want:    []Article{{ID: 1, Author: Author{ID: 1, Name: "testname"}, Title: "titletest", Body: "bodytest"}},
			wantErr: false,
		},
		{
			name:  "error mysql read articles",
			lfgar: FilterGetArticle{Limit: 10, Page: 0},
			setupMock: func(rmym *repositoryMySQLMock, rrdm *repositoryRedisMock) {
				rrdm.On("ReadArticles", ctx, mock.Anything).Return([]Article{}, errors.New(config.RedisErrKeyDoesNotExist)).Once()
				rmym.On("ReadArticles", ctx, mock.Anything).Return([]Article{}, errors.New("error read articles from mysql")).Once()
			},
			want:    []Article{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(rmym, rrdm)

			s := service{rmy: rmym, rre: rrdm}

			got, err := s.GetArticles(ctx, tt.lfgar)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.GetArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.GetArticles() = %v, want %v", got, tt.want)
			}
		})
	}
}
