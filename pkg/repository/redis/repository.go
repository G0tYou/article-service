package redis

import (
	"article/config"
	"article/pkg/listing"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type Storage struct {
	db *redis.Client
}

func NewStorage(cfgdb config.Redis) (*Storage, error) {
	var err error

	s := new(Storage)

	opts := &redis.Options{Addr: cfgdb.Addr, DB: cfgdb.Db}
	s.db = redis.NewClient(opts)

	_, err = s.db.Ping(context.Background()).Result()
	if err != nil {
		return s, err
	}

	log.Println("Redis connected", cfgdb.Addr, cfgdb.Db)

	return s, nil
}

func (s *Storage) CreateArticle(ctx context.Context, lars []listing.Article, lfgar listing.FilterGetArticle) {
	if lfgar.AuthorName == "" {
		lfgar.AuthorName = "-"
	}

	if lfgar.Search == "" {
		lfgar.Search = "-"
	}

	b, _ := json.Marshal(lars)

	s.db.Set(ctx, fmt.Sprintf(config.RedisKeyArticle, lfgar.AuthorName, lfgar.Search, lfgar.Limit, lfgar.Page), b, config.RedisTTLOneHour)
}

func (s *Storage) DeleteArticle(ctx context.Context) {
	s.db.FlushDB(ctx)
}

func (s *Storage) ReadArticles(ctx context.Context, lfgar listing.FilterGetArticle) ([]listing.Article, error) {
	var lars []listing.Article

	if lfgar.AuthorName == "" {
		lfgar.AuthorName = "-"
	}

	if lfgar.Search == "" {
		lfgar.Search = "-"
	}

	res := s.db.Get(ctx, fmt.Sprintf(config.RedisKeyArticle, lfgar.AuthorName, lfgar.Search, lfgar.Limit, lfgar.Page))
	json.Unmarshal([]byte(res.Val()), &lars)

	if res.Val() == "" {
		return lars, errors.New(config.RedisErrKeyDoesNotExist)
	}

	return lars, nil
}
