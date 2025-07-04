package redis

import (
	"article/config"
	"article/pkg/listing"
	"context"
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name    string
		cfgdb   config.Redis
		wantErr bool
	}{
		{
			name:    "success redis reachable",
			cfgdb:   config.Redis{Addr: startMiniRedis(t), Db: 0},
			wantErr: false,
		},
		{
			name:    "error redis unreachable",
			cfgdb:   config.Redis{Addr: "127.0.0.1:6399", Db: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStorage(tt.cfgdb)
			fmt.Println(tt.cfgdb)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStorage_CreateArticle(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: startMiniRedis(t), DB: 0})

	tests := []struct {
		name  string
		ctx   context.Context
		s     *Storage
		lars  []listing.Article
		lfgar listing.FilterGetArticle
	}{
		{
			name:  "succes create article",
			ctx:   context.Background(),
			s:     &Storage{db: rdb},
			lars:  []listing.Article{},
			lfgar: listing.FilterGetArticle{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.CreateArticle(tt.ctx, tt.lars, tt.lfgar)
		})
	}
}

func TestStorage_DeleteArticle(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: startMiniRedis(t), DB: 0})

	tests := []struct {
		name string
		s    *Storage
		ctx  context.Context
	}{
		{
			name: "success delete article",
			s:    &Storage{db: rdb},
			ctx:  context.Background(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.DeleteArticle(tt.ctx)
		})
	}
}

func TestStorage_ReadArticles(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{Addr: startMiniRedis(t), DB: 0})

	tests := []struct {
		name    string
		s       *Storage
		ctx     context.Context
		lfgar   listing.FilterGetArticle
		wantErr bool
	}{
		{
			name:    "success read articles",
			s:       &Storage{db: rdb},
			ctx:     context.Background(),
			lfgar:   listing.FilterGetArticle{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.s.ReadArticles(tt.ctx, tt.lfgar)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ReadArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func startMiniRedis(t *testing.T) string {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	t.Cleanup(func() { mr.Close() })
	return mr.Addr()
}
