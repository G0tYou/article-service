package adding

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

type repositoryMySQLMock struct {
	mock.Mock
}

func (rmym *repositoryMySQLMock) CreateArticle(ctx context.Context, a Article) (int, error) {
	args := rmym.Called(ctx, a)

	return args.Get(0).(int), args.Error(1)
}

type repositoryRedisMock struct {
	mock.Mock
}

func (rrm *repositoryRedisMock) DeleteArticle(ctx context.Context) {
	rrm.Called(ctx)
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

func Test_service_AddArticle(t *testing.T) {
	tests := []struct {
		name      string
		a         Article
		setupMock func(*repositoryMySQLMock, *repositoryRedisMock)
		want      int
		wantErr   bool
	}{
		{
			name: "success",
			a:    Article{AuthorID: 1, Title: "testtitle", Body: "testbody"},
			setupMock: func(rms *repositoryMySQLMock, rrdm *repositoryRedisMock) {
				rmym.On("CreateArticle", ctx, Article{AuthorID: 1, Title: "testtitle", Body: "testbody"}).Return(1, nil).Once()
				rrdm.On("DeleteArticle", ctx).Return().Once()
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error mysql create article ",
			a:    Article{AuthorID: 1, Title: "testtitle", Body: "testbody"},
			setupMock: func(rms *repositoryMySQLMock, rrdm *repositoryRedisMock) {
				rmym.On("CreateArticle", ctx, Article{AuthorID: 1, Title: "testtitle", Body: "testbody"}).Return(0, errors.New("error create article")).Once()
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(rmym, rrdm)

			s := service{rmy: rmym, rre: rrdm}
			got, err := s.AddArticle(ctx, tt.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.AddArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("service.AddArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}
